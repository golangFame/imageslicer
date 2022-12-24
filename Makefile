ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)
EXAMPLE_GO_MOD_DIRS := $(shell find ./examples/ -type f -name 'go.mod' -exec dirname {} \; | sort)

test:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go test in $${dir}"; \
	  (cd "$${dir}" && \
	    go test && \
	    env GOOS=linux GOARCH=386 go test && \
	    go vet); \
	done

mod:
	go get -u && go mod tidy -go=1.19
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go get -u ./... && \
	    go mod tidy -go=1.19); \
	donem

fmt:
	gofmt -w -s ./
	goimports -w  -local github.com/goferHiro/image-slicer ./

run-examples:
	set -e; for dir in $(EXAMPLE_GO_MOD_DIRS); do \
	  echo "go run . in $${dir}"; \
	  (cd "$${dir}" && go run .); \
	done

TAG := v0.8-alpha.1

tag:
	TAG=$(TAG) ./scripts/tag.sh

release:
	TAG=$(TAG) ./scripts/release.sh

deploy:
	make release
	make tag