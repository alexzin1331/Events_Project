FROM golang:1.24-alpine AS builder

WORKDIR /app

# Контекст сборки - директория backend, поэтому просто копируем все
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main cmd/url/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml . 
EXPOSE 50051
CMD ["./main"]