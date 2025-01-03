include .env
export

envConfigPath:=CONFIG_PATH=./config/local.yml
envDBPassword:=DB_PASSWORD=${DB_PASSWORD}	
envServerClientAddress:=SERVER_CLIENT_ADDRESS=87.228.37.67:8080
envServerClientAddressLocal:=SERVER_CLIENT_ADDRESS=localhost:8080
envS3SecretKey:=S3_SECRET_KEY=${S3_SECRET_KEY}
envDBTest:=DB_TEST=true
envDBPort:=DB_PORT=${DB_PORT}	
envDBUser:=DB_USER=${DB_USER}
server:	
	@echo "Running server"
	$(envConfigPath) $(envDBPassword) $(envS3SecretKey) $(envDBPort) $(envDBUser) go run ./cmd/server/
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
	$(envDBTest) $(envConfigPath) $(envS3SecretKey) go run ./cmd/server/
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

build-docker-container:
	@echo "Building db and app image"
	docker compose down || true # удаляет контейнеры и сети если они существуют
	docker compose -f docker-compose.yaml up --force-recreate -d  # создает и запускает контейнеры и сети
.PHONY: build-docker-container
