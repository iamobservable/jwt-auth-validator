FROM golang:1.24-alpine AS builder

LABEL org.opencontainers.image.source=https://github.com/iamobservable/jwt-auth-validator
LABEL org.opencontainers.image.description="This app provides signing validation of a JWT key. When used with an nginx proxy, this can provide verification the JWT was signed by a system a specific JWT_SECRET. If the JWT_SECRET was not used during the signing process, the key will be seen as invalid."
LABEL org.opencontainers.image.licenses="Apache License 2.0"

WORKDIR /app

COPY . .

RUN go mod vendor
RUN go build -o main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
