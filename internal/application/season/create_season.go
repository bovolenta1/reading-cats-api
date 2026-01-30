package season

import (
	"context"
	"fmt"
	"time"

	appUser "reading-cats-api/internal/application/user"
	domainSeason "reading-cats-api/internal/domain/season"

	"github.com/google/uuid"
)

type CreateSeasonUseCase struct {
	repo     Repository
	userRepo appUser.Repository
}

func NewCreateSeasonUseCase(repo Repository, userRepo appUser.Repository) *CreateSeasonUseCase {
	return &CreateSeasonUseCase{repo: repo, userRepo: userRepo}
}

func (uc *CreateSeasonUseCase) Execute(ctx context.Context, in CreateSeasonInput) (CreateSeasonOutput, error) {
	// Find user by CognitoSub to get UUID
	user, err := uc.userRepo.FindByCognitoSub(ctx, in.Claims.Sub)
	if err != nil {
		return CreateSeasonOutput{}, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return CreateSeasonOutput{}, fmt.Errorf("user not found")
	}

	// Validate timezone
	timezone, err := domainSeason.NewTimezone(in.Timezone)
	if err != nil {
		return CreateSeasonOutput{}, err
	}

	// Parse ends_at if provided
	var endsAt *time.Time
	if in.EndsAt != nil && *in.EndsAt != "" {
		t, err := time.Parse(time.RFC3339, *in.EndsAt)
		if err != nil {
			return CreateSeasonOutput{}, fmt.Errorf("invalid ends_at format: %w", err)
		}
		endsAt = &t
	}

	// Create season domain entity (DRAFT without started_at)
	seasonID := uuid.NewString()
	now := time.Now().UTC()
	s := domainSeason.New(
		seasonID,
		in.GroupID,
		domainSeason.StatusDraft,
		nil,    // startedAt - set later when activated
		endsAt, // endsAt - defined at creation
		timezone,
		domainSeason.MetricCheckinsPerDay,
		user.ID,
		now,
	)

	// Insert season
	if err := uc.repo.Insert(ctx, s); err != nil {
		return CreateSeasonOutput{}, fmt.Errorf("failed to insert season: %w", err)
	}

	// Build output
	output := CreateSeasonOutput{
		ID:              s.ID,
		GroupID:         s.GroupID,
		Status:          s.Status.String(),
		Timezone:        string(s.Timezone),
		Metric:          s.Metric.String(),
		CreatedByUserID: s.CreatedByUserID,
		CreatedAt:       s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       s.UpdatedAt.Format(time.RFC3339),
	}

	if s.EndsAt != nil {
		formatted := s.EndsAt.Format(time.RFC3339)
		output.EndsAt = &formatted
	}

	return output, nil
}
