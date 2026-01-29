package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

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

func readBody(event events.APIGatewayV2HTTPRequest) ([]byte, error) {
	s := strings.TrimSpace(event.Body)
	if s == "" {
		return nil, errors.New("empty body")
	}

	if event.IsBase64Encoded {
		// APIGW v2 costuma mandar base64 standard
		b, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			// fallback pra raw url encoding
			b2, err2 := base64.RawStdEncoding.DecodeString(s)
			if err2 != nil {
				return nil, err
			}
			return b2, nil
		}
		return b, nil
	}

	return []byte(s), nil
}
