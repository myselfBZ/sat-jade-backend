FROM golang:1.24.1 AS builder

WORKDIR /app

# copy go mod files and download deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# runtime stage
FROM alpine:3.19

WORKDIR /app

# copy binary from builder
COPY --from=builder /app/api .

# set default port (change if needed)
EXPOSE 8080

# run the binary
CMD ["./api"]