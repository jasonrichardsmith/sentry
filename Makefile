.DEFAULT_GOAL := build
VERSION="0.1"
REPO="jasonrichardsmith/Sentry"

build:
	docker build --no-cache -t jasonrichardsmith/Sentry:${VERSION} .
	
minikube: minikubecontext build

minikubecontext:
	eval $(minikube docker-env)
push:
	docker push ${REPO}:${VERSION}
