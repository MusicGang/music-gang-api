FROM golang:1.17 AS builder

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    GIT_TAG=$(git describe --tags --always) && \ 
    go build -ldflags "-X main.Commit=$GIT_COMMIT" -o mg-api cmd/main.go

FROM debian

WORKDIR /app

COPY --from=builder /app/mg-api .

CMD [ "./mg-api" ]