[![Build Status](https://travis-ci.org/jasonrichardsmith/sentry.svg?branch=master)](https://travis-ci.org/jasonrichardsmith/sentry)
[![Coverage Status](https://coveralls.io/repos/github/jasonrichardsmith/sentry/badge.svg?branch=master)](https://coveralls.io/github/jasonrichardsmith/sentry?branch=master)
[![GoDoc](https://godoc.org/github.com/jasonrichardsmith/sentry?status.svg)](https://godoc.org/github.com/jasonrichardsmith/sentry)
[![Go Report Card](https://goreportcard.com/badge/github.com/jasonrichardsmith/sentry)](https://goreportcard.com/report/github.com/jasonrichardsmith/sentry)

# Sentry

Sentry is a Webhook Validating Admission Controller that enforces rules cluster wide on objects in Kubernetes prior to admission.

## Rules
 
Sentry currently supports the below enforcement rules.

If they are not set in the config.yaml with "enabled" set to true, they will not be enforced.

Each rule can ignore a set of namespaces.

To enforce different configurations you can launch this admission controller under different names with different configurations.

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

By listing the entire image path with or without tag, you can allow specific images from a repository. So "gcr.io/google_containers/pause-amd64" would only allow the pause container.  Due to the matching strategy this also means "gcr.io/google_containers/pause-amd64foo" would also pass.

```yaml
source:
  type: Pod
  enabled: true
  ignoredNamespaces:
    - "test2"
    - "test3"
  allowed:
    - "this/isallowed"
    - "sois/this"
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

To run the e2e tests you can run

```bash
make e2etests
```

## Development

To develop a new module, you can copy over the "example" module.

It will need to be added to the mux/config.go.

```go

type Config struct {
	Limits  limits.Config  `yaml:"limits"`
	Healthz healthz.Config `yaml:"healthz"`
	Source  source.Config  `yaml:"source"`
	Tags    tags.Config    `yaml:"tags"`
	Example example.Config `yaml:"example"`
}

func New() *Config {
	l := limits.Config{}
	h := healthz.Config{}
	i := tags.Config{}
	s := source.Config{}
	e := example.Config{}
	return &Config{
		Limits:  l,
		Healthz: h,
		Tags:    i,
		Source:  s,
		Example: e,
	}
}
```

and to the mux/mux.go NewFromConfig function

```go

	if c.Example.Enabled {
		log.Info("Example enabled loading")
		s, err := c.Example.LoadSentry()
		if err != nil {
			return sm, err
		}
		mod := sentryModule{
			s,
			c.Source.IgnoredNamespaces,
		}
		log.Info("Ignoring Namespaces ", mod.ignored)
		sm.Sentries = append(sm.Sentries, mod)
	}
```

Hopefully (config loading will be improved)[https://github.com/jasonrichardsmith/sentry/issues/2] in the near future.


You can add e2e tests by adding a folder for your module in test-manifests, and aadding manifests named in the following convention.

```
description.expectation.yaml
```

Anything not titled with "pass" as an "expectation" will be expected to fail.

Then make sure your module is enabled in the manifest.yaml.
```yaml

    example:
      enabled: true
      ignoredNamespaces:
        - "kube-system"
```
