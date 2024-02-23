FROM golang:1.22-bookworm

WORKDIR /opt/rinha/api

COPY . .
RUN go mod download
RUN go build -o api .

ENTRYPOINT ["./api"]
