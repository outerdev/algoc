default: setup

depsrc:
	@echo
	@echo "Getting the go-algorand submodule..."
	@echo
	@git submodule init
	@git submodule update
	@echo "===================================="

buildsrc:
	@echo
	@echo "Building the libsodium library..."
	@echo
	$(eval ARCH := $(shell ./go-algorand/scripts/archtype.sh))
	$(eval OS_TYPE := $(shell ./go-algorand/scripts/ostype.sh))
	make -C go-algorand crypto/libs/$(OS_TYPE)/$(ARCH)/lib/libsodium.a
	@echo "================================="

setup: depsrc buildsrc
	@echo
	@echo "Build the algoc binary with 'go build'"
	@echo
