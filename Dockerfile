FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /src/server /app/server
COPY flag.txt /app/flag.txt

EXPOSE 5001

CMD ["/app/server"]
