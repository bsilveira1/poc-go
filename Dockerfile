# Build stage
FROM golang:1.23.4 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o golang-app .

# Run stage
FROM golang:1.18

WORKDIR /app

COPY --from=builder /app /app

CMD ["./golang-app"]
