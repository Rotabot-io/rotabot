PROG = bin/app
MODULE = github.com/rotabot-io/rotabot
GIT_SHA = $(shell git rev-parse --short HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_COMMAND = CGO_ENABLED=0 go build -ldflags "-X 'main.Sha=$(GIT_SHA)' -X 'main.Date=$(DATE)'"
LINT_COMMAND = golangci-lint run

.PHONY: clean
clean:
	rm -rvf $(PROG) $(PROG:%=%.linux_amd64) $(PROG:%=%.darwin_amd64)

.PHONY: build
build: clean $(PROG)

.PHONY: all darwin linux
all: darwin linux
darwin: $(PROG:=.darwin_amd64)
linux: $(PROG:=.linux_amd64)

bin/%.linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD_COMMAND) -a -o $@ cmd/$*/*.go

bin/%.darwin_amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(BUILD_COMMAND) -a -o $@ cmd/$*/*.go

bin/%:
	$(BUILD_COMMAND) -o $@ cmd/$*/*.go

.PHONY: generate
generate:
	go generate ./...
	sqlc generate --experimental
	goa gen $(MODULE)/design

.PHONY: test
test:
	go run github.com/onsi/ginkgo/v2/ginkgo -r -p \
		-randomize-all \
		-randomize-suites \
		-race \
		-trace \
		-poll-progress-after=10s \
		-poll-progress-interval=10s \
		--cover --coverprofile=cover.profile \
		--junit-report=report.xml

.PHONY: dev
dev: build
	$(PROG) --log-format=pretty --verbose=true serve --dev

.PHONY: lint
lint:
	$(LINT_COMMAND)

.PHONY: lint-fix
lint-fix:
	gofmt -s -w .
	$(LINT_COMMAND) --fix

.PHONY: install
install:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install goa.design/goa/v3/cmd/goa@v3
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
	go install github.com/rjeczalik/interfaces/cmd/interfacer@latest
	go install go.uber.org/mock/mockgen@latest

# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./assets/migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" down

## migrations/goto version=$1: migrate to a specific version number
.PHONY: migrations/goto
migrations/goto:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" goto ${version}

## migrations/force version=$1: force database migration
.PHONY: migrations/force
migrations/force:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" force ${version}

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" version
