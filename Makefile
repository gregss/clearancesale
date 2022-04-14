.PHONY: up
up:
	docker-compose up -d --build

.PHONY: down
down:
	docker-compose down

.PHONY: test
test:
	go test -race -count 100 ./...

.PHONY: test-integration
test-integration:
	go test -race -count 100 ./... # todo

.PHONY: protoc
protoc:
	rm -rf ./gen
	mkdir -p ./gen
	buf generate

.PHONY: pb-lint
lint:
	buf-lint

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run


.PHONY: lint
lint: golangci-lint buf-lint

.PHONY: fix-lint
fix-lint:
	gci write .
	gofumpt -l -w .

.PHONY: buf-lint
buf-lint:
	buf lint

.PHONY: build-app
build-app:
	CGO_ENABLED=0 go build -o ./bin/ ./cmd/app/app.go

.PHONY: build-server
build-server:
	CGO_ENABLED=0 go build -o ./bin/ ./cmd/server/server.go

.PHONY: build-creator
build-creator:
	CGO_ENABLED=0 go build -o ./bin/ ./cmd/creator/creator.go

.PHONY: build
build: build-app build-server build-creator
