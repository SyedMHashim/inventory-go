FROM golang:1.19-alpine

RUN mkdir /inventory-go
WORKDIR /inventory-go

ADD go.* /inventory-go/
RUN go mod download

ADD . /inventory-go/
RUN go build -o inventory ./cmd/main.go
ENTRYPOINT ["./inventory"]