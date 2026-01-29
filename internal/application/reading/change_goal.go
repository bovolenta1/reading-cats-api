package reading

import (
	"context"
	"time"

	appUser "reading-cats-api/internal/application/user"
	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/jackc/pgx/v5"
)

type ChangeGoalUseCase struct {
	repo     Repository
	userRepo appUser.Repository
	timezone string
}

func NewChangeGoalUseCase(repo Repository, userRepo appUser.Repository, timezone string) *ChangeGoalUseCase {
	return &ChangeGoalUseCase{
		repo:     repo,
		userRepo: userRepo,
		timezone: timezone,
	}
}

func (uc *ChangeGoalUseCase) Execute(ctx context.Context, in ChangeGoalInput) (ChangeGoalOutput, error) {
	// Lookup user by CognitoSub
	user, err := uc.userRepo.FindByCognitoSub(ctx, in.Claims.Sub)
	if err != nil {
		return ChangeGoalOutput{}, err
	}
	if user == nil {
		return ChangeGoalOutput{}, ErrUserNotFound
	}

	userID := user.ID

	loc, err := time.LoadLocation(uc.timezone)
	if err != nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)
	today := readingDomain.DateOf(now, loc)
	tomorrow := today.AddDays(1)

	var currentGoal *GoalRecord
	var nextGoal *GoalRecord

	err = uc.repo.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		currentPages, hasCurrentGoal, err := uc.repo.GetCurrentGoal(ctx, tx, userID)
		if err != nil {
			return err
		}

		if !hasCurrentGoal {
			if err := uc.repo.InsertGoal(ctx, tx, userID, 5, today); err != nil {
				return err
			}
			currentPages = 5
		}

		currentGoal = &GoalRecord{
			DailyPages: currentPages,
			ValidFrom:  today.String(),
		}

		nextPages, hasNextGoal, err := uc.repo.GetNextGoal(ctx, tx, userID, tomorrow)
		if err != nil {
			return err
		}

		if !hasNextGoal {
			if err := uc.repo.InsertGoal(ctx, tx, userID, int(in.Pages), tomorrow); err != nil {
				return err
			}
			nextPages = int(in.Pages)
		} else if nextPages != int(in.Pages) {
			if err := uc.repo.UpdateGoalPages(ctx, tx, userID, int(in.Pages), tomorrow); err != nil {
				return err
			}
			nextPages = int(in.Pages)
		}

		nextGoal = &GoalRecord{
			DailyPages: nextPages,
			ValidFrom:  tomorrow.String(),
		}

		return nil
	})

	if err != nil {
		return ChangeGoalOutput{}, err
	}

	return ChangeGoalOutput{
		CurrentGoal: currentGoal,
		NextGoal:    nextGoal,
	}, nil
}
