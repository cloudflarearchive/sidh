# Constants
MK_FILE_PATH = $(lastword $(MAKEFILE_LIST))
PRJ_DIR      = $(abspath $(dir $(MK_FILE_PATH)))
GOPATH_LOCAL = $(PRJ_DIR)/build
GOPATH_DIR   = github.com/cloudflare/p751sidh
CSHAKE_PKG   ?= github.com/henrydcase/nobs/hash/sha3
TARGETS      = p751toolbox sidh sike
GOARCH       ?=
OPTS_GCCGO   ?= -compiler gccgo -O2 -g
OPTS         ?=
OPTS_TAGS    ?= -tags=noasm
NOASM        ?=
# -run="NonExistent" is set to make sure tests are not run before benchmarking
BENCH_OPTS   ?= -bench=. -run="NonExistent"
# whether to be verbose
V            ?= 1

ifeq ($(NOASM),1)
	OPTS+=$(OPTS_TAGS)
endif

ifeq ($(V),1)
	OPTS += -v              # Be verbose
	BENCH_OPTS += -gcflags=-m     # Show results from inlining
endif

clean:
	rm -rf $(GOPATH_LOCAL)
	rm -rf coverage*.txt

build_env:
	GOPATH=$(GOPATH_LOCAL) go get $(CSHAKE_PKG)
	mkdir -p $(GOPATH_LOCAL)/src/$(GOPATH_DIR)
	cp -rf etc $(GOPATH_LOCAL)/src/$(GOPATH_DIR)

copy-target-%:
	cp -rf $* $(GOPATH_LOCAL)/src/$(GOPATH_DIR)

prep_targets: build_env $(addprefix copy-target-, $(TARGETS))

test-%: prep_targets
	GOPATH=$(GOPATH_LOCAL) go test $(OPTS) $(GOPATH_DIR)/$*

bench-%: prep_targets
	cd $*; GOPATH=$(GOPATH_LOCAL) go test $(OPTS) $(BENCH_OPTS)

cover-%: prep_targets
	GOPATH=$(GOPATH_LOCAL) go test \
		-race -coverprofile=coverage_$*.txt -covermode=atomic $(OPTS) $(GOPATH_DIR)/$*
	cat coverage_$*.txt >> coverage.txt
	rm coverage_$*.txt

test: $(addprefix test-, $(TARGETS))
bench: $(addprefix bench-, $(TARGETS))
cover: $(addprefix cover-, $(TARGETS))