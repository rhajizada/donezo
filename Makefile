VOLUME_NAME=donezo_data
IMAGE_NAME=donezo
CONTAINER_NAME=donezo


.PHONY: build
## build: Compile the packages
build:
	@go build -o server ./cmd/server
	@go build -o create-token ./cmd/create-token


.PHONY: swagger
## swaggger: Genearate swagger docs
swagger:
	@swag init -g cmd/server/main.go


.PHONY: sqlc
## sqlc: Generate repository using sqlc
sqlc:
	@sqlc generate

.PHONY: run
## run: Build and run in development mode
run: build
	@./server $(ARGS)


.PHONY: clean
## clean: Clean project and previous builds
clean:
	@rm -f server
	@rm -f create-token


.PHONY: deps
## deps: Download modules
deps:
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
run-container: build-image create-volume
	@docker run -d \
		--name $(CONTAINER_NAME) \
		-v $(VOLUME_NAME):/data \
		-p 8000:8000 \
		$(IMAGE_NAME)

.PHONY: create-token
## create-token: Create authentication token
create-token:
	docker exec -it $(CONTAINER_NAME) create-token

.PHONY: shell
## shell: Launch shell inside docker container
shell: build-image create-volume
	@docker run --rm -it \
		--name $(CONTAINER_NAME) \
		-v $(VOLUME_NAME):/path/in/container \
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
