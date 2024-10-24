envConfigPath:=CONFIG_PATH=./config/local.yml
envDBPassword:=DB_PASSWORD=postgres	

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