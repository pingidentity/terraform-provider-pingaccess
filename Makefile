SHELL := /bin/bash

.PHONY: install generate fmt vet starttestcontainer

default: install

install:
	go mod tidy
	go install .

generate:
	go generate ./...
	go fmt ./...
	go vet ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

starttestcontainer:
	docker run --name pingaccess_terraform_provider_container \
		-d -p 3000:3000 \
		-d -p 9000:9000 \
		--env-file "${HOME}/.pingidentity/config" \
		pingidentity/pingaccess:2211-7.1.3