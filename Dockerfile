# build stage
FROM golang:alpine AS build-stage
RUN apk update
WORKDIR /go/src/app
COPY . .
RUN mkdir -p /go/bin/app
RUN go build -o /go/bin/app -v ./...
# static files
COPY ./config /go/bin/app/config
COPY ./dataload /go/bin/app/dataload
COPY ./docs /go/bin/app/docs
COPY ./i18n /go/bin/app/i18n

# production stage
FROM alpine:latest AS production-stage
RUN apk update
COPY --from=build-stage /go/bin/app /app
ENTRYPOINT /app/cmd
WORKDIR /app