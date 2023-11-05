# Based on the Example from Joel Homes, author of "Shipping Go" at
# https://github.com/holmes89/hello-api/blob/main/ch10/Makefile

SHELL=/bin/bash

GO_VERSION := 1.21  # <1>

COVERAGE_AMT := 60  # should be 80

test:
	go test ./...

build:
	go test ./... && echo "---ok---" && go build 

clean:
	rm cover.html coverage.out cover.rpt

check-format: 
	test -z $$(go fmt ./...)

check-vet: 
	test -z $$(go vet ./...)

install-lint:
	# https://golangci-lint.run/usage/install/#local-installation to GOPATH
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(HEREGOPATH)/bin v1.54.2
	# report version
	golangci-lint --version

lint:
	golangci-lint run ./... 

module-update-tidy:
	go get -u ./...
	go mod tidy

check: test check-format check-vet lint

