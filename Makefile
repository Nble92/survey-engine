PROJECT_NAME := "survey-engine"
PKG := "github.com/dchote/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep lint vet test test-coverage build clean setup

all: build

setup: ## Install required packages and setup environment
	@echo "Updating package list"
	@sudo apt update
	@echo "Installing Go and Node.js"
	@sudo apt install -y golang npm
	@echo "Installing Yarn"
	@sudo npm install -g yarn
	@echo "Creating data directory and SQLite database"
	@mkdir -p data
	@touch data/survey.sqlite

dep: setup ## Get the dependencies
	@echo "Installing dependencies"
	@go mod download
	@yarn --cwd ./frontend install

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

build: dep ## Build the binary file
	@echo "Building Frontend code"
	@yarn --cwd ./frontend build
	@echo "Building rice-box.go of Frontend assets"
	@rice embed-go
	@echo "Building survey-engine binary"
	@go build -i -o build/survey-engine $(PKG)

copy: ##create three more of this folder for the other variances of the app
	@cp 

crossbuild: build
	@echo "Building survey-engine Linux binary"
	@GOOS=linux GOARCH=amd64 go build -o build/survey-engine-linux-amd64 $(PKG)

clean: ## Remove build related file
	@rm -f cover.out coverage.txt
	@rm -rf build


