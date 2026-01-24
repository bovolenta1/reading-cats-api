package user

import (
	"context"

	domain "reading-cats-api/internal/domain/user"
)

type EnsureMeUseCase struct {
	repo Repository
}

func NewEnsureMeUseCase(repo Repository) *EnsureMeUseCase {
	return &EnsureMeUseCase{repo: repo}
}

type Input struct {
	Claims domain.IDPClaims
}

func (uc *EnsureMeUseCase) Execute(ctx context.Context, in Input) (MeDTO, error) {
	existing, err := uc.repo.FindByCognitoSub(ctx, in.Claims.Sub)
	if err != nil {
		return MeDTO{}, err
	}

	if existing == nil {
		u := domain.NewFromIDP(in.Claims)
		if err := uc.repo.Insert(ctx, &u); err != nil {
			return MeDTO{}, err
		}
		return toMeDTO(u), nil
	}

	return toMeDTO(*existing), nil
}

func toMeDTO(u domain.User) MeDTO {
	return MeDTO{
		ID:          u.ID,
		CognitoSub:  string(u.CognitoSub),
		Email:       string(u.Email),
		DisplayName: string(u.DisplayName),
		AvatarURL:   string(u.AvatarURL),
		Source:      string(u.ProfileSource),
	}
}
