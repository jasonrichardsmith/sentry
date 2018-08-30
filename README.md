[![Build Status](https://travis-ci.org/jasonrichardsmith/sentry.svg?branch=master)](https://travis-ci.org/jasonrichardsmith/sentry)
[![Coverage Status](https://coveralls.io/repos/github/jasonrichardsmith/sentry/badge.svg?branch=master)](https://coveralls.io/github/jasonrichardsmith/sentry?branch=master)
[![GoDoc](https://godoc.org/github.com/jasonrichardsmith/sentry?status.svg)](https://godoc.org/github.com/jasonrichardsmith/sentry)
[![Go Report Card](https://goreportcard.com/badge/github.com/jasonrichardsmith/sentry)](https://goreportcard.com/report/github.com/jasonrichardsmith/sentry)

# Sentry

Sentry is a Webhook Validating Admission Controller that enforces rules cluster wide on objects in Kubernetes prior to admission.

This project is in Beta Release

## Rules
 
Sentry currently supports the below enforcement rules.

If they are not set in the config.yaml with "enabled" set to true, they will not be enforced.

Each can ignore a set of namespaces.

"type" targets the Kuberentes object type.


### Limits
 
Limits will insure all pods have limits for cpu and memory set and are within the range you provide.

```yaml
limits:
  enabled: true
  ignoredNamespaces:
    - "test2"
    - "test3"
  cpu:
    min: "1"
    max: "2"
  memory:
    min: 1G
    max: 2G
```
 
### Source

Source insures images are only pulled from allowed sources.  This is a very simple string match.  This will only check if your image string starts with strings provided in the config.  To insure your domain is not read as a subdomain, it is best to end your domain with a "/".

By listing the entire image path with or without tag, you can allow specific images from a repository. So "gcr.io/google_containers/pause-amd64" would only allow the pause container.

```yaml
source:
  type: Pod
  enabled: true
  ignoredNamespaces:
    - "test2"
    - "test3"
  allowed:
    - "thisdomain/isallowed"
    - "sois/thisone"
```


### Healthz
 
Healthz insures liveliness and readiness probes are set.

```yaml
healthz:
  enabled: true
  ignoredNamespaces:
    - "test1"
    - "test3"
```

 
### Tags

Tags insures no containers launch with 'latest' or with no tag set.

```yaml
tags:
  enabled: true
  ignoredNamespaces:
    - "test1"
    - "test2"
```


## Try out sentry
 
To build and test in minikube you can run

```bash
$ minikube start --kubernetes-version v1.11.1
$ make minikube
```

Please use Kubernetes version >= 1.1.0

This will build a container from source on your minikube server.

You can deploy by running:

```bash
$ make deployk8s
```

This create server certs, and makes them available in the deployment. It produces a manifest-ca.yaml which gets deployed.

To see the tests working you can deploy any of the manifests under the test-manifests folder.

