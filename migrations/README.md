# migrations

SQL migration files for the gowf Postgres schema.

Files are applied in lexical filename order by `internal/persistence`.
Each migration must be idempotent (e.g. `CREATE TABLE IF NOT EXISTS`) since
there is no version-tracking table in Phase 1.

The first migration (`workflow_definitions`, `workflow_instances`,
`task_queue`, `bookmarks`, `timers`, `instance_history`) will be added in
Issue 4.1 of the Phase 1 task list.
