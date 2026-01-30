package httpapi

import (
	"encoding/json"
	"errors"

	appGroup "reading-cats-api/internal/application/group"

	"github.com/aws/aws-lambda-go/events"
)

type createGroupBody struct {
	Name       string `json:"name"`
	IconID     string `json:"icon_id"`
	MaxMembers *int   `json:"max_members,omitempty"`
}

func BuildCreateGroupInput(event events.APIGatewayV2HTTPRequest) (appGroup.CreateGroupInput, error) {
	// Extract claims
	claims, err := ExtractClaims(event)
	if err != nil {
		return appGroup.CreateGroupInput{}, err
	}

	// Parse body
	var body createGroupBody
	if err := json.Unmarshal([]byte(event.Body), &body); err != nil {
		return appGroup.CreateGroupInput{}, errors.New("invalid request body")
	}

	if body.Name == "" || body.IconID == "" {
		return appGroup.CreateGroupInput{}, errors.New("name and icon_id are required")
	}

	return appGroup.CreateGroupInput{
		Claims:     claims,
		Name:       body.Name,
		IconID:     body.IconID,
		MaxMembers: body.MaxMembers,
	}, nil
}
