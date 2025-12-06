# Build stage
FROM golang:1.25.4-alpine AS build

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary (adjust main package path if needed)
RUN go build -o server .

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# Copy binary from build stage
COPY --from=build /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]
