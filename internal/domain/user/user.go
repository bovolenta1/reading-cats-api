package user

import (
	"time"

	"github.com/google/uuid"
)

type ProfileSource string

const (
	ProfileSourceIDP  ProfileSource = "idp"
	ProfileSourceUser ProfileSource = "user"
)

type User struct {
	ID            string
	CognitoSub    CognitoSub
	Email         Email
	DisplayName   DisplayName
	AvatarURL     AvatarURL
	ProfileSource ProfileSource

	CreatedAt time.Time
	UpdatedAt time.Time
}

type IDPClaims struct {
	Sub     CognitoSub
	Email   Email
	Name    DisplayName
	Picture AvatarURL
}

func NewFromIDP(claims IDPClaims) User {
	now := time.Now().UTC()
	return User{
		ID:            uuid.NewString(),
		CognitoSub:    claims.Sub,
		Email:         claims.Email,
		DisplayName:   claims.Name,
		AvatarURL:     claims.Picture,
		ProfileSource: ProfileSourceIDP,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
