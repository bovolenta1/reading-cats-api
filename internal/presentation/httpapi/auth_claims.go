package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"

	app "reading-cats-api/internal/application/user"
	domain "reading-cats-api/internal/domain/user"

	"github.com/aws/aws-lambda-go/events"
)

var ErrUnauthorized = errors.New("unauthorized")

// ExtractClaims extrai IDPClaims direto do evento (sem criar app.Input).
// Funciona tanto em PROD (API Gateway JWT Authorizer) quanto em DEV (SAM local).
func ExtractClaims(event events.APIGatewayV2HTTPRequest) (domain.IDPClaims, error) {
	// 1) PROD: claims do API Gateway JWT Authorizer (quando existir)
	if event.RequestContext.Authorizer != nil && event.RequestContext.Authorizer.JWT != nil {
		c := event.RequestContext.Authorizer.JWT.Claims

		sub := strings.TrimSpace(c["sub"])
		if sub != "" {
			name := strings.TrimSpace(c["name"])
			if name == "" {
				name = strings.TrimSpace(joinName(c["given_name"], c["family_name"]))
			}

			return buildIDPClaims(
				sub,
				strings.TrimSpace(c["email"]),
				name,
				strings.TrimSpace(c["picture"]),
			)
		}
	}

	// 2) SAM local: fallback decodificando o JWT do Authorization header
	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		token := bearerToken(event.Headers)
		if token == "" {
			return domain.IDPClaims{}, ErrUnauthorized
		}

		p, ok := decodeJwtPayload(token)
		if !ok || strings.TrimSpace(p.Sub) == "" {
			return domain.IDPClaims{}, ErrUnauthorized
		}

		name := strings.TrimSpace(p.Name)
		if name == "" {
			name = strings.TrimSpace(joinName(p.GivenName, p.FamilyName))
		}

		return buildIDPClaims(
			p.Sub,
			p.Email,
			name,
			p.Picture,
		)
	}

	return domain.IDPClaims{}, ErrUnauthorized
}

type jwtPayload struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

func BuildEnsureMeInput(event events.APIGatewayV2HTTPRequest) (app.Input, error) {
	// 1) PROD: claims do API Gateway JWT Authorizer (quando existir)
	if event.RequestContext.Authorizer != nil && event.RequestContext.Authorizer.JWT != nil {
		c := event.RequestContext.Authorizer.JWT.Claims

		sub := strings.TrimSpace(c["sub"])
		if sub != "" {
			name := strings.TrimSpace(c["name"])
			if name == "" {
				name = strings.TrimSpace(joinName(c["given_name"], c["family_name"]))
			}

			claims, err := buildIDPClaims(
				sub,
				strings.TrimSpace(c["email"]),
				name,
				strings.TrimSpace(c["picture"]),
			)

			if err != nil {
				return app.Input{}, err
			}

			return app.Input{Claims: claims}, nil
		}
	}

	// 2) SAM local: fallback decodificando o JWT do Authorization header
	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		token := bearerToken(event.Headers)
		if token == "" {
			return app.Input{}, ErrUnauthorized
		}

		p, ok := decodeJwtPayload(token)
		if !ok || strings.TrimSpace(p.Sub) == "" {
			return app.Input{}, ErrUnauthorized
		}

		name := strings.TrimSpace(p.Name)
		if name == "" {
			name = strings.TrimSpace(joinName(p.GivenName, p.FamilyName))
		}

		claims, err := buildIDPClaims(
			p.Sub,
			p.Email,
			name,
			p.Picture,
		)
		if err != nil {
			return app.Input{}, err
		}

		return app.Input{Claims: claims}, nil
	}

	return app.Input{}, ErrUnauthorized
}

// buildIDPClaims cria IDPClaims usando Value Objects.
// Regra prática:
// - sub: obrigatório e válido, senão -> unauthorized
// - email/nome/avatar: opcionais; se inválidos, a gente ignora (fica vazio) pra não quebrar login.
func buildIDPClaims(sub, email, name, picture string) (domain.IDPClaims, error) {
	subVO, err := domain.NewCognitoSub(sub)
	if err != nil {
		return domain.IDPClaims{}, ErrUnauthorized
	}

	emailVO, err := domain.NewEmail(email)
	if err != nil {
		emailVO = ""
	}

	nameVO, err := domain.NewDisplayName(name)
	if err != nil {
		nameVO = ""
	}

	picVO, err := domain.NewAvatarURL(picture)
	if err != nil {
		picVO = ""
	}

	return domain.IDPClaims{
		Sub:     subVO,
		Email:   emailVO,
		Name:    nameVO,
		Picture: picVO,
	}, nil
}

func bearerToken(headers map[string]string) string {
	h := headers["authorization"]
	if h == "" {
		h = headers["Authorization"]
	}

	// resolve "Bearer\nTOKEN" e outros whitespaces
	parts := strings.Fields(h)
	if len(parts) >= 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func decodeJwtPayload(token string) (jwtPayload, bool) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return jwtPayload{}, false
	}

	b, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return jwtPayload{}, false
	}

	var p jwtPayload
	if err := json.Unmarshal(b, &p); err != nil {
		return jwtPayload{}, false
	}

	return p, true
}

func joinName(given, family string) string {
	given = strings.TrimSpace(given)
	family = strings.TrimSpace(family)

	if given == "" {
		return family
	}
	if family == "" {
		return given
	}
	return given + " " + family
}
