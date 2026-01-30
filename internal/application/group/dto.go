package group

import (
	userDomain "reading-cats-api/internal/domain/user"
)

type CreateGroupInput struct {
	Claims     userDomain.IDPClaims
	Name       string `json:"name"`
	IconID     string `json:"icon_id"`
	MaxMembers *int   `json:"max_members,omitempty"`
}

type CreateGroupOutput struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	IconID          string `json:"icon_id"`
	Visibility      string `json:"visibility"`
	MaxMembers      int    `json:"max_members"`
	CreatedByUserID string `json:"created_by_user_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
