FROM golang:1.26.1

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .


CMD ["make", "run"]