# 📚 verso

This project was a Google Reader API clone based on [unofficial documentation](https://github.com/mihaip/google-reader-api). It is no longer updated.

The API is about 80% functional and usable in third-party apps like NetNewsWire or Reeder.

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
