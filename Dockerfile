FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o webapp .

EXPOSE 8080

CMD ["./webapp"]