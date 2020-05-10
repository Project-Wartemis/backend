FROM golang:alpine

WORKDIR /go/src/github.com/Project-Wartemis/pw-backend

COPY . .

RUN scripts/build.sh

CMD scripts/run.sh