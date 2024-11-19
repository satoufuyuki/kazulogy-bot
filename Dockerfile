FROM golang:1.22-alpine AS build-stage

WORKDIR /tmp/build

COPY . .

# Build the project
RUN go build .

FROM alpine:3

LABEL name "Kazulogy Bot"
LABEL maintainer "Satou Fuyuki <satoufuyuki@proton.me>"

WORKDIR /app

COPY --from=build-stage /tmp/build/kazulogy-bot /app/kazulogy-bot

CMD ["/app/kazulogy-bot"]