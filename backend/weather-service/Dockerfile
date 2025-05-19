FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o weather-service ./cmd/main.go

FROM gcr.io/distroless/base
COPY --from=builder /app/weather-service /
CMD ["/weather-service"]