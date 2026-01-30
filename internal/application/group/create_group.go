package group

import (
	"context"
	"fmt"
	"time"

	appUser "reading-cats-api/internal/application/user"
	domainGroup "reading-cats-api/internal/domain/group"

	"github.com/google/uuid"
)

type CreateGroupUseCase struct {
	repo     Repository
	userRepo appUser.Repository
}

func NewCreateGroupUseCase(repo Repository, userRepo appUser.Repository) *CreateGroupUseCase {
	return &CreateGroupUseCase{repo: repo, userRepo: userRepo}
}

func (uc *CreateGroupUseCase) Execute(ctx context.Context, in CreateGroupInput) (CreateGroupOutput, error) {
	user, err := uc.userRepo.FindByCognitoSub(ctx, in.Claims.Sub)
	if err != nil {
		return CreateGroupOutput{}, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return CreateGroupOutput{}, fmt.Errorf("user not found")
	}

	// Validate input
	name, err := domainGroup.NewGroupName(in.Name)
	if err != nil {
		return CreateGroupOutput{}, err
	}

	iconID, err := domainGroup.NewIconID(in.IconID)
	if err != nil {
		return CreateGroupOutput{}, err
	}

	// Create group domain entity
	groupID := uuid.NewString()
	now := time.Now().UTC()
	maxMembers := 5 // default
	if in.MaxMembers != nil && *in.MaxMembers > 0 {
		maxMembers = *in.MaxMembers
	}
	g := domainGroup.New(
		groupID,
		name,
		iconID,
		domainGroup.VisibilityInviteOnly,
		maxMembers,
		user.ID,
		now,
	)

	// Insert group
	if err := uc.repo.Insert(ctx, g); err != nil {
		return CreateGroupOutput{}, fmt.Errorf("failed to insert group: %w", err)
	}

	// Add creator as ADMIN member
	if err := uc.repo.AddMember(ctx, groupID, user.ID, "ADMIN"); err != nil {
		return CreateGroupOutput{}, fmt.Errorf("failed to add user as group member: %w", err)
	}

	// Return DTO
	return CreateGroupOutput{
		ID:              g.ID,
		Name:            string(g.Name),
		IconID:          string(g.IconID),
		Visibility:      g.Visibility.String(),
		MaxMembers:      g.MaxMembers,
		CreatedByUserID: g.CreatedByUserID,
		CreatedAt:       g.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       g.UpdatedAt.Format(time.RFC3339),
	}, nil
}
