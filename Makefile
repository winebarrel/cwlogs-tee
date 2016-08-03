VERSION=$(shell git tag | tail -n 1)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
RUNTIME_GOPATH=$(GOPATH):$(shell pwd)
TEST=$(wildcard src/cwlogs_tee/*_test.go)
SRC=$(filter-out $(TEST), $(wildcard src/cwlogs_tee/*.go))

ifeq ($(GOOS),windows)
	BIN=cwlogs-tee.exe
else
	BIN=cwlogs-tee
endif

all: cwlogs-tee

go-get:
	go get github.com/aws/aws-sdk-go
	go get github.com/golang/mock/gomock
	go get github.com/stretchr/testify
	go get github.com/bluele/go-timecop

cwlogs-tee: go-get main.go $(SRC)
	GOPATH=$(RUNTIME_GOPATH) go build -o $(BIN)

test: go-get $(SRC) $(TEST)
	GOPATH=$(RUNTIME_GOPATH) go test -v $(TEST) $(SRC)

clean:
	rm -f cwlogs-tee cwlogs-tee.exe *.gz *.zip

package: clean cwlogs-tee
ifeq ($(GOOS),windows)
	zip cwlogs-tee-$(VERSION)-$(GOOS)-$(GOARCH).zip cwlogs-tee.exe
else
	gzip -c cwlogs-tee > cwlogs-tee-$(VERSION)-$(GOOS)-$(GOARCH).gz
endif

mock:
	go get github.com/golang/mock/mockgen
	mkdir -p src/mockaws
	mockgen -source $(GOPATH)/src/github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface/interface.go -destination src/mockaws/cloudwatchlogsmock.go -package mockaws
