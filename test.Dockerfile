FROM golang:1.16 as BUILDER

RUN mkdir /app
ADD . /app
WORKDIR /app

CMD go test -v ./...