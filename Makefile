GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=irishub-sync
BINARY_UNIX=$(BINARY_NAME)-unix

all: get_tools get_deps build

get_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

######################################
## Tools

check_tools:
	cd tools && $(MAKE) check_tools

get_tools:
	cd tools && $(MAKE) get_tools

update_tools:
	cd tools && $(MAKE) update_tools