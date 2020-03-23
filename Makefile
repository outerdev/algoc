ARCH        := $(shell ./go-algorand/scripts/archtype.sh)
OS_TYPE     := $(shell ./go-algorand/scripts/ostype.sh)

default: setup

depsrc:
	@echo "Getting the go-algorand submodule..."
	rm -rf go-algorand
	git clone https://github.com/algorand/go-algorand

setup: depsrc
	@echo "Building the libsodium library..."
	make -C go-algorand crypto/libs/$(OS_TYPE)/$(ARCH)/lib/libsodium.a

