package httpapi

import (
	"encoding/json"
	"errors"

	app "reading-cats-api/internal/application/reading"
	readingDomain "reading-cats-api/internal/domain/reading"

	"github.com/aws/aws-lambda-go/events"
)

type registerReadingBody struct {
	Pages int `json:"pages"`
}

func BuildRegisterReadingInput(event events.APIGatewayV2HTTPRequest) (app.RegisterReadingInput, error) {
	// Extract Claims from event
	claims, err := ExtractClaims(event)
	if err != nil {
		return app.RegisterReadingInput{}, err
	}

	// Parse body
	var body registerReadingBody
	if err := json.Unmarshal([]byte(event.Body), &body); err != nil {
		return app.RegisterReadingInput{}, errors.New("invalid request body")
	}

	pagesVO, err := readingDomain.NewPages(body.Pages)
	if err != nil {
		return app.RegisterReadingInput{}, err
	}

	return app.RegisterReadingInput{
		Claims: claims,
		Pages:  pagesVO,
	}, nil
}
