# citra

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
task goose:create
task goose:drop
task goose:gen -- migration_name
task goose:migrate
task goose:rollback
```
