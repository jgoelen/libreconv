FROM golang:1.9-alpine3.6 AS builder
RUN apk add --no-cache libreoffice git
WORKDIR /go/src/github.com/jgoelen/libreconv
COPY . .
RUN go get ./...
RUN go test
RUN CGO_ENABLED=0 GOOS=linux go build -a -o libreconv

FROM alpine:3.6
RUN apk add --no-cache libreoffice msttcorefonts-installer && update-ms-fonts && fc-cache -f
WORKDIR /app
COPY --from=builder /go/src/github.com/jgoelen/libreconv /app/
ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ./libreconv