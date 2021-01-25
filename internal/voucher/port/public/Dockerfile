FROM golang:1.10 AS build
WORKDIR /go/src
COPY http ./http
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o http .

FROM scratch AS runtime
COPY --from=build /go/src/http ./
EXPOSE 8080/tcp
ENTRYPOINT ["./http"]
