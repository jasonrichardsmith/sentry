[![Build Status](https://travis-ci.org/jasonrichardsmith/sentry.svg?branch=master)](https://travis-ci.org/jasonrichardsmith/sentry)
[![Coverage Status](https://coveralls.io/repos/github/jasonrichardsmith/sentry/badge.svg?branch=master)](https://coveralls.io/github/jasonrichardsmith/sentry?branch=master)
[![GoDoc](https://godoc.org/github.com/jasonrichardsmith/sentry?status.svg)](https://godoc.org/github.com/jasonrichardsmith/sentry)
[![Go Report Card](https://goreportcard.com/badge/github.com/jasonrichardsmith/sentry)](https://goreportcard.com/report/github.com/jasonrichardsmith/sentry)

# Sentry

Sentry is a Webhook Admission Controller that enforces rules on objects in Kubernetes prior to admission.

## Rules

Sentry currently supports the following rules:

### Limits
Limits will insure all pods have limits for cpu and memory set.
