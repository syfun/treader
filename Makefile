# This repo's root import path (under GOPATH).
PKG := github.com/syfun/treader

# Which architecture to build - see $(ALL_ARCH) for options.
ARCH ?= amd64

# This version-strategy uses git tags to set the version string
VERSION := 0.1.0
#
# This version-strategy uses a manual value to set the version string
#VERSION := 1.2.3

###
### These variables should not need tweaking.
###

SRC_DIRS := cmd pkg # directories which hold app source (not vendored)

ALL_ARCH := amd64 arm arm64 ppc64le

# If you want to build all binaries, see the 'all-build' rule.
all: build

build-%:
	@$(MAKE) --no-print-directory ARCH=$* build

run-search:
	./bin/$(ARCH)/treadsearch -book=$(BOOK)

all-build: $(addprefix build-, $(ALL_ARCH))

build: bin/$(ARCH)/$(BIN)

bin/$(ARCH)/$(BIN): build-dirs
	@echo "building: $@"
	ARCH=$(ARCH) VERSION=$(VERSION) PKG=$(PKG) ./build/build.sh 

test: build-dirs
	./build/test.sh $(SRC_DIRS)  

build-dirs:
	@mkdir -p bin/$(ARCH)

clean: 
	rm -rf bin
	
