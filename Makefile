default: setup

depsrc:
	@echo "Getting the go-algorand submodule..."
	@echo
	rm -rf go-algorand
	git clone https://github.com/algorand/go-algorand
	@echo "===================================="

setup: depsrc
	@echo "Building the libsodium library..."
	@echo
	$(eval ARCH := $(shell ./go-algorand/scripts/archtype.sh))
	$(eval OS_TYPE := $(shell ./go-algorand/scripts/ostype.sh))
	make -C go-algorand crypto/libs/$(OS_TYPE)/$(ARCH)/lib/libsodium.a
	@echo "================================="

