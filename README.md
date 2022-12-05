# go-graphql-ent-template
This is a golang project template using entgo and gqlgen heavaly inspired by this
[blog post](https://betterprogramming.pub/clean-architecture-with-ent-and-gqlgen-a789933a3665).

## Setup for Development

Make sure `docker-compose` is installed on your system
```bash
docker-compose up -d
```

The development server supports live reloads using [air](https://github.com/cosmtrek/air)

```bash
make start
```

## Unit/Integration Testing

First setup the database for testing

```bash
make test_setup_db
```

The tests can be executed via `make`

```bash
make test       # without coverage
make test_cov   # with code coverage saved to ./coverage
```

## E2E Testing

First setup the database for E2E testing

```bash
make e2e_setup_db
```

The E2E tests can be executed via `make`

```bash
make e2e
```