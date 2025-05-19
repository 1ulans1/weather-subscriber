FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o email-service ./cmd/main.go

FROM gcr.io/distroless/base
COPY --from=builder /app/email-service /
CMD ["/email-service"]