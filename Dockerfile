# Stage1: OpenAPI code generation
FROM public.ecr.aws/docker/library/node:lts-alpine AS openapi

RUN apk add openjdk17-jre && npm install @openapitools/openapi-generator-cli -g && apk add --no-cache font-dejavu

RUN find /usr/share/fonts -type f

WORKDIR /src
COPY contract/openapi ./contract/openapi

RUN openapi-generator-cli generate \
    --generator-name go-gin-server \
    --input-spec contract/openapi/captcha-service.yaml \
    --output generated/openapi-gen \
    --additional-properties=interfaceOnly=true,packageName=openapi,generateMetadata=false,generateGoMod=false &&\
    mkdir -p generated/openapi && cp -r generated/openapi-gen/go/* generated/openapi/ && rm -rf generated/openapi-gen

# Stage2: gRPC code generation
FROM public.ecr.aws/docker/library/golang:alpine AS proto

RUN apk update && apk add --no-cache make protobuf-dev

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src
COPY contract/proto ./contract/proto

RUN mkdir -p generated/proto && \
    for file in contract/proto/*.proto; do \
        protoc -I=contract/proto \
            --go_out=generated/proto --go_opt=paths=source_relative \
            --go-grpc_out=generated/proto --go-grpc_opt=paths=source_relative \
            "$file"; \
    done

# Stage3: build
FROM public.ecr.aws/docker/library/golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=openapi /src/generated/openapi ./generated/openapi
COPY --from=proto /src/generated/proto ./generated/proto

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/captcha-service ./cmd/captcha-service

# Stage4: Final image
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/captcha-service .
COPY --from=openapi /usr/share/fonts/dejavu/DejaVuSans-Bold.ttf /usr/share/fonts/dejavu/DejaVuSans-Bold.ttf

# Default command
ENTRYPOINT ["/app/captcha-service"]