# 📚 verso

![build](https://github.com/versolabs/verso/actions/workflows/build.yml/badge.svg)

Syndication summarized.

## Dependencies

```
brew install golang goose go-task postgresql redis sqlc

go install github.com/bokwoon95/wgo@latest
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
