VERSION=$(shell git tag | tail -n 1)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
RUNTIME_GOPATH=$(GOPATH):$(shell pwd)

ifeq ($(GOOS),windows)
	BIN=cwlogs-tee.exe
else
	BIN=cwlogs-tee
endif

all: cwlogs-tee

go-get:
	go get github.com/aws/aws-sdk-go

cwlogs-tee: go-get main.go $(SRC)
	GOPATH=$(RUNTIME_GOPATH) go build -o $(BIN)

clean:
	rm -f cwlogs-tee cwlogs-tee.exe *.gz *.zip

package: clean cwlogs-tee
ifeq ($(GOOS),windows)
	zip cwlogs-tee-$(VERSION)-$(GOOS)-$(GOARCH).zip cwlogs-tee.exe
else
	gzip -c cwlogs-tee > cwlogs-tee-$(VERSION)-$(GOOS)-$(GOARCH).gz
endif
