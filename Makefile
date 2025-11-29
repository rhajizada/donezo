VERSION := $(shell \
  tag=$$(git describe --tags --exact-match 2>/dev/null || true); \
  if [ -n "$$tag" ] && echo "$$tag" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
    echo "$$tag" | sed 's/^v//'; \
  else \
    git rev-parse --short HEAD; \
  fi)


.PHONY: build
## build: Build project
build:
	go build -ldflags "-s -w -X 'main.Version=$(VERSION)'" -o bin/donezo


.PHONY: clean
## clean: Remove previous builds
clean:
	@rm bin/*


.PHONY: mod
## mod: Install dependecies
mod:
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest \
	@go mod download


.PHONY: install
## install: Install donezo
install: build
	@mkdir -p $(HOME)/.local/bin
	@cp ./bin/donezo $(HOME)/.local/bin
	@echo "Installed donezo to '$(HOME)/.local/bin'. Please add '$(HOME)/.local/bin' to your PATH."


.PHONY: lint
## lint: Lint source code
lint:
	@golangci-lint run

.PHONY: test
## test: Run tests
test:
	@go tool gotestsum

.PHONY: run
## run: Build and run in development mode
run:
	@go run main.go


.PHONY: sqlc
## sqlc: Generate repository code using sqlc
sqlc:
	@go tool sqlc generate


.PHONY: uninstall
## uninstall: Uninstall TUI
uninstall:
	@rm $(HOME)/.local/bin/donezo


.PHONY: help
all: help
# help: show help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
