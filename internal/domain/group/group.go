package group

import "time"

type Group struct {
	ID              string
	Name            GroupName
	IconID          IconID
	Visibility      Visibility
	MaxMembers      int
	CreatedByUserID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func New(
	id string,
	name GroupName,
	iconID IconID,
	visibility Visibility,
	maxMembers int,
	createdByUserID string,
	createdAt time.Time,
) *Group {
	return &Group{
		ID:              id,
		Name:            name,
		IconID:          iconID,
		Visibility:      visibility,
		MaxMembers:      maxMembers,
		CreatedByUserID: createdByUserID,
		CreatedAt:       createdAt,
		UpdatedAt:       createdAt,
	}
}
