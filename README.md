# citra

Syndication summarized.

## Dependencies

```
brew install goose go-task sqlc
```

## Server

```
task server
```

## Tasks

```
task --list-all
```

## Migrations

```
task goose:create
task goose:drop
task goose:gen -- migration_name
task goose:migrate
task goose:rollback
```
