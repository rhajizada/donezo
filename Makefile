VOLUME_NAME=donezo-data
IMAGE_NAME=donezo
CONTAINER_NAME=donezo
SHELL_CONTAINER_NAME=donezo-debug


define prepreqs
	@if ! command -v yq >/dev/null 2>&1; then \
		echo "Error: yq is not installed. Please install yq and try again." >&2; \
		exit 1; \
	fi; \
	if [ ! -f "$(CONFIG_FILE)" ]; then \
		echo "Error: CONFIG_FILE \"$(CONFIG_FILE)\" does not exist. Please provide a valid file path." >&2; \
		exit 1; \
	fi;
endef


.PHONY: build
## build: Compile the packages
build:
	@go build -o bin/server ./cmd/server
	@go build -o bin/create-token ./cmd/create-token
	@go build -o bin/cli ./cmd/cli/


.PHONY: swagger
## swaggger: Genearate swagger docs
swagger:
	@swag init -g cmd/server/main.go


.PHONY: sqlc
## sqlc: Generate repository using sqlc
sqlc:
	@sqlc generate


.PHONY: generate-config
## generate-config: Generates a compatible config.yaml
generate-config:
	echo "port: 8000" > $(CONFIG_FILE) && \
	echo "database: /data/db.sqlite" >> $(CONFIG_FILE) && \
	echo "jwt:" >> $(CONFIG_FILE) && \
	echo "  secret: $$(head -c 32 /dev/urandom | base64)" >> $(CONFIG_FILE) && \
	echo "  expiration: 24h" >> $(CONFIG_FILE) && \
	echo "Configuration file created at $(CONFIG_FILE)"


.PHONY: run
## run: Build and run in development mode
run: pre
	@go run cmd/server/main.go -config "$(CONFIG_FILE)"


.PHONY: clean
## clean: Clean project and previous builds
clean:
	@rm builds/*


.PHONY: deps
## deps: Download modules
deps:
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
  @go install github.com/swaggo/swag/cmd/swag@latest && \
	@go mod download


.PHONY: build-image
## build-image: Build docker image
build-image:
	@docker build . -t $(IMAGE_NAME):latest


.PHONY: create-volume
## create-volume: Create docker volume
create-volume:
	@if [ "$(shell docker volume ls -q -f name=$(VOLUME_NAME))" = "$(VOLUME_NAME)" ]; then \
		echo "Volume $(VOLUME_NAME) already exists"; \
	else \
		echo "Creating volume $(VOLUME_NAME)"; \
		docker volume create $(VOLUME_NAME); \
	fi


.PHONY: run-container
## run-container: Launch a docker container
run-container: pre build-image create-volume
	@PORT=$$(yq '.port' $(CONFIG_FILE)); \
	docker run -d \
		--name $(CONTAINER_NAME) \
		-v $(VOLUME_NAME):/data \
		--mount type=bind,source=$(CONFIG_FILE),target=/etc/donezo/config.yaml \
		-p $$PORT:$$PORT \
		$(IMAGE_NAME)


.PHONY: rm-container
## rm-container: Stops and deletes container
rm-container:
	@docker stop $(CONTAINER_NAME)
	@docker rm $(CONTAINER_NAME)


.PHONY: create-token
## create-token: Create authentication token
create-token: pre
	@go run ./cmd/create-token/main.go -config $(CONFIG_FILE) $(ARGS)

.PHONY: cli
## cli: Launch CLI
cli:
	@go run ./cmd/cli/main.go $(ARGS)


.PHONY: create-token-container
## create-token-container: Create authentication token in running docker container
create-token-container:
	@docker exec -it $(CONTAINER_NAME) create-token $(ARGS)


.PHONY: shell
## shell: Launch shell inside docker container
shell: build-image create-volume
	@docker run --rm -it \
		--name $(SHELL_CONTAINER_NAME) \
		-v $(VOLUME_NAME):/data \
		--mount type=bind,source=$(CONFIG_FILE),target=/etc/donezo/config.yaml \
		-p 8000:8000 \
		$(IMAGE_NAME) /bin/sh


.PHONY: help
all: help
# help: show help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

pre:
	$(prereqs)
