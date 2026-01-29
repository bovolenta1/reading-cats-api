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

func (r *PostgresRepository) ExistsDay(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate) (bool, error) {
	q := `SELECT EXISTS(SELECT 1 FROM user_checkins WHERE user_id=$1::uuid AND local_date=$2::date)`
	var exists bool
	err := tx.QueryRow(ctx, q, userID, date.String()).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) GetDay(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate) (app.DayRow, bool, error) {
	q := `SELECT pages_total, streak_days FROM user_checkins WHERE user_id=$1::uuid AND local_date=$2::date LIMIT 1`
	var pages, streak int
	err := tx.QueryRow(ctx, q, userID, date.String()).Scan(&pages, &streak)
	if errors.Is(err, pgx.ErrNoRows) {
		return app.DayRow{}, false, nil
	}
	if err != nil {
		return app.DayRow{}, false, err
	}
	return app.DayRow{Date: date, Pages: pages, StreakDays: streak}, true, nil
}

func (r *PostgresRepository) GetLastDayBefore(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate) (app.LastDayRow, bool, error) {
	q := `
SELECT local_date::text, streak_days
FROM user_checkins
WHERE user_id=$1::uuid AND local_date < $2::date
ORDER BY local_date DESC
LIMIT 1`
	var d string
	var streak int
	err := tx.QueryRow(ctx, q, userID, date.String()).Scan(&d, &streak)
	if errors.Is(err, pgx.ErrNoRows) {
		return app.LastDayRow{}, false, nil
	}
	if err != nil {
		return app.LastDayRow{}, false, err
	}
	return app.LastDayRow{Date: readingDomain.LocalDate(d), StreakDays: streak}, true, nil
}

func (r *PostgresRepository) AddPages(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate, delta int) (app.DayRow, error) {
	q := `
UPDATE user_checkins
SET pages_total = pages_total + $3
WHERE user_id=$1::uuid AND local_date=$2::date
RETURNING pages_total, streak_days`
	var pages, streak int
	if err := tx.QueryRow(ctx, q, userID, date.String(), delta).Scan(&pages, &streak); err != nil {
		return app.DayRow{}, err
	}
	return app.DayRow{Date: date, Pages: pages, StreakDays: streak}, nil
}

func (r *PostgresRepository) InsertDay(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate, pagesTotal int, streakDays int) (app.DayRow, error) {
	q := `
INSERT INTO user_checkins (user_id, local_date, pages_total, streak_days, created_at, updated_at)
VALUES ($1::uuid, $2::date, $3, $4, now(), now())
RETURNING pages_total, streak_days`
	var pages, streak int
	err := tx.QueryRow(ctx, q, userID, date.String(), pagesTotal, streakDays).Scan(&pages, &streak)
	if err == nil {
		return app.DayRow{Date: date, Pages: pages, StreakDays: streak}, nil
	}

	// Race: outra request inseriu o mesmo dia entre GetDay e InsertDay
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return r.AddPages(ctx, tx, userID, date, pagesTotal)
	}

	return app.DayRow{}, err
}

func (r *PostgresRepository) GetDaysBetween(ctx context.Context, tx pgx.Tx, userID string, start, end readingDomain.LocalDate) (map[readingDomain.LocalDate]int, error) {
	q := `
SELECT local_date::text, pages_total
FROM user_checkins
WHERE user_id=$1::uuid AND local_date BETWEEN $2::date AND $3::date`
	rows, err := tx.Query(ctx, q, userID, start.String(), end.String())
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

// GetCurrentGoal retorna o goal vigente (start_date <= agora) ou nil se não existe
func (r *PostgresRepository) GetCurrentGoal(ctx context.Context, tx pgx.Tx, userID string) (int, bool, error) {
	q := `
SELECT daily_pages 
FROM reading_goal 
WHERE user_id=$1::uuid AND start_date <= now()
ORDER BY start_date DESC
LIMIT 1
`
	var pages int
	err := tx.QueryRow(ctx, q, userID).Scan(&pages)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return pages, true, nil
}

// GetNextGoal retorna o goal para uma data futura ou nil se não existe
func (r *PostgresRepository) GetNextGoal(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate) (int, bool, error) {
	q := `
SELECT daily_pages 
FROM reading_goal 
WHERE user_id=$1::uuid AND start_date::date = $2::date
LIMIT 1
`
	var pages int
	err := tx.QueryRow(ctx, q, userID, date.String()).Scan(&pages)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return pages, true, nil
}

// InsertGoal inserts a new goal record
func (r *PostgresRepository) InsertGoal(ctx context.Context, tx pgx.Tx, userID string, pages int, startDate readingDomain.LocalDate) error {
	q := `
INSERT INTO reading_goal (user_id, daily_pages, start_date, created_at)
VALUES ($1::uuid, $2, $3::date, now())
`
	_, err := tx.Exec(ctx, q, userID, pages, startDate.String())
	return err
}

// UpdateGoalPages updates the pages for a specific goal date
func (r *PostgresRepository) UpdateGoalPages(ctx context.Context, tx pgx.Tx, userID string, pages int, startDate readingDomain.LocalDate) error {
	q := `
UPDATE reading_goal
SET daily_pages = $1
WHERE user_id = $2::uuid AND start_date::date = $3::date
`
	_, err := tx.Exec(ctx, q, pages, userID, startDate.String())
	return err
}
