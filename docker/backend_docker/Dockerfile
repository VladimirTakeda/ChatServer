FROM golang:1.21.1-alpine

RUN go version
ENV GOPATH=/

# install psql
RUN apk update
RUN apk add postgresql-client

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY pkg/go.mod pkg/go.sum ./
#RUN go mod download && go mod verify

COPY ../../go.mod ./

RUN go mod download
RUN go mod tidy

COPY ../.. ./
RUN go build -v -o chatServer ./cmd/main.go

CMD ["./chatServer"]