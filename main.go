package main

import (
	"context"

	"reading-cats-api/internal/application/user"
	"reading-cats-api/internal/config"
	"reading-cats-api/internal/infra/db"
	infraUser "reading-cats-api/internal/infra/user"
	"reading-cats-api/internal/presentation/httpapi"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var router *httpapi.Router

func init() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	repo := infraUser.NewPostgresRepository(pool)
	uc := user.NewEnsureMeUseCase(repo)
	meHandler := httpapi.NewMeHandler(uc)

	router = httpapi.NewRouter(meHandler)
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return router.Route(ctx, event)
}

func main() {
	lambda.Start(handler)
}
