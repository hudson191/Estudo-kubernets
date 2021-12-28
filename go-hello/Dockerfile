FROM golang:1.15 as build

WORKDIR /go/src/app

COPY . .

RUN go build -o server .

FROM scratch

COPY --from=build /go/src/app/server ./

CMD ["./server"]