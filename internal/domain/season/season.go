package season

import "time"

type Season struct {
	ID              string
	GroupID         string
	Status          Status
	StartedAt       *time.Time
	EndsAt          *time.Time
	Timezone        Timezone
	Metric          Metric
	CreatedByUserID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func New(
	id string,
	groupID string,
	status Status,
	startedAt *time.Time,
	endsAt *time.Time,
	timezone Timezone,
	metric Metric,
	createdByUserID string,
	createdAt time.Time,
) *Season {
	return &Season{
		ID:              id,
		GroupID:         groupID,
		Status:          status,
		StartedAt:       startedAt,
		EndsAt:          endsAt,
		Timezone:        timezone,
		Metric:          metric,
		CreatedByUserID: createdByUserID,
		CreatedAt:       createdAt,
		UpdatedAt:       createdAt,
	}
}
