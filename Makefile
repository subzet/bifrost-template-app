ENV ?= local
-include .env.$(ENV)
export

dev:
	BIFROST_DEV=1 air -c .air.toml
	

build:
	go run github.com/3-lines-studio/bifrost/cmd/build@latest main.go
	go build -o ./tmp/app main.go

start: build
	./tmp/app

doctor:
	go run github.com/3-lines-studio/bifrost/cmd/doctor@latest .

ATLAS := atlas
LOCAL_DB := file:./data/app.db?_fk=1
PROD_DB  := libsql://$$TURSO_DB_URL?authToken=$$TURSO_AUTH_TOKEN

.PHONY: migrations-generate
migrations-generate:
	@if [ -z "$(name)" ]; then \
		echo "Error: provide a name, e.g. make migrations-generate name=initial_schema"; \
		exit 1; \
	fi
	atlas migrate diff $(name) \
		--env gorm \
		--dir file://migrations \
		--dev-url "$(LOCAL_DB)"

# Default: make migrations-generate name=initial or name=add_bio_field
# If name not set → atlas migrate diff --env gorm

.PHONY: migrations-apply-local
migrations-apply-local:
	$(ATLAS) migrate apply \
		--dir file://migrations \
		--url "$(LOCAL_DB)"

.PHONY: migrations-apply-prod
migrations-apply-prod:
	$(ATLAS) migrate apply \
		--dir file://migrations \
		--url $(PROD_DB)

# Optional: inspect what Atlas sees from your models
.PHONY: inspect-models
inspect-models:
	$(ATLAS) schema inspect --env gorm

# Optional: quick reset local DB
.PHONY: dev-reset
dev-reset:
	rm -f ./data/dev.db
	make migrations-apply-local