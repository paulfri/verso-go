# citra

Syndication summarized.

## Dependencies

```
brew install golang goose go-task sqlc

go install github.com/cosmtrek/air@latest
```

## Usage

```
task --list-all

# Auto-reloading in development:
task server
task worker
task --parallel server worker
task -p s w
```

## Migrations

```
task goose:create
task goose:drop
task goose:gen -- migration_name
task goose:migrate
task goose:rollback
```
