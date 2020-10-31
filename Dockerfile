FROM golang:1.15

LABEL mantainers="Andrea Zanin <andrea@igloo.ooo>"
WORKDIR /usr/src/app
COPY . .
RUN go get

CMD ["go","run","."]