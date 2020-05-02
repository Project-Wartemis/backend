FROM golang:1.14.2

WORKDIR /go/src/github.com/Project-Wartemis/pw-backend

COPY . .

RUN scripts/build.sh

CMD scripts/run.sh