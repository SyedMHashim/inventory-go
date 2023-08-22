FROM golang:1.19-alpine

RUN mkdir /app
WORKDIR /app

ADD go.* /app/
RUN go mod download

ADD . /app/
RUN go build -o inventory ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/inventory ./
CMD ["./inventory"]
