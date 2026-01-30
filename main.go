package main

import (
	"context"
	"log"
	"os"

	appGroup "reading-cats-api/internal/application/group"
	appReading "reading-cats-api/internal/application/reading"
	appSeason "reading-cats-api/internal/application/season"
	appUser "reading-cats-api/internal/application/user"
	"reading-cats-api/internal/config"
	"reading-cats-api/internal/infra/db"
	infraGroup "reading-cats-api/internal/infra/group"
	infraReading "reading-cats-api/internal/infra/reading"
	infraSeason "reading-cats-api/internal/infra/season"
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

	// group/create
	groupRepo := infraGroup.NewPostgresRepository(pool)
	createGroupUC := appGroup.NewCreateGroupUseCase(groupRepo, userRepo)
	createGroupHandler := httpapi.NewCreateGroupHandler(createGroupUC)

	// season/create
	seasonRepo := infraSeason.NewPostgresRepository(pool)
	createSeasonUC := appSeason.NewCreateSeasonUseCase(seasonRepo, userRepo)
	createSeasonHandler := httpapi.NewCreateSeasonHandler(createSeasonUC)

	router = httpapi.NewRouter(meHandler, registerReadingHandler, getReadingProgressHandler, changeGoalHandler, createGroupHandler, createSeasonHandler)
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return router.Route(ctx, event)
}

func main() {
	lambda.Start(handler)
}
