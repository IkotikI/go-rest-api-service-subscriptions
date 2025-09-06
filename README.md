# REST API subscriptions service

Service provide REST API for CRUDL operations over *Subscriptions* for relation between *Users* and *Services*.

See swagger documentation provided by default route:
[http://localhost:8020/api/v1/swagger/index.html](http://localhost:8020/api/v1/swagger/index.html)

## Run

1. Provide correct configuration of PostgresSQL user, password and database name to ENV (`.env`) and `./config/config.yaml`.

2. Execute
```
make run
```

## Build Docker
1. Provide correct ENV (`.env`), specify `config.yaml`.
   
2. Run  
```
docker compose up
```

### Possible docker build issues
1. `config.yaml: httpServer.port` should match microservice port `docker-compose.yml`
   
2. `wait-for-postgres.sh` should have `LF` (UNIX-system's) line sequence 

## Project structure

### Enteties
Enteties described in files:
- General: [`subscription.go`](./subscription.go)

- Database model: [`./internal/storage/model.go`](./internal/storage/model.go)

- Responses: [`./internal/server/http/handler/util.go`](./internal/server/http/handler/util.go)

### Architecture
Project separated by layers, using interfaces:

- Storage/Repository layer ([`./internal/storage`](./internal/storage/))

- Service layer ([`./internal/service`](./internal/service/))

- Handlers/Controller layer (for HTTP [`./internal/server/http/handlers`](./internal/server/http/handlers/))

### Migrations
Database migrations implements with [`goose`](https://github.com/pressly/goose) package.

Migrations runs from tool-file [`./tools/migrations/goose.go`](./tools/migrations/goose.go) and user [`./schema/`](./schema/) directory for SQL migration files.

### Documentation
API handler provide [Swagger](https://swagger.io/) documentation. By default path, you can find it for docker:
[`http://localhost:8020/api/v1/swagger/index.html`](http://localhost:8020/api/v1/swagger/index.html)
or specify your parameters for host, port.

### API Client
Besides [Swagger](https://swagger.io/), there also [Bruno](https://www.usebruno.com/) directory in [`./api/subscription_microservice`](./api/subscription_microservice)

### Utilities packages
Other packages, that used in the project:
- [air](https://github.com/air-verse/air) - auto-rebuild project
- [testify](https://github.com/stretchr/testify), [mockery](https://github.com/vektra/mockery), [sqlmock](https://github.com/DATA-DOG/go-sqlmock) - mock & testing
- [zerolog](https://github.com/rs/zerolog) - primary logging package
- 