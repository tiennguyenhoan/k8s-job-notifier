# Builder
FROM golang:alpine as builder

WORKDIR /go/cronjob-fail-notifier
ADD . /go/cronjob-fail-notifier

RUN go build -o /app .

# Runner 
FROM alpine

COPY --from=builder /app /app
ENTRYPOINT ["/app"]
