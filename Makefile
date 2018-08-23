.DEFAULT_GOAL := build
VERSION="0.1"
REPO="jasonrichardsmith/sentry"
HASH=$(shell git log --pretty=format:'%H' -n 1)

build:
	docker build --no-cache -t ${REPO}:${VERSION} .
	
buildhash:
	docker build --no-cache -t ${REPO}:${HASH} .

minikube: minikubecontext build

minikubecontext:
	eval $(shell minikube docker-env)
push:
	docker push ${REPO}:${VERSION}
pushhash:
	docker push ${REPO}:${HASH}
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
deploydindk8s: deployk8s
	kubectl set image deployment/sentry -n sentry webhook=jasonrichardsmith/sentry:${HASH}
	kubectl rollout status -w -n sentry deployment/sentry
e2etests:
	cd test-manifests && ./e2etest.py
travise2e:
	./dind-cluster-v1.10.sh up 
	${MAKE} buildhash
	echo "${DOCKERPASSWORD}" | docker login -u "${DOCKERUSERNAME}" --password-stdin
	${MAKE} pushhash
	${MAKE} deploydindk8s
	${MAKE} e2etests
