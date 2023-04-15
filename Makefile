.PHONY: start
start:
	docker-compose up -d --build --remove-orphans

.PHONY: stop
stop:
	docker-compose down

.PHONY: generate-docs
generate-docs:
	swag init -dir cmd/api --o docs/api --ot yaml --instanceName api --pd true

.PHONY: serve-docs
serve-docs:
	python3 -m http.server 9000 --directory docs/api

.PHONY: e2e
e2e:
	go clean -testcache && go test -p 1 ./... -v -tags e2e

.PHONY: tools
tools:
	go install github.com/swaggo/swag/cmd/swag@v1.8.6
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v.1.11.0
	docker pull quay.io/goswagger/swagger:v0.30.3
	alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger'

.PHONY: docker-stop
docker-stop:
	docker stop $(shell docker ps -aq)
	docker rm $(shell docker ps -aq)