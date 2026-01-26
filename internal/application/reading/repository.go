package reading

import (
	"context"

	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/jackc/pgx/v5"
)

type DayRow struct {
	Date       readingDomain.LocalDate
	Pages      int
	StreakDays int
}

type LastDayRow struct {
	Date       readingDomain.LocalDate
	StreakDays int
}

type Repository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error

	// reads
	ExistsDay(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate) (bool, error)
	GetDay(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate) (DayRow, bool, error)
	GetLastDayBefore(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate) (LastDayRow, bool, error)
	GetGoalPagesOrDefault(ctx context.Context, tx pgx.Tx, subID string, def int) (int, error)
	GetDaysBetween(ctx context.Context, tx pgx.Tx, subID string, start, end readingDomain.LocalDate) (map[readingDomain.LocalDate]int, error)

	// writes
	AddPages(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate, delta int) (DayRow, error)
	InsertDay(ctx context.Context, tx pgx.Tx, subID string, date readingDomain.LocalDate, pagesTotal int, streakDays int) (DayRow, error)
	UpdateGoal(ctx context.Context, subID string, pages readingDomain.Pages, validFrom readingDomain.LocalDate) error
}
