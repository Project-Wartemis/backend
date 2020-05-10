FROM golang:alpine

WORKDIR /go/src/github.com/Project-Wartemis/pw-backend

COPY . .
COPY ./nginx.conf /etc/nginx/

RUN scripts/build.sh

CMD scripts/run.sh