SHELL          := /bin/bash
PROGRAM        := cwlogs-tee
VERSION        := v0.1.4
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
TEST_SRC       := $(wildcard src/cwlogs_tee/*_test.go)
SRC            := $(filter-out $(TEST_SRC),$(wildcard src/cwlogs_tee/*.go))

UBUNTU_IMAGE          := docker-go-pkg-build-ubuntu-trusty
UBUNTU_CONTAINER_NAME := docker-go-pkg-build-ubuntu-trusty-$(shell date +%s)
CENTOS_IMAGE          := docker-go-pkg-build-centos6
CENTOS_CONTAINER_NAME := docker-go-pkg-build-centos6-$(shell date +%s)

.PHONY: all
all: $(PROGRAM)

.PHONY: go-get
go-get:
	go get github.com/aws/aws-sdk-go
	go get github.com/golang/mock/gomock
	go get github.com/stretchr/testify
	go get github.com/bluele/go-timecop
	go get github.com/cenkalti/backoff

$(PROGRAM): $(SRC)
ifeq ($(GOOS),linux)
	GOPATH=$(RUNTIME_GOPATH) go build -ldflags "-X cwlogs_tee.version=$(VERSION)" -a -tags netgo -installsuffix netgo -o $(PROGRAM)
else
	GOPATH=$(RUNTIME_GOPATH) go build -ldflags "-X cwlogs_tee.version=$(VERSION)" -o $(PROGRAM)
endif

.PHONY: mock
mock:
	go get github.com/golang/mock/mockgen
	mkdir -p src/mockaws
	mockgen -source $(GOPATH)/src/github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface/interface.go \
	  -destination src/mockaws/cloudwatchlogsmock.go \
	  -package mockaws

.PHONY: test
test: $(SRC) $(TEST_SRC)
	GOPATH=$(RUNTIME_GOPATH) go test -v $(TEST_SRC) $(SRC)

.PHONY: clean
clean:
	rm -f $(PROGRAM)
	rm -f pkg/*
	rm -f debian/$(PROGRAM).debhelper.log
	rm -f debian/$(PROGRAM).substvars
	rm -f debian/files
	rm -rf debian/$(PROGRAM)/

.PHONY: package
package: clean test $(PROGRAM)
	gzip -c $(PROGRAM) > pkg/$(PROGRAM)-$(VERSION)-$(GOOS)-$(GOARCH).gz
	rm -f $(PROGRAM)

.PHONY: package/linux
package/linux:
	docker run \
	  --name $(UBUNTU_CONTAINER_NAME) \
	  -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) \
	  make -C /tmp/src go-get package
	docker rm $(UBUNTU_CONTAINER_NAME)

.PHONY: deb
deb:
	docker run \
	--name $(UBUNTU_CONTAINER_NAME) \
	-v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) \
	make -C /tmp/src go-get deb/docker
	docker rm $(UBUNTU_CONTAINER_NAME)

.PHONY: deb/docker
deb/docker: clean
	dpkg-buildpackage -us -uc
	mv ../$(PROGRAM)_* pkg/

.PHONY: rpm
rpm:
	docker run \
	  --name $(CENTOS_CONTAINER_NAME) \
	  -v $(shell pwd):/tmp/src $(CENTOS_IMAGE) \
	  make -C /tmp/src go-get rpm/docker
	docker rm $(CENTOS_CONTAINER_NAME)

.PHONY: rpm/docker
rpm/docker: clean
	cd ../ && tar zcf $(PROGRAM).tar.gz src
	mv ../$(PROGRAM).tar.gz /root/rpmbuild/SOURCES/
	cp $(PROGRAM).spec /root/rpmbuild/SPECS/
	rpmbuild -ba /root/rpmbuild/SPECS/$(PROGRAM).spec
	mv /root/rpmbuild/RPMS/x86_64/$(PROGRAM)-*.rpm pkg/
	mv /root/rpmbuild/SRPMS/$(PROGRAM)-*.src.rpm pkg/

.PHONY: docker/build/ubuntu-trusty
docker/build/ubuntu-trusty:
	docker build -f etc/Dockerfile.ubuntu-trusty -t $(UBUNTU_IMAGE) .

.PHONY: docker/build/centos6
docker/build/centos6:
	docker build -f etc/Dockerfile.centos6 -t $(CENTOS_IMAGE) .
