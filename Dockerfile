FROM golang:1.24.0-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY . .

RUN go work use ./src/api ./src/carwise ./src/infra ./src/docs
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o carwise-backend ./src/api

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

RUN update-ca-certificates

ENV TZ=UTC

COPY --from=builder /app/carwise-backend .
COPY --from=builder /app/.well-known ./well-known

ENV GIN_MODE=release

EXPOSE 4040

CMD ["./carwise-backend"]
