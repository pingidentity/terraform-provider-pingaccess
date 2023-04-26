SHELL := /bin/bash

.PHONY: install generate fmt vet starttestcontainer removetestcontainer spincontainer clearstates kaboom testacc testacccomplete generateresource openlocalwebapi

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
# Wait for the instance to become ready
	sleep 1
	duration=0
	while (( duration < 240 )) && ! docker logs pingaccess_terraform_provider_container 2>&1 | grep -q "PingAccess is up"; \
	do \
	    duration=$$((duration+1)); \
		sleep 1; \
	done
# Fail if the container didn't become ready in time
	docker logs pingaccess_terraform_provider_container 2>&1 | grep -q "PingAccess is up" || \
		{ echo "PingAccess container did not become ready in time. Logs:"; docker logs pingaccess_terraform_provider_container; exit 1; }
		
removetestcontainer:
	docker rm -f pingaccess_terraform_provider_container
	
spincontainer: removetestcontainer starttestcontainer

testacc:
	PINGACCESS_PROVIDER_HTTPS_HOST=https://localhost:9000 \
	PINGACCESS_PROVIDER_USERNAME=administrator \
	PINGACCESS_PROVIDER_PASSWORD=2Access \
	TF_ACC=1 go test -timeout 10m -v ./... -p 1

testacccomplete: removetestcontainer starttestcontainer install testacc

clearstates:
	find . -name "*tfstate*" -delete
	
kaboom: clearstates spincontainer install

generateresource: spincontainer
	PINGACCESS_PROVIDER_HTTPS_HOST="https://localhost:9000" \
	PINGACCESS_GENERATED_RESOURCE=virtualhosts \
	PINGACCESS_PROVIDER_USERNAME=administrator \
	PINGACCESS_PROVIDER_PASSWORD=2Access \
	python3 scripts/generate_resource.py

openlocalwebapi:
	open "https://localhost:9000/pa-admin-api/v3/api-docs/"
