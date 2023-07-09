.DEFAULT_GOAL := run

.PHONY: run
run:
	go run ./cmd/app/main.go input.txt
 
.PHONY: build
build:
	go build -v -o app_service ./cmd/app/main.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: build-docker
build-docker:
	docker build --rm -t computer-club-administrator-service . ; docker image prune -f

.PHONY: run-docker
run-docker:
	docker run --rm -i computer-club-administrator-service /app/tests/test_base.txt

#docker run --rm -i -v /home/user/github.com/computer-club-administrator/input.txt:/app/input.txt computer-club-administrator-service /app/input.txt
