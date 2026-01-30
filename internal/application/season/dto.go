package season

import (
	userDomain "reading-cats-api/internal/domain/user"
)

type CreateSeasonInput struct {
	Claims    userDomain.IDPClaims
	GroupID   string  `json:"group_id"`
	StartedAt *string `json:"started_at,omitempty"`
	EndsAt    *string `json:"ends_at,omitempty"`
	Timezone  string  `json:"timezone"`
}

type CreateSeasonOutput struct {
	ID              string  `json:"id"`
	GroupID         string  `json:"group_id"`
	Status          string  `json:"status"`
	EndsAt          *string `json:"ends_at,omitempty"`
	Timezone        string  `json:"timezone"`
	Metric          string  `json:"metric"`
	CreatedByUserID string  `json:"created_by_user_id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
