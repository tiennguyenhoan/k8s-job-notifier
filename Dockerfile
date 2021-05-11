FROM golang:alpine as builder

WORKDIR /go/cronjob-fail-notifier
ADD . /go/cronjob-fail-notifier

RUN go build -o /app .

FROM alpine

COPY --from=builder /app /app
ENTRYPOINT ["/app"]
