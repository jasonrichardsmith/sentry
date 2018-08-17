.DEFAULT_GOAL := build
VERSION="0.1"
REPO="jasonrichardsmith/sentry"

build:
	docker build --no-cache -t ${REPO}:${VERSION} .
	
minikube: minikubecontext build

minikubecontext:
	eval $(shell minikube docker-env)
push:
	docker push ${REPO}:${VERSION}
dep:
	glide install
test: dep
	go test ./...
goveralls: dep
	go test -coverprofile=coverage.out ./...
	${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=travis-ci

deployk8s:
	kubectl apply -f sentry-ns.yaml
	./gen-cert.sh
	./ca-bundle.sh
	kubectl apply -f manifest-ca.yaml
