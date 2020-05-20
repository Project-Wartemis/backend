FROM golang:alpine

WORKDIR /go/src/github.com/Project-Wartemis/pw-backend

COPY . .

CMD ./backend