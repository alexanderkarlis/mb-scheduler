FROM golang:1.14 AS build

WORKDIR /usr/src/app/server
RUN mkdir -p /usr/src/app/server
RUN mkdir -p /usr/local/bin/
COPY . /usr/src/app/server

RUN go mod download
RUN GOOS=linux go build -o backend main.go

# COPY /usr/src/app/server/backend /usr/local/bin/backend

ENV HOST=0.0.0.0

CMD ["./backend"]