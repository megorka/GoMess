FROM golang:1.23.6-alpine3.20 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .
RUN go mod download


# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /goproject ./cmd/authorization/main.go

EXPOSE 8080

# Run
CMD ["/goproject"]

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /goproject .

COPY config/config.yaml /app/config/config.yaml

EXPOSE 8080

CMD ["./goproject"]