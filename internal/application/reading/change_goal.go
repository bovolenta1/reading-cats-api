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
	nextDay := today.AddDays(1)

	var currentGoal *GoalRecord
	var nextGoal *GoalRecord

	err = uc.repo.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		goalPages, _ := uc.repo.GetGoalPagesOrDefault(ctx, tx, in.CognitoSub, 5)
		currentGoal = &GoalRecord{
			DailyPages: goalPages,
			ValidFrom:  today.String(),
		}

		if err := uc.repo.UpdateGoal(ctx, in.CognitoSub, in.Pages, nextDay); err != nil {
			return err
		}

		nextGoal = &GoalRecord{
			DailyPages: int(in.Pages),
			ValidFrom:  nextDay.String(),
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
