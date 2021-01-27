FROM golang:1.10 AS build
WORKDIR /go/src
COPY http ./http
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o http_port .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/http_port ./
EXPOSE 8080/tcp
ENTRYPOINT ["./http_port"]
