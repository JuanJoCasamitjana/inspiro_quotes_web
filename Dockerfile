FROM golang:latest

RUN apt-get update && apt-get install -y git build-essential

WORKDIR /app

RUN git clone https://github.com/JuanJoCasamitjana/inspiro_quotes_web.git

RUN go env -w CGO_ENABLED=1
RUN go build -o app ./cmd/app/main.go

EXPOSE 8080

CMD ["./app -action init-run -port :8080"]