package reading

import (
	"context"
	"time"

	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/jackc/pgx/v5"
)

type ChangeGoalUseCase struct {
	repo     Repository
	timezone string
}

func NewChangeGoalUseCase(repo Repository, timezone string) *ChangeGoalUseCase {
	return &ChangeGoalUseCase{
		repo:     repo,
		timezone: timezone,
	}
}

func (uc *ChangeGoalUseCase) Execute(ctx context.Context, in ChangeGoalInput) (ChangeGoalOutput, error) {
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
		currentPages, hasCurrentGoal, err := uc.repo.GetCurrentGoal(ctx, tx, in.CognitoSub)
		if err != nil {
			return err
		}

		if !hasCurrentGoal {
			if err := uc.repo.InsertGoal(ctx, tx, in.CognitoSub, 5, today); err != nil {
				return err
			}
			currentPages = 5
		}

		currentGoal = &GoalRecord{
			DailyPages: currentPages,
			ValidFrom:  today.String(),
		}

		nextPages, hasNextGoal, err := uc.repo.GetNextGoal(ctx, tx, in.CognitoSub, tomorrow)
		if err != nil {
			return err
		}

		if !hasNextGoal {
			if err := uc.repo.InsertGoal(ctx, tx, in.CognitoSub, int(in.Pages), tomorrow); err != nil {
				return err
			}
			nextPages = int(in.Pages)
		} else if nextPages != int(in.Pages) {
			if err := uc.repo.UpdateGoalPages(ctx, tx, in.CognitoSub, int(in.Pages), tomorrow); err != nil {
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
