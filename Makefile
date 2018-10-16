# Constants
MK_FILE_PATH = $(lastword $(MAKEFILE_LIST))
PRJ_DIR      = $(abspath $(dir $(MK_FILE_PATH)))
GOPATH_LOCAL = $(PRJ_DIR)/build
GOPATH_DIR   = github.com/cloudflare/sidh
VENDOR_DIR   = build/vendor
CSHAKE_PKG   ?= github.com/henrydcase/nobs/hash/sha3
TARGETS      = p503 p751 sidh sike
GO           ?= go
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

all: test
clean:
	rm -rf $(GOPATH_LOCAL)
	rm -rf coverage*.txt

build_env:
	GOPATH=$(GOPATH_LOCAL) $(GO) get $(CSHAKE_PKG)
	mkdir -p $(GOPATH_LOCAL)/src/$(GOPATH_DIR)
	cp -rf internal $(GOPATH_LOCAL)/src/$(GOPATH_DIR)
	cp -rf etc $(GOPATH_LOCAL)/src/$(GOPATH_DIR)

copy-target-%:
	cp -rf $* $(GOPATH_LOCAL)/src/$(GOPATH_DIR)

prep_targets: build_env $(addprefix copy-target-, $(TARGETS))

install-%: prep_targets
	GOPATH=$(GOPATH_LOCAL) GOARCH=$(GOARCH) $(GO) install $(OPTS) $(GOPATH_DIR)/$*

test-%: prep_targets
	GOPATH=$(GOPATH_LOCAL) $(GO) vet $(GOPATH_DIR)/$*
	GOPATH=$(GOPATH_LOCAL) GOARCH=$(GOARCH) $(GO) test $(OPTS) $(GOPATH_DIR)/$*

bench-%: prep_targets
	GOMAXPROCS=1 GOPATH=$(GOPATH_LOCAL) $(GO) test $(OPTS) $(GOPATH_DIR)/$* $(BENCH_OPTS)

cover-%: prep_targets
	GOPATH=$(GOPATH_LOCAL) $(GO) test \
		-race -coverprofile=coverage_$*.txt -covermode=atomic $(OPTS) $(GOPATH_DIR)/$*
	cat coverage_$*.txt >> coverage.txt
	rm coverage_$*.txt

vendor: clean
	mkdir -p $(VENDOR_DIR)/github_com/cloudflare/sidh/
	rsync -a . $(VENDOR_DIR)/github_com/cloudflare/sidh/ \
		--exclude=$(VENDOR_DIR) \
		--exclude=.git          \
		--exclude=.travis.yml   \
		--exclude=README.md     \
		--exclude=Makefile      \
		--exclude=build
	# This swaps all imports with github.com to github_com, so that standard library doesn't
	# try to access external libraries.
	find $(VENDOR_DIR) -type f -iname "*.go" -print0  | xargs -0 sed -i 's/github\.com/github_com/g'

bench:   $(addprefix bench-,   $(TARGETS))
cover:   $(addprefix cover-,   $(TARGETS))
install: $(addprefix install-, $(TARGETS))
test:    $(addprefix test-,    $(TARGETS))
