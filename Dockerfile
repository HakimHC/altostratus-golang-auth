FROM golang:1.22 as builder

WORKDIR /app
COPY go.* /app/
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/auth_service -ldflags "-X main.build=." ./cmd

FROM alpine:3.18

COPY --from=builder /app/bin/auth_service /usr/bin/auth_service
RUN chmod +x /usr/bin/auth_service

EXPOSE 80

CMD ["/usr/bin/auth_service"]

