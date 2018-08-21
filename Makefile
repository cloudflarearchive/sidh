# Constants
MK_FILE_PATH = $(lastword $(MAKEFILE_LIST))
PRJ_DIR      = $(abspath $(dir $(MK_FILE_PATH)))
GOPATH_LOCAL = $(PRJ_DIR)/build
GOPATH_PKG   = src/github.com/cloudflare/p751sidh
CSHAKE_PKG   = github.com/henrydcase/nobs/hash/sha3
TARGETS      = p751toolbox sidh sike

clean:
	rm -rf $(GOPATH_LOCAL)
	rm -rf coverage*.txt

prep:
	GOPATH=$(GOPATH_LOCAL) go get $(CSHAKE_PKG)
	mkdir -p $(GOPATH_LOCAL)/$(GOPATH_PKG)
	cp -rf p751toolbox $(GOPATH_LOCAL)/$(GOPATH_PKG)
	cp -rf sidh $(GOPATH_LOCAL)/$(GOPATH_PKG)

test-%: clean prep
	GOPATH=$(GOPATH_LOCAL) go test -race -v ./$*

bench-%: clean prep
	cd $*; GOPATH=$(GOPATH_LOCAL) go test -v -bench=.

cover-%: clean prep
	GOPATH=$(GOPATH_LOCAL) go test -race -coverprofile=coverage_$*.txt -covermode=atomic ./$*
	cat coverage_$*.txt >> coverage.txt

test: $(addprefix test-, $(TARGETS))
bench: $(addprefix bench-, $(TARGETS))
cover: $(addprefix cover-, $(TARGETS))
