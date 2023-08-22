MODULE_NAME := myapp
SRC := $(shell find . -name '*.go')
BUILD_DIR := ./build
BINARY_NAME := $(BUILD_DIR)/$(MODULE_NAME)
SONAR_HOST_URL := https://sonarcloud.io
SONAR_SECRET := $(shell cat .sonar.secret)
BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)
CHANGE_TARGET := $(shell git rev-parse --abbrev-ref --symbolic-full-name @{u} | sed 's/.*\///')
# CHANGE_ID := $(shell git rev-parse --short=8 HEAD)
CHANGE_ID := $(shell whoami)

ifeq ($(OS), Windows_NT)
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	SHELL_VERSION = $(shell (Get-Host | Select-Object Version | Format-Table -HideTableHeaders | Out-String).Trim())
	OS = $(shell "{0} {1}" -f "windows", (Get-ComputerInfo -Property OsVersion, OsArchitecture | Format-Table -HideTableHeaders | Out-String).Trim())
	PACKAGE = $(shell (Get-Content go.mod -head 1).Split(" ")[1])
	CHECK_DIR_CMD = if (!(Test-Path $@)) { $$e = [char]27; Write-Error "$$e[31mDirectory $@ doesn't exist$${e}[0m" }
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
	RM_F_CMD = Remove-Item -erroraction silentlycontinue -Force
	RM_RF_CMD = ${RM_F_CMD} -Recurse
	SERVER_BIN = ${SERVER_DIR}.exe
	CLIENT_BIN = ${CLIENT_DIR}.exe
else
	SHELL := bash
	SHELL_VERSION = $(shell echo $$BASH_VERSION)
	UNAME := $(shell uname -s)
	VERSION_AND_ARCH = $(shell uname -rm)
	ifeq ($(UNAME),Darwin)
		OS = macos ${VERSION_AND_ARCH}
	else ifeq ($(UNAME),Linux)
		OS = linux ${VERSION_AND_ARCH}
	else
		$(error OS not supported by this Makefile)
	endif
	PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
	CHECK_DIR_CMD = test -d $@ || (echo "\033[31mDirectory $@ doesn't exist\033[0m" && false)
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	RM_F_CMD = rm -f
	RM_RF_CMD = ${RM_F_CMD} -r
	SERVER_BIN = ${SERVER_DIR}
	CLIENT_BIN = ${CLIENT_DIR}
endif

.DEFAULT_GOAL := help

.PHONY: all

build: ## Build the binary
	go build -v -ldflags="-X 'version.Version=v1.0.0' -X 'version.GitCommit=$(shell git rev-parse --short=8 HEAD)' -X 'build.User=$(shell id -u -n)' -X 'build.Time=$(shell date)'" -o $(BINARY_NAME)

clean: ## Clean build files
	go clean
	rm -f $(BINARY_NAME)

unit-test: ## Run unit tests
	@echo "Running unit tests"
	rm -f unit-test-coverage.out && \
	go test -v $(shell go list ./... | grep -v /test) \
	-count=1 \
	-cover \
	-coverpkg=./... \
	-coverprofile=./unit-test-coverage.out

api-test: ## Run api tests
	@echo "Running api tests"
	rm -f api-test-coverage.out && \
	go test -v ./test/ \
	-count=1 \
	-cover \
	-coverpkg=./... \
	-coverprofile=./interface-test-coverage.out

cov-html: ## Generate coverage html
	go tool cover -html=api-test-coverage.out -html=unit-test-coverage.out -o merged-coverage.html

sonarqube-pr: ## Run sonarqube analysis for pull request
	rm -rf .scannerwork && \
	sonar-scanner \
		-Dsonar.host.url="$(SONAR_HOST_URL)" \
		-Dsonar.working.directory=".scannerwork" \
		-Dsonar.pullrequest.key="$(CHANGE_ID)" \
		-Dsonar.pullrequest.branch="$(BRANCH_NAME)" \
		-Dsonar.pullrequest.base="$(CHANGE_TARGET)" \
		-Dsonar.login="$(SONAR_SECRET)"

sonarqube-branch: ## Run sonarqube analysis on branch
	rm -rf .scannerwork && \
	sonar-scanner \
		-Dsonar.host.url="$(SONAR_HOST_URL)" \
		-Dsonar.working.directory=".scannerwork" \
		-Dsonar.branch.name="$(BRANCH_NAME)" \
		-Dsonar.login="$(SONAR_SECRET)"

coverage: ## Run unit and api tests with coverage
	make unit-test && make api-test

vet: ## Run go vet
	go vet $(SRC)

lint: ## Run golint
	go get golang.org/x/lint/golint
	$(GOPATH)/bin/golint ./...


gen-proto: ## Generate protobuf files
	protoc -I proto/ proto/*.proto \
	--go_out=proto \
	--go-grpc_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative

bump: all ## Update packages version
	go get -u ./...

about: ## Display info related to the build
	@echo "OS: ${OS}"
	@echo "Shell: ${SHELL} ${SHELL_VERSION}"
	@echo "Protoc version: $(shell protoc --version)"
	@echo "Go version: $(shell go version)"
	@echo "Go package: ${PACKAGE}"
	@echo "Openssl version: $(shell openssl version)"

help: ## Show this help
	@${HELP_CMD}