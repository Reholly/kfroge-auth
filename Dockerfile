FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gettext

COPY go.* ./
RUN go mod download

COPY ./ ./

RUN go build -o ./bin/sso-service cmd/app/main.go

FROM alpine AS runner

COPY --from=builder /app/bin/sso-service /
COPY --from=builder /app/migrator/migrations /migrator/migrations
COPY --from=builder /app/config.yaml config.yaml

EXPOSE "8090:8090"

ENTRYPOINT ["./sso-service", "-config=./config.yaml"]