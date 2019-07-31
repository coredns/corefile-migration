ARCH:=amd64 arm arm64 ppc64le s390x

.PHONY: test
test:
	cd migration; GO111MODULE=on go test -v -coverprofile=cover.out; cat cover.out >> ../coverage.txt
	cd corefile-tool/cmd; GO111MODULE=on go test -v -coverprofile=cover.out; cat cover.out >> ../coverage.txt

.PHONY: build
build:
	cd corefile-tool; GO111MODULE=on go build -o corefile-tool

.PHONY: release-binaries
release-binaries:
	mkdir -p build
	cd corefile-tool; for arch in $(ARCH); do \
	    GO111MODULE=on GOOS=linux GOARCH=$$arch go build -o ../build/corefile-tool-$$arch;\
	done
