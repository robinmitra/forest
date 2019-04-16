SHELL := /bin/bash

VERSION ?=
TYPE ?= patch

CURRENT_VERSION := $(shell cat VERSION)

GOOS ?= $(word 1, $(subst /, " ", $(word 4, $(shell go version))))

SOURCES := $(wildcard *.go **/*.go) $(MAKEFILE)

BUILD_DIR := build
BIN_DIR := bin

BINARY32 := forest_$(GOOS)_386
BINARY64 := forest_$(GOOS)_amd64

DOCKER_IMAGE := robinmitra/forest
DOCKER_APP_PATH := /go/src/app
DOCKER_VOL := $$(pwd):$(DOCKER_APP_PATH)

UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
	BINARY := $(BINARY64)
else ifeq ($(UNAME_M),amd64)
	BINARY := $(BINARY64)
else ifeq ($(UNAME_M),i686)
	BINARY := $(BINARY32)
else ifeq ($(UNAME_M),i386)
	BINARY := $(BINARY32)
else
$(error "Build on $(UNAME_M) is not supported yet.")
endif

.PHONY: ask_version check_version_provided check_version_not_current check_type_provided
.PHONY: check_prerequisites update_files commit tag bump push release

#########
# Build #
#########

# Build binary corresponding to the current architecture.
build: $(BUILD_DIR)/$(BINARY)

# Build the docker image.
docker-build:
	docker build -t $(DOCKER_IMAGE) .

$(BUILD_DIR)/$(BINARY32): $(SOURCES)
	docker run --rm -v $(DOCKER_VOL) -e GOARCH=386 -e GOOS=$(GOOS) $(DOCKER_IMAGE) go build -v -o $@

$(BUILD_DIR)/$(BINARY64): $(SOURCES)
	docker run --rm -v $(DOCKER_VOL) -e GOARCH=amd64 -e GOOS=$(GOOS) $(DOCKER_IMAGE) \
		go build -v -o $@

install: $(BIN_DIR)/forest

$(BIN_DIR):
	mkdir -p $@

# Ensure that target doesn't get rebuilt if $(BIN_DIR) gets updated, while still having it as a
# pre-requisite.
$(BIN_DIR)/forest: $(BUILD_DIR)/$(BINARY) | $(BIN_DIR)
	cp -f $< $@

########
# Test #
########

docker-test:
	docker run --rm -v $$(pwd):/go/src/app robinmitra/forest go test -v ./...

###########
# Version #
###########

ask_version:
	$(eval VERSION=$(shell read -p "* What's the new version? " version; echo $$version))
	$(eval TYPE=$(shell read -p "* What's the release type [major/minor/patch]? " t; echo $$t))

check_version_provided:
#	Too bad `ifndef` doesn't work with dynamically defined variables!
	@if [ "$(VERSION)" == "" ]; then echo "VERSION must be specified"; exit 1; fi

check_version_not_current:
#	Too bad `ifndef` doesn't work with dynamically defined variables!
	@if [ "$(VERSION)" == "$(CURRENT_VERSION)" ]; then \
	echo "VERSION cannot be same as CURRENT_VERSION"; exit 1; fi

check_type_provided:
#	Too bad `ifndef` doesn't work with dynamically defined variables!
	@if [ "$(TYPE)" == "" ]; then echo "TYPE must be specified"; exit 1; fi

check_prerequisites: check_version_provided check_version_not_current check_type_provided

# Update files with the new version.
update_files:
	@echo "==> Bumping version to $(VERSION)"
	@for file in "main.go" "VERSION"; do \
		sed -i '' "s/$(CURRENT_VERSION)/$(VERSION)/g" $$file; \
	done

# Commit version bump changes.
commit:
	@echo "==> Commiting version bump changes"
	git add main.go VERSION
	git commit -m "Bump version for $(TYPE) release to $(VERSION)."

# Create tag for the new version.
tag:
	@echo "==> Tagging $(TYPE) release v$(VERSION)"
	git tag v$(VERSION)

# Bump the version.
bump: ask_version check_prerequisites update_files commit tag

push:
	@echo "* Do you want to push the changes? [y/N]: "
	@read sure; \
	case "$$sure" in \
		[yY]) echo "==> Pushing changes up" && git push origin master && git push origin --tags;; \
	esac
	@echo "==> Finished"

release: bump push
