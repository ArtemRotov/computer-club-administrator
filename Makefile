.DEFAULT_GOAL := run

.PHONY: run
run:
	go run ./cmd/app/main.go tests/test.txt
 
.PHONY: build
build:
	go build -v -o app_service ./cmd/app/main.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: build-docker
build-docker:
	docker build --rm -t computer-club-manager . ; docker image prune -f

.PHONY: run-docker
run-docker:
	docker run --rm -i computer-club-manager /app/tests/test.txt

#docker run --rm -i -v /home/user/github.com/computer-club-manager/input.txt:/app/input.txt computer-club-manager /app/input.txt
