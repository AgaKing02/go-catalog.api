FROM golang

ADD . /go/src/

WORKDIR /app

COPY . /app

EXPOSE 9090

RUN go build -o app main.go

ENTRYPOINT  /app/app