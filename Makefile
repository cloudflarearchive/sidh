# Constants
MK_FILE_PATH = $(lastword $(MAKEFILE_LIST))
PRJ_DIR      = $(abspath $(dir $(MK_FILE_PATH)))
GOPATH_LOCAL = $(PRJ_DIR)/build
GOPATH_DIR   = github.com/cloudflare/p751sidh
CSHAKE_PKG   = github.com/henrydcase/nobs/hash/sha3
CPU_PKG      = golang.org/x/sys/cpu
TARGETS      = p751toolbox sidh sike
GO           ?= go
GOARCH       ?=
OPTS_GCCGO   ?= -compiler gccgo -O2 -g
OPTS_TAGS    ?= -tags=noasm
OPTS         ?=
NOASM        ?=

ifeq ($(NOASM),1)
	OPTS+=$(OPTS_TAGS)
endif

clean:
	rm -rf $(GOPATH_LOCAL)
	rm -rf coverage*.txt

prep:
	GOPATH=$(GOPATH_LOCAL) $(GO) get $(CSHAKE_PKG) $(CPU_PKG)
	mkdir -p $(GOPATH_LOCAL)/src/$(GOPATH_DIR)
	cp -rf p751toolbox $(GOPATH_LOCAL)/src/$(GOPATH_DIR)
	cp -rf sidh $(GOPATH_LOCAL)/src/$(GOPATH_DIR)
	cp -rf sike $(GOPATH_LOCAL)/src/$(GOPATH_DIR)
	cp -rf etc $(GOPATH_LOCAL)/src/$(GOPATH_DIR)

test-%: prep
	GOPATH=$(GOPATH_LOCAL) $(GO) test -v $(OPTS) $(GOPATH_DIR)/$*

bench-%: prep
	cd $*; GOPATH=$(GOPATH_LOCAL) $(GO) test -v $(OPTS) -bench=.

cover-%: prep
	GOPATH=$(GOPATH_LOCAL) $(GO) test \
		-race -coverprofile=coverage_$*.txt -covermode=atomic $(OPTS) $(GOPATH_DIR)/$*
	cat coverage_$*.txt >> coverage.txt
	rm coverage_$*.txt

test: $(addprefix test-, $(TARGETS))
bench: $(addprefix bench-, $(TARGETS))
cover: $(addprefix cover-, $(TARGETS))
