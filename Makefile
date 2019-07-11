ARCH:=amd64 arm arm64 ppc64le s390x

.PHONY: test
test:
	cd migration; GO111MODULE=on go test -v
	cd corefile-tool; GO111MODULE=on go test -v

.PHONY: build
build:
	cd corefile-tool; GO111MODULE=on go build -o corefile-tool

.PHONY: linux-archs
linux-archs:
	cd corefile-tool; for arch in $(ARCH); do \
	    mkdir -p build/linux/$$arch && GO111MODULE=on GOOS=linux GOARCH=$$arch go build -o build/linux/$$arch/corefile-tool;\
	done
