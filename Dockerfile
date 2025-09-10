FROM golang:1.25-alpine

COPY . /app
WORKDIR /app

RUN go build

CMD ["./backend"]

