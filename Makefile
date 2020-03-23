GOPATH		:= $(shell go env GOPATH)
GOPATH1		:= $(firstword $(subst :, ,$(GOPATH)))
export GOPATH
GO111MODULE	:= on
export GO111MODULE
UNAME		:= $(shell uname)
SRCPATH     := $(shell pwd)/go-algorand
ARCH        := $(shell ./go-algorand/scripts/archtype.sh)
OS_TYPE     := $(shell ./go-algorand/scripts/ostype.sh)

# If build number already set, use it - to ensure same build number across multiple platforms being built
BUILDNUMBER      ?= $(shell ./go-algorand/scripts/compute_build_number.sh)
COMMITHASH       := $(shell ./go-algorand/scripts/compute_build_commit.sh)
BUILDBRANCH      ?= $(shell ./go-algorand/scripts/compute_branch.sh)
BUILDCHANNEL     ?= $(shell ./go-algorand/scripts/compute_branch_channel.sh $(BUILDBRANCH))
DEFAULTNETWORK   ?= $(shell ./go-algorand/scripts/compute_branch_network.sh $(BUILDBRANCH))
DEFAULT_DEADLOCK ?= $(shell ./go-algorand/scripts/compute_branch_deadlock_default.sh $(BUILDBRANCH))

GOTAGSLIST          := sqlite_unlock_notify sqlite_omit_load_extension

ifeq ($(UNAME), Linux)
EXTLDFLAGS := -static-libstdc++ -static-libgcc
ifeq ($(ARCH), amd64)
# the following predicate is abit misleading; it tests if we're not in centos.
ifeq (,$(wildcard /etc/centos-release))
EXTLDFLAGS  += -static
endif
GOTAGSLIST  += osusergo netgo static_build
GOBUILDMODE := -buildmode pie
endif
ifeq ($(ARCH), arm)
ifneq ("$(wildcard /etc/alpine-release)","")
EXTLDFLAGS  += -static
GOTAGSLIST  += osusergo netgo static_build
GOBUILDMODE := -buildmode pie
endif
endif
endif

GOTAGS      := --tags "$(GOTAGSLIST)"
GOTRIMPATH	:= $(shell go help build | grep -q .-trimpath && echo -trimpath)

GOLDFLAGS_BASE  := -X github.com/outerdev/algoc/go-algorand/config.BuildNumber=$(BUILDNUMBER) \
		 -X github.com/outerdev/algoc/go-algorand/config.CommitHash=$(COMMITHASH) \
		 -X github.com/outerdev/algoc/go-algorand/config.Branch=$(BUILDBRANCH) \
		 -X github.com/outerdev/algoc/go-algorand/config.DefaultDeadlock=$(DEFAULT_DEADLOCK) \
		 -extldflags \"$(EXTLDFLAGS)\"

GOLDFLAGS := $(GOLDFLAGS_BASE) \
		 -X github.com/outerdev/algoc/go-algorand/config.Channel=$(BUILDCHANNEL)

default: build

# tools

generate: deps
	PATH=$(GOPATH1)/bin:$$PATH go generate ./...

ALWAYS:

# build our fork of libsodium, placing artifacts into crypto/lib/ and crypto/include/
crypto/libs/$(OS_TYPE)/$(ARCH)/lib/libsodium.a:
	mkdir -p go-algorand/crypto/copies/$(OS_TYPE)/$(ARCH)
	cp -R go-algorand/crypto/libsodium-fork go-algorand/crypto/copies/$(OS_TYPE)/$(ARCH)/libsodium-fork
	cd go-algorand/crypto/copies/$(OS_TYPE)/$(ARCH)/libsodium-fork && \
		./autogen.sh --prefix $(SRCPATH)/crypto/libs/$(OS_TYPE)/$(ARCH) && \
		./configure --disable-shared --prefix="$(SRCPATH)/crypto/libs/$(OS_TYPE)/$(ARCH)" && \
		$(MAKE) && \
		$(MAKE) install

deps:
	./go-algorand/scripts/check_deps.sh

# develop

build: buildsrc gen

buildsrc: crypto/libs/$(OS_TYPE)/$(ARCH)/lib/libsodium.a node_exporter NONGO_BIN deps
	go install $(GOTRIMPATH) $(GOTAGS) $(GOBUILDMODE) -ldflags="$(GOLDFLAGS)" ./...

NONGO_BIN_FILES=$(GOPATH1)/bin/find-nodes.sh $(GOPATH1)/bin/update.sh $(GOPATH1)/bin/COPYING $(GOPATH1)/bin/ddconfig.sh

NONGO_BIN: $(NONGO_BIN_FILES)

$(GOPATH1)/bin/find-nodes.sh: go-algorand/scripts/find-nodes.sh

$(GOPATH1)/bin/update.sh: go-algorand/cmd/updater/update.sh

$(GOPATH1)/bin/COPYING: go-algorand/COPYING

$(GOPATH1)/bin/ddconfig.sh: go-algorand/scripts/ddconfig.sh

$(GOPATH1)/bin/%:
	cp -f $< $@

clean:
	go clean -i ./...
	rm -f $(GOPATH1)/bin/node_exporter
	cd crypto/libsodium-fork && \
		test ! -e Makefile || make clean
	rm -rf crypto/lib
	rm -rf crypto/libs
	rm -rf crypto/copies

# assign the phony target node_exporter the dependency of the actual executable.
node_exporter: $(GOPATH1)/bin/node_exporter

# The recipe for making the node_exporter is by extracting it from the gzipped&tar file.
# The file is was taken from the S3 cloud and it traditionally stored at
# /travis-build-artifacts-us-ea-1.algorand.network/algorand/node_exporter/latest/node_exporter-stable-linux-x86_64.tar.gz
$(GOPATH1)/bin/node_exporter:
	tar -xzvf installer/external/node_exporter-stable-$(shell ./scripts/ostype.sh)-$(shell uname -m | tr '[:upper:]' '[:lower:]').tar.gz -C $(GOPATH1)/bin

.PRECIOUS: gen/%/genesis.json

# devnet & testnet
NETWORKS = testnet devnet

gen/%/genesis.dump: gen/%/genesis.json
	./scripts/dump_genesis.sh $< > $@

gen/%/genesis.json: gen/%.json gen/generate.go buildsrc
	$(GOPATH1)/bin/genesis -q -n $(shell basename $(shell dirname $@)) -c $< -d $(subst .json,,$<)

gen: $(addsuffix gen, $(NETWORKS)) mainnetgen

$(addsuffix gen, $(NETWORKS)): %gen: gen/%/genesis.dump

# mainnet

gen/mainnet/genesis.dump: gen/mainnet/genesis.json
	./scripts/dump_genesis.sh gen/mainnet/genesis.json > gen/mainnet/genesis.dump

mainnetgen: gen/mainnet/genesis.dump

gen/mainnet/genesis.json: gen/pregen/mainnet/genesis.csv buildsrc
	mkdir -p gen/mainnet
	cat gen/pregen/mainnet/genesis.csv | $(GOPATH1)/bin/incorporate -m gen/pregen/mainnet/metadata.json > gen/mainnet/genesis.json

.PHONY: default fmt vet lint check_license check_shell sanity cover prof deps build test fulltest shorttest clean cleango deploy node_exporter install %gen gen NONGO_BIN

