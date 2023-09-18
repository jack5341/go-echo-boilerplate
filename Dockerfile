# Builder stage
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /bin/backend

# Final stage
FROM gcr.io/distroless/static:nonroot
WORKDIR /app/
COPY --from=builder /bin/backend /app/backend
ENTRYPOINT ["/app/backend"]
