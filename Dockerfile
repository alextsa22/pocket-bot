FROM golang:alpine AS builder

COPY . /github.com/alextsa22/pocket-bot
WORKDIR /github.com/alextsa22/pocket-bot

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go


FROM alpine:3.12

COPY --from=0 /github.com/alextsa22/pocket-bot/bin/bot .
COPY --from=0 /github.com/alextsa22/pocket-bot/configs configs/
COPY --from=0 /github.com/alextsa22/pocket-bot/.env .

EXPOSE 80

CMD ["./bot"]
