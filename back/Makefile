TEST_COUNT ?= 1
TEST_ARGS ?=
COVERPROFILE ?= coverage.out
BIN_DIR = bin
export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

.PHONY: lint
lint: ## run linter
	${BIN_DIR}/golangci-lint --color=always run ./... -v --timeout 5m

golangci: ## install golangci-linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${BIN_DIR} v1.27.0

go-acc: ## install coverage tool
	go get github.com/ory/go-acc@v0.2.3

swag: ## install swaggo
	go get -u github.com/swaggo/swag/cmd/swag

install-deps: golangci go-acc swag ## install necessary dependencies

.PHONY: test
test:
	go test -v ./... -race -count $(TEST_COUNT) $(TEST_ARGS)

.PHONY: cover
cover:
	go-acc ./... && go tool cover -html=coverage.txt

.PHONY: test_e2e
test_e2e: stop start_rebuild
	cd ../front && yarn run cypress run


.PHONY: docker
docker:
	docker build -t seismographd .
	docker tag seismographd seismographd:$(tag)

.PHONY: build_docker
build_docker:
	docker build -t seismographd .

.PHONY: build
build:
	go build -o seismographd cmd/server/main.go

.PHONY: swagger
swagger:
	swag init --dir back/cmd/server --parseDependency

.PHONY: start
start:
	docker-compose up -d && cd ../front && yarn dev

.PHONY: start_rebuild
start_rebuild: build_docker start

.PHONY: stop
stop:
	docker-compose down --volumes
