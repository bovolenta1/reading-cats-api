package reading

import (
	"context"
	"time"

	appUser "reading-cats-api/internal/application/user"
	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/jackc/pgx/v5"
)

type RegisterReadingUseCase struct {
	repo        Repository
	userRepo    appUser.Repository
	defaultTZ   string
	graceHour   int
	goalDefault int
	clock       func() time.Time
}

func NewRegisterReadingUseCase(repo Repository, userRepo appUser.Repository, defaultTZ string) *RegisterReadingUseCase {
	return &RegisterReadingUseCase{
		repo:        repo,
		userRepo:    userRepo,
		defaultTZ:   defaultTZ,
		graceHour:   2,
		goalDefault: 5,
		clock:       time.Now,
	}
}

func (uc *RegisterReadingUseCase) Execute(ctx context.Context, in RegisterReadingInput) (RegisterReadingOutput, error) {
	// Lookup user by CognitoSub
	user, err := uc.userRepo.FindByCognitoSub(ctx, in.Claims.Sub)
	if err != nil {
		return RegisterReadingOutput{}, err
	}
	if user == nil {
		return RegisterReadingOutput{}, ErrUserNotFound
	}

	userID := user.ID
	loc, err := time.LoadLocation(uc.defaultTZ)
	if err != nil {
		return RegisterReadingOutput{}, err
	}
	now := uc.clock().In(loc)

	realDate := readingDomain.DateOf(now, loc)
	yesterday := realDate.AddDays(-1)

	var out RegisterReadingOutput

	err = uc.repo.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		hasYesterday, err := uc.repo.ExistsDay(ctx, tx, userID, yesterday)
		if err != nil {
			return err
		}

		targetDate := readingDomain.TargetDatePolicy{GraceHour: uc.graceHour}.Resolve(now, loc, hasYesterday)

		day, found, err := uc.repo.GetDay(ctx, tx, userID, targetDate)
		if err != nil {
			return err
		}

		if found {
			// adicionar páginas no mesmo dia, streak não muda
			day, err = uc.repo.AddPages(ctx, tx, userID, targetDate, int(in.Pages))
			if err != nil {
				return err
			}
		} else {
			last, hasLast, err := uc.repo.GetLastDayBefore(ctx, tx, userID, targetDate)
			if err != nil {
				return err
			}

			var lastDate readingDomain.LocalDate
			var lastStreak readingDomain.StreakDays
			if hasLast {
				lastDate = last.Date
				lastStreak = readingDomain.StreakDays(last.StreakDays)
			}

			newStreak := readingDomain.StreakPolicy{}.Next(targetDate, lastDate, lastStreak, hasLast)

			day, err = uc.repo.InsertDay(ctx, tx, userID, targetDate, int(in.Pages), int(newStreak))
			if err != nil {
				return err
			}
		}

		goal, hasGoal, err := uc.repo.GetCurrentGoal(ctx, tx, userID)
		if err != nil {
			return err
		}
		if !hasGoal {
			goal = uc.goalDefault
		}

		start := targetDate.AddDays(-6)
		byDate, err := uc.repo.GetDaysBetween(ctx, tx, userID, start, targetDate)
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

		out = RegisterReadingOutput{
			Progress: readingDomain.ReadingProgress{
				Day: readingDomain.DayProgress{
					Date:      targetDate.String(),
					Pages:     day.Pages,
					GoalPages: goal,
				},
				Streak: readingDomain.StreakProgress{
					CurrentDays: day.StreakDays,
				},
				Week: week,
			},
		}

		return nil
	})

	return out, err
}
