SHELL := /bin/bash
VERSION ?=
TYPE ?= patch

CURRENT_VERSION := $(shell cat VERSION)

.PHONY: ask_version check_version_provided check_version_not_current check_type_provided
.PHONY: check_prerequisites update_files commit tag bump push release

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
