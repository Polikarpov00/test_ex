FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o /app ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
