FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6 && swag init -g cmd/main.go -o cmd/docs
RUN CGO_ENABLED=0 go build -o /app/bin ./cmd/main.go

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 1001 app
COPY --from=builder /app/bin /app/bin
USER app
EXPOSE 8080
CMD ["/app/bin"]
