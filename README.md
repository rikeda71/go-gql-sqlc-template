# go-gql-sqlc-template

Template Project of Web Backend API

Use Following Techstack

- [Go](https://github.com/golang/go) : Modern Programing Language
- [gqlgen](https://github.com/99designs/gqlgen) : Use GraphQL Schema
- [sqlc](https://github.com/sqlc-dev/sqlc) : SQL based OR Mapper
- [dbmate](https://github.com/amacneil/dbmate) : Migration Tool
- [task](https://github.com/go-task/task) : Task Runner

## Setup

### Installation

Assuming the use of MacOS

1. Install followings

- Go 1.23
- Docker Environmet (ex. Docker Desktop)

2. Execute following commands

    ```bash
    $ go install github.com/go-task/task/v3/cmd/task@latest
    $ task setup
    # add 'export PATH="/opt/homebrew/opt/libpq/bin:$PATH"' to .zshrc etc.
    ```

## Start API

```bash
$ task up
# wait a minutes...
$ task migrate # setup db
```

## Test

```bash
$ task test
```

### Integration Test

```bash
$ task integration-test
```

## Development

### Generate API I/F From Graphql Schema

```bash
$ task gql-gen
```

### Generate DB I/F From SQL

```bash
# OR Mapper Code from DML
$ task sqlc-gen

# Model Code from DDL
$ task migrate # generate `db/schema.sql` through `dbmate migrate`
$ task sqlc-gen
```

### DB Operations

```bash
# migrate
$ task migrate

# rollback
$ task rollback
```
