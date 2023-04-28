# Migrations

Use [goose](https://github.com/pressly/goose) to manage database migrations.

See `task --list-all` for commands.

## Conventions

1. Table and column names are Rails style: snake_case, inflected, etc.
2. The first four columns of all tables are:

```sql
id bigint primary key generated always as identity,
uuid uuid unique not null default gen_random_uuid(),
created_at timestamptz not null default now(),
updated_at timestamptz not null default now(),
```

3. Indexes are named `{tablename}_{columnname(s)}_{suffix}`, where `suffix` is
   in `pkey`, `key` (unique), `excl`, `idx`, `fkey` or `check`.
4. Add the `touch_updated_at()` trigger to every table.
5. Organize tables by schema.
6. Cascade deletes.
