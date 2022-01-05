FROM golang:1.15 as build

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./...

CMD ["/go/src/app/server"]

FROM scratch

WORKDIR /var

COPY --from=build /go/src/app/server /var

ENTRYPOINT ["/var/server"]