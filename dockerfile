FROM golang:1.25

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd

EXPOSE 8080

CMD ["/app/main"]
