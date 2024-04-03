FROM golang:1.22.1-bullseye

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN export CGO_ENABLED=1 &&  go build -o ./app

EXPOSE 4000

CMD ["./app"]
