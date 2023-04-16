.PHONY: start
start:
	docker-compose up -d --build --remove-orphans

.PHONY: stop
stop:
	docker-compose down

.PHONY: e2e
e2e:
	go clean -testcache && go test -p 1 ./... -v -tags e2e

.PHONY: docker-stop
docker-stop:
	docker stop $(shell docker ps -aq)
	docker rm $(shell docker ps -aq)