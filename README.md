# donezo

![Go](https://img.shields.io/badge/Go-1.22-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Docker](https://img.shields.io/badge/Docker-20.10.7-blue.svg)

**donezo** is a simple todo list web application built with Go. Designed as an
educational experiment, donezo leverages Go's standard `net/http` library
enhanced with the new `mux` capabilities introduced in Go 1.22. It features a
RESTful API, a Go API and a text-based user interface (TUI) built
with the [`bubbletea`](https://github.com/charmbracelet/bubbletea) package. Additionally, donezo includes comprehensive API
documentation available via `Swagger`.

## Table of Contents

- [Features](#features)
- [Project Details](#project-details)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Build](#build)
  - [Configuration](#configuration)
- [Running the Application](#running-the-application)
  - [Locally](#locally)
  - [Using Docker](#using-docker)
- [Makefile Targets](#makefile-targets)
- [API Documentation](#api-documentation)
- [Usage](#usage)
  - [Client API](#client-api)

## Features

- **RESTful API**: Built with Go's `net/http` package, providing endpoints to
  manage boards and items.
- **Client API**: A Go-based client to interact with the REST API.
- **TUI**: A text-based user interface built with the [`bubbletea`](https://github.com/charmbracelet/bubbletea) package for managing tasks directly from the terminal.
- **Swagger Documentation**: Comprehensive API documentation available at the `/swagger` endpoint.
- **Database Migrations**: Managed using [`goose`](https://github.com/pressly/goose) with SQLite as the database.
- **Docker Support**: Easily deployable using Docker with pre-configured `Make` targets.

## Project Details and Objectives

donezo was developed as a hands-on project to enhance practical skills in
building a complete web application from the ground up. Inspired by the
common practice among developers to build simple "to-do" apps
when learning new programming languages, donezo extends this concept by
implementing diverse application components:

- Backend Development: Building a RESTful API to handle CRUD operations for boards and items, ensuring efficient and secure data handling.
- Client Interaction: Creating a command-line client to interact with the API, facilitating programmatic access and manipulation of data.
- User Interface Design: Developing a text-based user interface (TUI) to offer an intuitive and accessible way for users to manage their tasks directly from the terminal.
- Documentation and Testing: Integrating comprehensive API documentation and implementing tests to validate functionality and ensure reliability.

## Installation

### Prerequisites

- **Go**: Version 1.22 or higher. [Install Go](https://golang.org/doc/install)
- **Make**: For running Makefile targets. [Install Make](https://www.gnu.org/software/make/)
- **Docker**: For containerized deployment. [Install Docker](https://docs.docker.com/get-docker/)
- **SQLC**: SQL to Go code generator. [Install SQLC](https://sqlc.dev/)

### Build
Run `make build` to compile the server, and create-token binaries located in the `bin/` directory.

### Configuration

Create a configuration file `config.yaml` in `/etc/donezo/config.yaml` or specify a different path when running commands. An example configuration is provided:

```yaml
port: 8000
database: /data/db.sqlite
jwt:
  secret: your_jwt_secret_key_here
  expiration: 24h

# port: The port on which the server will run.
# database: Path to the SQLite database file.
# jwt.secret: Secret key for signing JWT tokens.
# jwt.expiration: Defaul token expiration duration.
```

Ensure that the `config.yaml` file has appropriate permissions (`0600`) to secure sensitive information, especially the JWT secret.

## Running the Application
### Locally
Start the server using the Makefile:
```bash
make run
```

This command will:
- Build executables.
- Apply database migrations using goose.
- Start the REST API server on the configured port.

### Using Docker
You can run the following `make` target to just build the image.

```bash
make build-image
```

Run the Docker Container

```bash
make run-container
```

The Makefile handles building the Docker image and running the container with the necessary configurations.

## Makefile Targets
```
build                    Compile the executables
swaggger                 Genearate swagger docs
sqlc                     Generate repository using sqlc
generate-config          Generates a compatible config.yaml
run                      Build and run in development mode
clean                    Clean project and previous builds
deps                     Download modules
build-image              Build docker image
create-volume            Create docker volume
run-container            Launch a docker container
rm-container             Stops and deletes container
create-token             Create authentication token
create-token-container   Create authentication token in running docker container
shell                    Launch shell inside docker container
```

Refer to the Makefile for additional targets and customization options.

## API Documentation
`donezo` provides comprehensive API documentation using Swagger. After running the server, access the Swagger UI at:

http://localhost:8000/swagger


## Usage
### Client API

The client API allows you to interact with donezo programmatically.
Initialize the Client

```go
package main

import (
    "log"
    "time"

    "github.com/rhajizada/donezo/pkg/client"
)

func main() {
    c := client.New(
        "http://localhost:8000",
        "your_api_token_here",
        time.Second*5,
    )

    // Check health
    if err := c.Healthy(); err != nil {
        log.Fatalf("Health check failed: %v", err)
    }

    // Validate token
    if err := c.ValidateToken(); err != nil {
        log.Fatalf("Token validation failed: %v", err)
    }

    // List boards
    boards, err := c.ListBoards()
    if err != nil {
        log.Fatalf("Failed to list boards: %v", err)
    }
    log.Printf("Boards: %+v", boards)
}
```
