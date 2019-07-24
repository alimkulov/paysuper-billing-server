FROM golang:1.11-alpine AS builder

RUN apk add bash ca-certificates git tzdata

WORKDIR /application

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/paysuper_billing_service .

FROM alpine:3.9
RUN apk update && apk add ca-certificates tzdata && rm -rf /var/cache/apk/*

WORKDIR /application

COPY --from=builder /application /application
ENTRYPOINT /application/bin/paysuper_billing_service
