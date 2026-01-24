-include .env.local
export

.PHONY: build start clean migrate-create migrate-up migrate-down migrate-version

MIGRATIONS_DIR=migrations
MIGRATE=migrate

build:
	sam build

copy-env:
	powershell -NoProfile -Command "Copy-Item -Force .\.env.local .\.aws-sam\build\HelloFunction\.env.local"

start: build copy-env
	sam local start-api -p 3001 --debug

clean:
	@if exist .aws-sam rmdir /s /q .aws-sam

# uso: make migrate-create name=create_user
MIGRATIONS_DIR=migrations
MIGRATE=migrate
DATABASE_URL_MIGRATE := $(patsubst postgresql://%,postgres://%,$(DATABASE_URL))

migrate-create:
	@if "$(name)"=="" (echo use: make migrate-create name=create_user & exit /b 1)
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	@if "$(DATABASE_URL)"=="" (echo DATABASE_URL not set. Put it in .env.local & exit /b 1)
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL_MIGRATE)" up

migrate-down:
	@if "$(DATABASE_URL)"=="" (echo DATABASE_URL not set. Put it in .env.local & exit /b 1)
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL_MIGRATE)" down 1

migrate-version:
	@if "$(DATABASE_URL)"=="" (echo DATABASE_URL not set. Put it in .env.local & exit /b 1)
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL_MIGRATE)" version