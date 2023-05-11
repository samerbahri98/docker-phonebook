FROM golang:1.20 as build

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build .

EXPOSE 3333

ENTRYPOINT ["/app/docker-phonebook"]
