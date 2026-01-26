package reading

import (
	"context"
	"errors"

	app "reading-cats-api/internal/application/reading"
	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) WithTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) ExistsDay(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate) (bool, error) {
	q := `SELECT EXISTS(SELECT 1 FROM reading_day WHERE user_id=$1 AND reading_date=$2::date)`
	var exists bool
	err := tx.QueryRow(ctx, q, subID, date.String()).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) GetDay(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate) (app.DayRow, bool, error) {
	q := `SELECT pages_total, streak_days FROM reading_day WHERE user_id=$1 AND reading_date=$2::date LIMIT 1`
	var pages, streak int
	err := tx.QueryRow(ctx, q, subID, date.String()).Scan(&pages, &streak)
	if errors.Is(err, pgx.ErrNoRows) {
		return app.DayRow{}, false, nil
	}
	if err != nil {
		return app.DayRow{}, false, err
	}
	return app.DayRow{Date: date, Pages: pages, StreakDays: streak}, true, nil
}

func (r *PostgresRepository) GetLastDayBefore(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate) (app.LastDayRow, bool, error) {
	q := `
SELECT reading_date::text, streak_days
FROM reading_day
WHERE user_id=$1 AND reading_date < $2::date
ORDER BY reading_date DESC
LIMIT 1`
	var d string
	var streak int
	err := tx.QueryRow(ctx, q, subID, date.String()).Scan(&d, &streak)
	if errors.Is(err, pgx.ErrNoRows) {
		return app.LastDayRow{}, false, nil
	}
	if err != nil {
		return app.LastDayRow{}, false, err
	}
	return app.LastDayRow{Date: readingDomain.LocalDate(d), StreakDays: streak}, true, nil
}

func (r *PostgresRepository) AddPages(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate, delta int) (app.DayRow, error) {
	q := `
UPDATE reading_day
SET pages_total = pages_total + $3
WHERE user_id=$1 AND reading_date=$2::date
RETURNING pages_total, streak_days`
	var pages, streak int
	if err := tx.QueryRow(ctx, q, subID, date.String(), delta).Scan(&pages, &streak); err != nil {
		return app.DayRow{}, err
	}
	return app.DayRow{Date: date, Pages: pages, StreakDays: streak}, nil
}

func (r *PostgresRepository) InsertDay(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate, pagesTotal int, streakDays int) (app.DayRow, error) {
	q := `
INSERT INTO reading_day (user_id, reading_date, pages_total, streak_days, created_at, updated_at)
VALUES ($1, $2::date, $3, $4, now(), now())
RETURNING pages_total, streak_days`
	var pages, streak int
	err := tx.QueryRow(ctx, q, subID, date.String(), pagesTotal, streakDays).Scan(&pages, &streak)
	if err == nil {
		return app.DayRow{Date: date, Pages: pages, StreakDays: streak}, nil
	}

	// Race: outra request inseriu o mesmo dia entre GetDay e InsertDay
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return r.AddPages(ctx, tx, subID, date, pagesTotal)
	}

	return app.DayRow{}, err
}

func (r *PostgresRepository) GetGoalPagesOrDefault(ctx context.Context, tx pgx.Tx, subID string, def int) (int, error) {
	// Get the most recent goal where start_date <= today
	q := `
SELECT daily_pages 
FROM reading_goal 
WHERE user_id=$1 AND start_date <= now()
ORDER BY start_date DESC
LIMIT 1
`
	var v int
	err := tx.QueryRow(ctx, q, subID).Scan(&v)
	if errors.Is(err, pgx.ErrNoRows) {
		return def, nil
	}
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *PostgresRepository) GetDaysBetween(ctx context.Context, tx pgx.Tx, subID string, start, end readingDomain.LocalDate) (map[readingDomain.LocalDate]int, error) {
	q := `
SELECT reading_date::text, pages_total
FROM reading_day
WHERE user_id=$1 AND reading_date BETWEEN $2::date AND $3::date`
	rows, err := tx.Query(ctx, q, subID, start.String(), end.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[readingDomain.LocalDate]int{}
	for rows.Next() {
		var d string
		var p int
		if err := rows.Scan(&d, &p); err != nil {
			return nil, err
		}
		out[readingDomain.LocalDate(d)] = p
	}
	return out, rows.Err()
}

// UpdateGoal inserts a new goal record with start_date = tomorrow
func (r *PostgresRepository) UpdateGoal(ctx context.Context, subID string, pages readingDomain.Pages, validFrom readingDomain.LocalDate) error {
	// Insert new goal record with start_date at 00:00 of validFrom date
	q := `
INSERT INTO reading_goal (user_id, daily_pages, start_date, created_at)
VALUES ($1, $2, $3::date::timestamptz, now())
ON CONFLICT (user_id, start_date) DO UPDATE
SET daily_pages = $2
`
	_, err := r.pool.Exec(ctx, q, subID, int(pages), validFrom.String())
	return err
}
