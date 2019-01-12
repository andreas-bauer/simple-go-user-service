# First stage: image to build Go application
FROM golang:1.11 as builder

WORKDIR /build
COPY . .

# Build application with flags for Alpine Linux
WORKDIR /build/cmd/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

# Second stage: image to run Go application
FROM alpine:latest

RUN mkdir /app
WORKDIR /app

# Pull the binary from the builder container
COPY --from=builder /build/cmd/user-service .

# Run the binary
ENTRYPOINT ["./user-service"]
