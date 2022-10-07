build:
	cp -R docs deployment/dist
	go build -o deployment/dist/lbc cmd/lbc/lbc.go
	go build -o deployment/dist/server cmd/server/server.go

build-linux:
	cp -R docs deployment/dist
	GOOS=linux go build -o deployment/dist/lbc cmd/lbc/lbc.go
	GOOS=linux go build -o deployment/dist/server cmd/server/server.go

run-lbc: build
	deployment/dist/lbc

run-server: build
	deployment/dist/server

dist-%: build-linux
	@if [ -z "$$VERSION" ]; then \
	echo "VERSION must be set"; \
	exit 2; \
	fi; \
	tag=$$(echo $$VERSION | sed -n "s/refs\/tags\///p"); \
	if [ $$tag ]; then \
		VERSION=$$tag; \
	else \
		echo "Only tag can be pushed to registry"; \
		exit 0; \
	fi; \
	docker build -t xefiji/$* -f deployment/.$*.Dockerfile deployment; \
	docker tag xefiji/$* xefiji/$*:$$VERSION; \

lint:
	golangci-lint run
