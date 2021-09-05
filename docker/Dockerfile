# build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN mkdir -p /go/bin/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...
# static files
COPY ./config /go/bin/app/config
COPY ./dataload /go/bin/app/dataload

# final stage
FROM alpine:latest
RUN apk --no-cache add -U ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app/cmd
LABEL Name=fimsprotoms Version=0.0.1
EXPOSE 8080
WORKDIR /app