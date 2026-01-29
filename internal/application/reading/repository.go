package reading

import (
	"context"
	"errors"

	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

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
	ExistsDay(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate) (bool, error)
	GetDay(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate) (DayRow, bool, error)
	GetLastDayBefore(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate) (LastDayRow, bool, error)
	GetCurrentGoal(ctx context.Context, tx pgx.Tx, userID string) (int, bool, error)
	GetNextGoal(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate) (int, bool, error)
	GetDaysBetween(ctx context.Context, tx pgx.Tx, userID string, start, end readingDomain.LocalDate) (map[readingDomain.LocalDate]int, error)

	// writes
	AddPages(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate, delta int) (DayRow, error)
	InsertDay(ctx context.Context, tx pgx.Tx, userID string, date readingDomain.LocalDate, pagesTotal int, streakDays int) (DayRow, error)
	InsertGoal(ctx context.Context, tx pgx.Tx, userID string, pages int, startDate readingDomain.LocalDate) error
	UpdateGoalPages(ctx context.Context, tx pgx.Tx, userID string, pages int, startDate readingDomain.LocalDate) error
}
