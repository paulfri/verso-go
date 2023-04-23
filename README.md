# 📚 verso

![build](https://github.com/versolabs/verso/actions/workflows/build.yml/badge.svg)

Syndication summarized.

## Dependencies

```
brew install golang goose go-task sqlc postgresql redis

go install github.com/cosmtrek/air@latest
```

## Usage

```
task
task --list-all
```

## Migrations

```
task db:create
task db:migrate
task db:reset
task db:gen -- migration_name
```

## Compiling SQL queries

```
task db:compile
```
