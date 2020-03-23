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
DEFAULT_DEADLOCK ?= "enable"

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

default: setup

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

srcdep:
	rm -rf go-algorand
	git clone https://github.com/algorand/go-algorand

setup: srcdep buildsrc

buildsrc: crypto/libs/$(OS_TYPE)/$(ARCH)/lib/libsodium.a deps
	go install $(GOTRIMPATH) $(GOTAGS) $(GOBUILDMODE) -ldflags="$(GOLDFLAGS)" ./...

clean:
	go clean -i ./...
	cd go-algorand/crypto/libsodium-fork && \
		test ! -e Makefile || make clean
	rm -rf go-algorand/crypto/lib
	rm -rf go-algorand/crypto/libs
	rm -rf go-algorand/crypto/copies

.PHONY: default fmt vet lint check_license check_shell sanity cover prof deps build test fulltest shorttest clean cleango deploy node_exporter install %gen gen NONGO_BIN

