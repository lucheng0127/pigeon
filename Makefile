GOCMD=go
GOGET=$(GOCMD) get

all: build

build:
	$(GOCMD) build -o pigeond pigeond/pigeond.go
	$(GOCMD) build -o pigeon pigeon/pigeon.go

clean:
	rm -f pigeond/pigeond
	rm -f pigeon/pigeon

deps:
	$(GOGET) "github.com/sirupsen/logrus"
	$(GOGET) "github.com/spf13/cobra