package reading

import (
	"context"
	"time"

	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/jackc/pgx/v5"
)

type GetReadingProgressUseCase struct {
	repo        Repository
	defaultTZ   string
	graceHour   int
	goalDefault int
	clock       func() time.Time
}

func NewGetReadingProgressUseCase(repo Repository, defaultTZ string) *GetReadingProgressUseCase {
	return &GetReadingProgressUseCase{
		repo:        repo,
		defaultTZ:   defaultTZ,
		graceHour:   2,
		goalDefault: 5,
		clock:       time.Now,
	}
}

func (uc *GetReadingProgressUseCase) Execute(ctx context.Context, in GetReadingProgressInput) (GetReadingProgressOutput, error) {
	sub := string(in.Claims.Sub)

	loc, err := time.LoadLocation(uc.defaultTZ)
	if err != nil {
		return GetReadingProgressOutput{}, err
	}
	now := uc.clock().In(loc)

	realDate := readingDomain.DateOf(now, loc)
	yesterday := realDate.AddDays(-1)

	var out GetReadingProgressOutput

	err = uc.repo.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		hasYesterday, err := uc.repo.ExistsDay(ctx, tx, sub, yesterday)
		if err != nil {
			return err
		}

		targetDate := readingDomain.TargetDatePolicy{GraceHour: uc.graceHour}.Resolve(now, loc, hasYesterday)

		day, found, err := uc.repo.GetDay(ctx, tx, sub, targetDate)
		if err != nil {
			return err
		}

		goal, err := uc.repo.GetGoalPagesOrDefault(ctx, tx, sub, uc.goalDefault)
		if err != nil {
			return err
		}

		// streak pra snapshot:
		streak := 0
		pagesToday := 0

		if found {
			streak = day.StreakDays
			pagesToday = day.Pages
		} else {
			last, hasLast, err := uc.repo.GetLastDayBefore(ctx, tx, sub, targetDate)
			if err != nil {
				return err
			}
			if hasLast && last.Date == targetDate.AddDays(-1) {
				streak = last.StreakDays
			}
		}

		start := targetDate.AddDays(-6)
		byDate, err := uc.repo.GetDaysBetween(ctx, tx, sub, start, targetDate)
		if err != nil {
			return err
		}

		week := make([]readingDomain.WeekDayProgress, 0, 7)
		for i := 0; i < 7; i++ {
			d := start.AddDays(i)
			p := byDate[d]
			week = append(week, readingDomain.WeekDayProgress{
				Date:    d.String(),
				Pages:   p,
				Checked: p > 0,
			})
		}

		out = GetReadingProgressOutput{
			Progress: readingDomain.ReadingProgress{
				Day: readingDomain.DayProgress{
					Date:      targetDate.String(),
					Pages:     pagesToday,
					GoalPages: goal,
				},
				Streak: readingDomain.StreakProgress{
					CurrentDays: streak,
				},
				Week: week,
			},
		}

		return nil
	})

	return out, err
}
