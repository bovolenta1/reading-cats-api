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

type body struct {
	Pages int `json:"pages"`
}

func BuildRegisterReadingInput(event events.APIGatewayV2HTTPRequest) (app.RegisterReadingInput, error) {
	// reaproveita claims do core httpapi
	meIn, err := BuildEnsureMeInput(event)
	if err != nil {
		return app.RegisterReadingInput{}, err
	}

	raw, err := readBody(event)
	if err != nil {
		return app.RegisterReadingInput{}, err
	}

	var b body
	if err := json.Unmarshal(raw, &b); err != nil {
		return app.RegisterReadingInput{}, err
	}

	pagesVO, err := readingDomain.NewPages(b.Pages)
	if err != nil {
		return app.RegisterReadingInput{}, readingDomain.ErrInvalidPages
	}

	return app.RegisterReadingInput{
		Claims: meIn.Claims,
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
