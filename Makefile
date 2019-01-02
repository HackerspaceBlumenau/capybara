DOCKER_VENDOR=hsbnu
DOCKER_PROJECT=capybara
DOCKER_IMAGE="${DOCKER_VENDOR}/${DOCKER_PROJECT}"

default: build

build:
	docker build -t ${DOCKER_IMAGE} .

.PHONY: publish
push: build
	docker push ${DOCKER_IMAGE}

.PHONY: test
test:
	docker run -it --rm -v ${PWD}:/go/src/github.com/hackerspaceblumenau/capybara -w /go/src/github.com/hackerspaceblumenau/capybara golang go test -failfast ./...

.PHONY: ci
ci: test build push
