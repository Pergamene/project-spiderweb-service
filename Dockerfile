FROM golang:alpine

RUN mkdir -p /go/src/github.com/Pergamene/project-spiderweb-service

WORKDIR /go/src/github.com/Pergamene/project-spiderweb-service

CMD ["go", "run", "cmd/server/main.go"]