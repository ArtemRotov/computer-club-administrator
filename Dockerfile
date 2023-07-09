#1: Modules caching
FROM golang:1.20-alpine as modules
COPY go.mod /modules/
WORKDIR /modules
RUN go mod download

#2: Builder
FROM golang:1.20-alpine as builder
COPY . /app
WORKDIR /app
RUN GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app/main.go

# Step 3: Final
FROM golang:1.20-alpine
COPY --from=builder /bin/app /app/app
COPY --from=builder /app/tests /app/tests
ENTRYPOINT ["/app/app"]