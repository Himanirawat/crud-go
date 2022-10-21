FROM golang:1.17.8-alpine3.15 AS builder


WORKDIR /app

COPY ./ ./
RUN apk add --update --no-cache ca-certificates git
RUN go mod download
RUN go build -o /crud



EXPOSE 8000
CMD [ "/crud" ]