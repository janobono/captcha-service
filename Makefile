.PHONY: tools clean generate generate-openapi generate-proto build fmt test vet

default: build

tools:
	@echo "  >  Installing openapi generator"
	@npm install @openapitools/openapi-generator-cli -g
	@echo "  >  Installing sqlc"
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "  >  Installing grpc"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "  >  Installing modules"
	@go mod tidy

clean:
	@echo "  >  Cleaning build cache"
	@go clean ./...
	@rm -rf bin
	@rm -rf generated

generate-openapi:
	@echo "  > Generate openapi source files"
	openapi-generator-cli generate \
	--generator-name go-gin-server \
	--input-spec contract/openapi/captcha-service.yaml \
	--output generated/openapi-gen \
	--additional-properties=interfaceOnly=true,packageName=openapi,generateMetadata=false,generateGoMod=false &&\
	mkdir -p generated/openapi && cp -r generated/openapi-gen/go/* generated/openapi/ && rm -rf generated/openapi-gen

generate-proto:
	@echo "  > Generate proto source files"
	mkdir -p generated/proto && \
	for file in contract/proto/*.proto; do \
		protoc -I=contract/proto \
			--go_out=generated/proto --go_opt=paths=source_relative \
			--go-grpc_out=generated/proto --go-grpc_opt=paths=source_relative \
			"$$file"; \
	done

generate: generate-openapi generate-proto

build: clean generate
	go build -o bin/captcha-service ./cmd/captcha-service/main.go

fmt:
	@echo "  >  Formatting code"
	@go fmt ./...

test:
	@echo "  >  Executing unit tests"
	@go test -v ./...

vet:
	@echo "  >  Checking code with vet"
	@go vet ./...