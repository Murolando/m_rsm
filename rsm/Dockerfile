FROM golang:1.23-alpine

RUN go version
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY ./ ./

RUN go build -o cmd/app ./cmd/main.go
EXPOSE 8080
CMD ["/app/cmd/app", "--app-cfg=./configs/test/config.yaml"]