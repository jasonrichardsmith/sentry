.PHONY: build minikube minikubecontext push test goveralls deployk8s deploydindk8s e2etests e2eclean travise2e dindup buildpushhash
SHELL=/bin/bash -eo pipefail
.DEFAULT_GOAL := build
REPO=jasonrichardsmith/sentry
VERSION=$(shell cat VERSION)
HASH=$(shell git log --pretty=format:'%H' -n 1)
TAG=${VERSION}
build:
	docker build --no-cache -t ${REPO}:${TAG} .

push:
	docker push ${REPO}:${TAG}

minikube: minikubecontext build

minikubecontext:
	eval $(shell minikube docker-env)
test:
	go test ./...
goveralls:
	go test -coverprofile=coverage.out ./...
	${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=travis-ci

deployk8s:
	$(eval export TAG)
	kubectl apply -f sentry-ns.yaml
	./gen-cert.sh
	./ca-bundle.sh
	kubectl apply -f manifest-ca.yaml

deploydindk8s: hashtag deployk8s
	kubectl rollout status -w -n sentry deployment/sentry

e2etests:
	cd test-manifests && ./e2etest.py
e2eclean:
	cd test-manifests && ./e2eclean.py

travise2e: | dindup deploydindk8s e2etests

dindup:
	./dind-cluster-v1.10.sh up 

hashtag:
	$(eval export TAG=${HASH})

buildpushhash: | hashtag build push
