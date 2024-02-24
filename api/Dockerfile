FROM golang:1.22-bookworm

WORKDIR /opt/rinha

# need dumb-init otherwise fiber's prefork fails
# see: https://github.com/gofiber/fiber/issues/1036#issuecomment-841763449
RUN apt-get update && apt-get install dumb-init -y

COPY go.mod go.sum ./
RUN go mod download

COPY ./api ./api
RUN go build -o /usr/local/bin/api ./api/cmd

ENV GOMAXPROCS 3

ENTRYPOINT ["/usr/bin/dumb-init", "--", "api"]
