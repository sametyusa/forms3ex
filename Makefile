.DEFAULT_GOAL := default

clean:
	@rm ./bin/*

default: install-deps build test

build: install-deps
	@./bin/hey go build -o bin/innsecure ./cmd/innsecure
	@./bin/hey go build -o bin/token ./cmd/token

bin_dir:
	@mkdir -p ./bin

install-deps: install-hey install-goimports

install-hey: bin_dir
	@curl -L --insecure https://github.com/rossmcf/hey/releases/download/v1.0.0/installer.sh | bash
	@mv hey bin

install-goimports:
	@if [ ! -f ./goimports ]; then \
		cd ~ && go get -u golang.org/x/tools/cmd/goimports; \
	fi

test:
	@echo "executing tests..."
	go test github.com/form3tech/innsecure

# package for release to candidates (ignore for test exercise)
package-%:
	@echo $*
	@cd ..&& pwd && tar -czvf innsecure-$*.tar.gz --exclude={".git",".github","bin","releases"} innsecure
	@mkdir -p releases
	@mv ../innsecure-$*.tar.gz releases 

.PHONY: clean build test package-%
