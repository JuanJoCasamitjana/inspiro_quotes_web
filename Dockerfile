FROM golang:1.21.6-alpine

RUN apk update && apk add --no-cache gcc

ENV CGO_ENABLED=1 \
    CC=gcc 

WORKDIR /app

COPY . .

RUN go build -o app ./cmd/app/main.go

EXPOSE 8080

CMD ["./app", "-action", "init-run", "-port", ":8080"]
