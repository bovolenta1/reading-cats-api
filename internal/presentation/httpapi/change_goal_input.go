package httpapi

import (
	"encoding/json"
	"errors"
	"strings"

	app "reading-cats-api/internal/application/reading"
	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/aws/aws-lambda-go/events"
)

type changeGoalBody struct {
	Pages int `json:"pages"`
}

func BuildChangeGoalInput(event events.APIGatewayV2HTTPRequest) (app.ChangeGoalInput, error) {
	authInput, err := BuildEnsureMeInput(event)
	if err != nil {
		return app.ChangeGoalInput{}, err
	}

	s := strings.TrimSpace(event.Body)
	if s == "" {
		return app.ChangeGoalInput{}, errors.New("empty body")
	}

	var body changeGoalBody
	if err := json.Unmarshal([]byte(s), &body); err != nil {
		return app.ChangeGoalInput{}, errors.New("invalid json")
	}

	pages, err := readingDomain.NewPages(body.Pages)
	if err != nil {
		return app.ChangeGoalInput{}, err
	}

	return app.ChangeGoalInput{
		CognitoSub: string(authInput.Claims.Sub),
		Pages:      pages,
	}, nil
}
