package reading

import (
	"context"
	"time"

	appUser "reading-cats-api/internal/application/user"
	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/jackc/pgx/v5"
)

type GetReadingProgressUseCase struct {
	repo        Repository
	userRepo    appUser.Repository
	defaultTZ   string
	graceHour   int
	goalDefault int
	clock       func() time.Time
}

func NewGetReadingProgressUseCase(repo Repository, userRepo appUser.Repository, defaultTZ string) *GetReadingProgressUseCase {
	return &GetReadingProgressUseCase{
		repo:        repo,
		userRepo:    userRepo,
		defaultTZ:   defaultTZ,
		graceHour:   2,
		goalDefault: 5,
		clock:       time.Now,
	}
}

func (uc *GetReadingProgressUseCase) Execute(ctx context.Context, in GetReadingProgressInput) (GetReadingProgressOutput, error) {
	// Lookup user by CognitoSub
	user, err := uc.userRepo.FindByCognitoSub(ctx, in.Claims.Sub)
	if err != nil {
		return GetReadingProgressOutput{}, err
	}
	if user == nil {
		return GetReadingProgressOutput{}, ErrUserNotFound
	}

	userID := user.ID

	loc, err := time.LoadLocation(uc.defaultTZ)
	if err != nil {
		return GetReadingProgressOutput{}, err
	}
	now := uc.clock().In(loc)

	realDate := readingDomain.DateOf(now, loc)
	yesterday := realDate.AddDays(-1)

	var out GetReadingProgressOutput

	err = uc.repo.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		hasYesterday, err := uc.repo.ExistsDay(ctx, tx, userID, yesterday)
		if err != nil {
			return err
		}

		// check if user loses yesterday's reading due to grace period
		targetDate := readingDomain.TargetDatePolicy{GraceHour: uc.graceHour}.Resolve(now, loc, hasYesterday)

		day, found, err := uc.repo.GetDay(ctx, tx, userID, targetDate)
		if err != nil {
			return err
		}

		goal, hasGoal, err := uc.repo.GetCurrentGoal(ctx, tx, userID)
		if err != nil {
			return err
		}
		if !hasGoal {
			goal = uc.goalDefault
		}

		currentGoal := &GoalRecord{
			DailyPages: goal,
			ValidFrom:  targetDate.String(),
		}

		nextGoal := uc.buildNextGoal(ctx, tx, userID, realDate, goal)

		streak := 0
		pagesToday := 0

		if found {
			streak = day.StreakDays
			pagesToday = day.Pages
		} else {
			last, hasLast, err := uc.repo.GetLastDayBefore(ctx, tx, userID, targetDate)
			if err != nil {
				return err
			}
			if hasLast && last.Date == targetDate.AddDays(-1) {
				streak = last.StreakDays
			}
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
			CurrentGoal: currentGoal,
			NextGoal:    nextGoal,
		}

		return nil
	})

	return out, err
}

func (uc *GetReadingProgressUseCase) buildNextGoal(ctx context.Context, tx pgx.Tx, userID string, targetDate readingDomain.LocalDate, currentGoalPages int) *GoalRecord {
	nextDate := targetDate.AddDays(1)
	nextGoalPages, hasNextGoal, _ := uc.repo.GetNextGoal(ctx, tx, userID, nextDate)
	if !hasNextGoal {
		nextGoalPages = 0
	}

	if nextGoalPages != currentGoalPages {
		return &GoalRecord{
			DailyPages: nextGoalPages,
			ValidFrom:  nextDate.String(),
		}
	}

	return nil
}
