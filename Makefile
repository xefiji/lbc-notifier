build:
	cp -R docs deployment/dist
	go build -o deployment/dist/lbc cmd/lbc.go

build-linux:
	cp -R docs deployment/dist
	GOOS=linux go build -o deployment/dist/lbc cmd/lbc.go

run: build
	deployment/dist/lbc

dist: build-linux
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
	docker build -t xefiji/lbc deployment; \
	docker tag xefiji/lbc xefiji/lbc:$$VERSION; \

lint:
	golangci-lint run