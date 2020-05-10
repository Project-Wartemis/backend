FROM golang:alpine

WORKDIR /go/src/github.com/Project-Wartemis/pw-backend

COPY . .
ADD ./nginx.conf /etc/nginx/nginx.conf

RUN scripts/build.sh

CMD scripts/run.sh