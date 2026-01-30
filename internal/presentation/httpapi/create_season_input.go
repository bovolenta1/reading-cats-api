package httpapi

import (
	"encoding/json"
	"errors"
	"strings"

	appSeason "reading-cats-api/internal/application/season"

	"github.com/aws/aws-lambda-go/events"
)

type createSeasonBody struct {
	EndsAt   *string `json:"ends_at,omitempty"`
	Timezone string  `json:"timezone"`
}

func BuildCreateSeasonInput(event events.APIGatewayV2HTTPRequest) (appSeason.CreateSeasonInput, error) {
	// Extract claims
	claims, err := ExtractClaims(event)
	if err != nil {
		return appSeason.CreateSeasonInput{}, err
	}

	// Extract groupID from path: /v1/groups/{groupId}/seasons
	groupID := extractGroupIDFromPath(event.RawPath)
	if groupID == "" {
		return appSeason.CreateSeasonInput{}, errors.New("invalid group_id in path")
	}

	// Parse body
	var body createSeasonBody
	if err := json.Unmarshal([]byte(event.Body), &body); err != nil {
		return appSeason.CreateSeasonInput{}, errors.New("invalid request body")
	}

	if body.Timezone == "" {
		return appSeason.CreateSeasonInput{}, errors.New("timezone is required")
	}

	return appSeason.CreateSeasonInput{
		Claims:   claims,
		GroupID:  groupID,
		EndsAt:   body.EndsAt,
		Timezone: body.Timezone,
	}, nil
}

func extractGroupIDFromPath(path string) string {
	// Path format: /v1/groups/{groupId}/seasons
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 && parts[0] == "v1" && parts[1] == "groups" {
		return parts[2]
	}
	return ""
}
