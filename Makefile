ifeq ($(CD_MODE),true)
	ENVFILE_BASE=.env.cd
else
	ENVFILE_BASE=.env.example
endif

.DEFAULT_GOAL := startDev

env:
	cp $(ENVFILE_BASE) .env
.PHONY: env

deps:
	docker-compose run --rm golang go mod tidy
.PHONY: deps

fmt:
	docker-compose run --rm golang go fmt ./...
.PHONY: fmt

vet:
	docker-compose run --rm golang shadow ./...
.PHONY: vet

build: .env
	docker-compose run --rm --entrypoint ./scripts/build-app.sh golang
.PHONY: build

test: .env
	docker-compose run --rm golang go test -v ./...
.PHONY: test

testWithIntegrations: .env
	$(MAKE) startTestDB
	$(MAKE) migrateDB
	docker-compose run --rm golang go test --tags=integration -v ./...
	$(MAKE) stopContainers
.PHONY: testWithIntegrations

startTestDB: .env
	docker-compose up -d postgres-test
	sleep 5
.PHONY: startTestDB

migrateDB: .env
	docker-compose run --rm -v $(PWD)/scripts/migrations/wallet:/flyway/sql:delegated -e SCHEMA=wallet flyway
.PHONY: migrateDB

stopContainers:
	docker-compose down --remove-orphans
.PHONY: stopContainers

# dev-specific targets/derivatives

# run this to ensure interface changes are reflected by mocks
genMocks: .env
	docker-compose run --rm mockery --all --output=./mocks/autogen --inpackage --keeptree
.PHONY: genMocks

runDev: .env
	docker-compose run --rm -p "8080:8080" golang go run main.go
.PHONY: runDev

startDev:
	$(MAKE) fmt
	$(MAKE) vet
	$(MAKE) build
	docker-compose up app-dev
.PHONY: startDev