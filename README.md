# Tezos Delegation Service

## Overview

In this exercise, I built a Golang service that gathers new delegations made on the Tezos protocol and exposes them through a public API.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.22 or later)
- [Docker](https://www.docker.com/products/docker-desktop) (for containerization)
- [Docker Compose](https://docs.docker.com/compose/) (optional, for managing multi-container Docker applications)
- [PostgreSQL](https://www.postgresql.org/) (for local development)
- [Goose](https://github.com/pressly/goose) (to run migration)

## Architecture

### cmd

cmd contains the file main.go, we will define how the server/database/redis runs and start the gRPC and HTTP server if we need it.

### internal

internal is composed by several folders. It's the core of the microservice.

#### config

Like his name, this folder contain all the config that the microservice need.

#### database

The database folder contain all the **migrations** and the **postgres** logics.

#### entities

I define all the entities I need in this folder.

#### handlers

I define here the http server and all the handlers methods.

#### service

It's the service layer, our handlers will call the service methods and the service call directly database methods.

### pkg

pkg contains all the libraries that we use in the different micro-services. These libraries are non-product ones. It can be anything from a client connecting to our database to a utility function improving the way we handle the errors.

## Getting Started

Make sure you have goose installed on your computer

```go install github.com/pressly/goose/v3/cmd/goose@latest```

### Update environments variables

**docker-compose.yml** is already configured but feel free to update environments variables.

### Run the project

```make run```

### Retrieve delegations

Fetch all delegations:
```curl -X GET http://localhost:3000/xtz/delegations```

Fetch delegations by year:
```curl -X GET http://localhost:3000/xtz/delegations?year=2024```

Optional parameters:

limit (default 100)
offset (default 0)

### Connnect to the database

```postgresql://root:root@127.0.0.1/app-db?tLSMode=0```

### Run test coverage

```make run-tests```

