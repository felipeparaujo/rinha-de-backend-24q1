FROM golang:1.22-bookworm

WORKDIR /opt/rinha

# install gow and dlv first so we can start the watcher and debugger
RUN GOBIN=/usr/local/bin go install github.com/mitranim/gow@latest

COPY go.mod go.sum ./
RUN go mod download

COPY ./api ./api
# RUN go build -o /usr/local/bin/api ./api/cmd

ENTRYPOINT ["gow", "-s", "-r=false", "-w=/opt/rinha/api", "run", "./api/cmd", "--"]
