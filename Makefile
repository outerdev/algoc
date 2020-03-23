default: setup

depsrc:
	@echo
	@echo "Getting the go-algorand submodule..."
	@echo
	rm -rf go-algorand
	git clone https://github.com/algorand/go-algorand
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
