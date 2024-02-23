FROM golang:1.22-bookworm

WORKDIR /opt/rinha

COPY go.mod go.sum ./
RUN go mod download

COPY ./api ./api
RUN go build -o /usr/local/bin/api ./api/cmd

ENTRYPOINT ["api"]
