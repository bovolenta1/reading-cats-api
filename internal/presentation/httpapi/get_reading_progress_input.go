package httpapi

import (
	app "reading-cats-api/internal/application/reading"

	"github.com/aws/aws-lambda-go/events"
)

func BuildGetReadingProgressInput(event events.APIGatewayV2HTTPRequest) (app.GetReadingProgressInput, error) {
	meIn, err := BuildEnsureMeInput(event)
	if err != nil {
		return app.GetReadingProgressInput{}, err
	}
	return app.GetReadingProgressInput{
		Claims: meIn.Claims,
	}, nil
}
