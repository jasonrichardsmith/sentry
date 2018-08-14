.DEFAULT_GOAL := build
VERSION="0.1"
REPO="jasonrichardsmith/sentry"

build:
	docker build --no-cache -t ${REPO}:${VERSION} .
	
minikube: minikubecontext build

minikubecontext:
	eval $(minikube docker-env)
push:
	docker push ${REPO}:${VERSION}
dep:
	glide install
test: dep
	go test ./...
goveralls: dep
	${GOPATH}/bin/goveralls -service=travis-ci
