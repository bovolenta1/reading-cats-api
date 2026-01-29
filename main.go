package main

import (
	"context"
	"log"
	"os"

	appReading "reading-cats-api/internal/application/reading"
	appUser "reading-cats-api/internal/application/user"
	"reading-cats-api/internal/config"
	"reading-cats-api/internal/infra/db"
	infraReading "reading-cats-api/internal/infra/reading"
	infraUser "reading-cats-api/internal/infra/user"
	"reading-cats-api/internal/presentation/httpapi"
	httpReading "reading-cats-api/internal/presentation/httpapi"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var router *httpapi.Router

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ctx := context.Background()
	cfg := config.Load()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// user/me
	userRepo := infraUser.NewPostgresRepository(pool)
	userUC := appUser.NewEnsureMeUseCase(userRepo)
	meHandler := httpapi.NewMeHandler(userUC)

	// reading/logs
	readingRepo := infraReading.NewPostgresRepository(pool)
	readingUC := appReading.NewRegisterReadingUseCase(readingRepo, userRepo, "America/Sao_Paulo")
	getReadingProgressUC := appReading.NewGetReadingProgressUseCase(readingRepo, userRepo, "America/Sao_Paulo")
	changeGoalUC := appReading.NewChangeGoalUseCase(readingRepo, userRepo, "America/Sao_Paulo")
	registerReadingHandler := httpReading.NewRegisterReadingHandler(readingUC)
	getReadingProgressHandler := httpReading.NewGetReadingProgressHandler(getReadingProgressUC)
	changeGoalHandler := httpReading.NewChangeGoalHandler(changeGoalUC)

	router = httpapi.NewRouter(meHandler, registerReadingHandler, getReadingProgressHandler, changeGoalHandler)
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return router.Route(ctx, event)
}

func main() {
	lambda.Start(handler)
}
