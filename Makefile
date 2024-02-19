run:
	docker-compose -f deployments/docker-compose.yml up

build-dev:
	docker build -f build/development/Dockerfile . -t server

build-pkg:
	docker build -f build/package/Dockerfile .
