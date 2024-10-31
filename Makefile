envConfigPath:=CONFIG_PATH=./config/local.yml
envDBPassword:=DB_PASSWORD=postgres	
envServerClientAddress:=SERVER_CLIENT_ADDRESS=localhost:8080

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

lint:
	@echo "Running lint"
	golangci-lint run \
		-c .golangci.yml \
		> ./golangci-lint/lint.log
.PHONY: lint

testdb:
	@echo "Running docker compose for tests database"
	docker compose -f docker-compose.test.yaml up -d --build
.PHONY: testdb
