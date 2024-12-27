envConfigPath:=CONFIG_PATH=./config/local.yml
envDBPassword:=DB_PASSWORD=postgres	
envServerClientAddress:=SERVER_CLIENT_ADDRESS=87.228.37.67:8080
envServerClientAddressLocal:=SERVER_CLIENT_ADDRESS=localhost:8080
envDBTest:=DB_TEST=true

server:	
	@echo "Running server"
	$(envConfigPath) $(envDBPassword) go run ./cmd/server/
.PHONY: server

proto:
	@echo "Generating proto"
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	internal/proto/registry.proto
.PHONY: proto

migration:
	@echo "Running migration"
	go run ./cmd/migrator \
 		--storage-path="postgres://postgres:postgres@localhost:5432/gophkeeper?sslmode=disable" \
 		--migrations-path=./migrations
.PHONY: migration

client:
	@echo "Running client"
	$(envServerClientAddress) go run ./cmd/client/
.PHONY: client

client-local:
	@echo "Running client"
	$(envServerClientAddressLocal) go run ./cmd/client/
.PHONY: client-local

lint:
	@echo "Running lint"
	golangci-lint run \
		-c .golangci.yml \
		> ./golangci-lint/lint.log
.PHONY: lint

testdb:
	@echo "Running docker compose for tests database"
	docker compose -f docker-compose.test.yaml up -d --build
	@echo "Database is ready"
	@echo "Run server with test database"
	$(envDBTest) $(envConfigPath) go run ./cmd/server/
.PHONY: testdb

test-functional:
	@echo "Running functional tests"
	go test -v -count=1 ./tests/...
.PHONY: test-functional

test:
	@echo "Running unit tests"
	go test -count=1 ./internal/app/... ./internal/client/... ./internal/lib/...
.PHONY: test

build-client:
	@echo "Building client"
	./build.sh
.PHONY: build-client

build-db-image:
	@echo "Building db image"
	docker compose -f docker-compose.yaml up --force-recreate -d
.PHONY: build-db-image
