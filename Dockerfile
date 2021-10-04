FROM golang:1.16 as BUILDER

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/apiserver/main.go

FROM alpine:latest AS prodaction
COPY --from=builder /app .
CMD ["./app"]