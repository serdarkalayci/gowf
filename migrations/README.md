# migrations

SQL migration files for the gowf Postgres schema.

Files are applied in lexical filename order by `internal/persistence`.
Each migration must be idempotent (e.g. `CREATE TABLE IF NOT EXISTS`) since
there is no version-tracking table in Phase 1.

## Files

| File | Contents |
|------|----------|
| `0001_init.sql` | `workflow_definitions`, `workflow_instances`, `task_queue`, `bookmarks`, `timers`, `instance_history` plus supporting indexes |

## Verifying a migration locally

Start the disposable test database (see `build/testdb/Containerfile`), apply
the file twice, then tear it down; both applies must succeed:

```sh
task testdb:up
podman exec -i gowf-testdb psql -U gowf -d gowf_test -v ON_ERROR_STOP=1 < migrations/0001_init.sql
podman exec -i gowf-testdb psql -U gowf -d gowf_test -v ON_ERROR_STOP=1 < migrations/0001_init.sql
task testdb:down
```

(Use `docker` instead of `podman` if that is your engine.)
