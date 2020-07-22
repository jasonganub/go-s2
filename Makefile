GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

test:
	$(GOTEST) ./... -v

build:
	$(GOBUILD)

clean:
	$(GOCLEAN)