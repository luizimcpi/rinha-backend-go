FROM golang:1.21.5

WORKDIR /app

COPY . .

EXPOSE 8080

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

ENTRYPOINT [ "/app/server" ]