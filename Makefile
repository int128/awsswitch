# CircleCI specific variables
CIRCLE_TAG ?= latest
GITHUB_USERNAME := $(CIRCLE_PROJECT_USERNAME)
GITHUB_REPONAME := $(CIRCLE_PROJECT_REPONAME)

TARGET := awsswitch
TARGET_OSARCH := linux_amd64 darwin_amd64 windows_amd64 linux_arm linux_arm64
VERSION ?= $(CIRCLE_TAG)
LDFLAGS := -X main.version=$(VERSION)

all: $(TARGET)

$(TARGET): $(wildcard **/*.go)
	go build -o $@ -ldflags "$(LDFLAGS)"

.PHONY: check
check:
	golangci-lint run

.PHONY: dist
dist: dist/output
dist/output:
	# make the zip files for GitHub Releases
	VERSION=$(VERSION) CGO_ENABLED=0 goxzst -d dist/output -i "LICENSE" -o "$(TARGET)" -osarch "$(TARGET_OSARCH)" -t "dist/awsswitch.rb" -- -ldflags "$(LDFLAGS)"
	# test the zip file
	zipinfo dist/output/awsswitch_linux_amd64.zip

.PHONY: release
release: dist
	# publish the binaries
	ghcp release -u "$(GITHUB_USERNAME)" -r "$(GITHUB_REPONAME)" -t "$(VERSION)" dist/output/
	# publish the Homebrew formula
	ghcp commit -u "$(GITHUB_USERNAME)" -r "homebrew-$(GITHUB_REPONAME)" -b "bump-$(VERSION)" -m "Bump the version to $(VERSION)" -C dist/output/ awsswitch.rb
	ghcp pull-request -u "$(GITHUB_USERNAME)" -r "homebrew-$(GITHUB_REPONAME)" -b "bump-$(VERSION)" --title "Bump the version to $(VERSION)"

.PHONY: clean
clean:
	-rm $(TARGET)
	-rm -r dist/output/
