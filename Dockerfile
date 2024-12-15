FROM golang:1.23

WORKDIR /usr/src/user-registration-service

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/user-registration-service ./cmd/user-registration-service/...

CMD user-registration-service