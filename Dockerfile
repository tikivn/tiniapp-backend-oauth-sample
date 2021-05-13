# Build
FROM golang:1.16-alpine AS builder

ARG CGO_ENABLED=0
ARG GO111MODULE=on
ARG GOARCH=amd64
ARG GOOS=linux

RUN apk add --update --no-cache ca-certificates curl git tzdata
RUN ln -fs /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime

WORKDIR /repo

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN /repo/scripts/bin.sh build

# Deploy
FROM alpine:3.13 as deployer

RUN apk add --update --no-cache ca-certificates curl tzdata
RUN ln -fs /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime

WORKDIR /repo

COPY --from=builder /repo/_build /repo/
RUN chmod +x /repo/run.sh

CMD ["/repo/run.sh"]
