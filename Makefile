run:
	docker-compose -f deployments/docker-compose.yml up

build-dev:
	docker build -f build/development/Dockerfile . -t server

build-pkg:
	docker build -f build/package/Dockerfile .

test-coverage:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	go tool cover -func coverage/coverage.out