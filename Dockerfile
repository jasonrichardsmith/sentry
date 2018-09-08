# build stage
FROM golang:1.10-stretch AS build-env
RUN mkdir -p /go/src/github.com/jasonrichardsmith/sentry
WORKDIR /go/src/github.com/jasonrichardsmith/sentry
COPY  . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o sentrywebhook

FROM scratch
COPY --from=build-env /go/src/github.com/jasonrichardsmith/sentry/sentrywebhook .
ENTRYPOINT ["/sentrywebhook"]
