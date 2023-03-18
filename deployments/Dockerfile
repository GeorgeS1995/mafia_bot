FROM golang:1.19.4-bullseye as base
WORKDIR code
COPY .. .

FROM base as building
RUN go build ./cmd/mafia_bot/mafia_bot.go && ls | grep -xv "mafia_bot" | xargs rm -rf
