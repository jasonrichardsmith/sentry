# build stage
FROM golang:alpine AS build-env
RUN apk update && apk add curl git
RUN mkdir -p /go/src/github.com/jasonrichardsmith/Sentry
WORKDIR /go/src/github.com/jasonrichardsmith/Sentry
COPY  . .
RUN curl https://glide.sh/get | sh
RUN glide install
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o sentrywebhook

FROM alpine:latest
COPY --from=build-env /go/src/github.com/jasonrichardsmith/Sentry/sentrywebhook .
ENTRYPOINT ["/sentrywebhook"]
