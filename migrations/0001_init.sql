-- Phase 1 core schema for gowf. Idempotent: safe to apply repeatedly.

CREATE TABLE IF NOT EXISTS workflow_definitions (
    id         TEXT        NOT NULL,
    version    INT         NOT NULL,
    name       TEXT        NOT NULL DEFAULT '',
    spec       JSONB       NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, version)
);

CREATE TABLE IF NOT EXISTS workflow_instances (
    instance_id        UUID        PRIMARY KEY,
    definition_id      TEXT        NOT NULL,
    definition_version INT         NOT NULL,
    status             TEXT        NOT NULL
        CHECK (status IN ('running', 'waiting', 'completed', 'faulted', 'cancelled')),
    variables          JSONB       NOT NULL DEFAULT '{}'::jsonb,
    current_nodes      JSONB       NOT NULL DEFAULT '[]'::jsonb,
    fault_message      TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (definition_id, definition_version)
        REFERENCES workflow_definitions (id, version)
);
CREATE INDEX IF NOT EXISTS idx_workflow_instances_definition
    ON workflow_instances (definition_id, definition_version);
CREATE INDEX IF NOT EXISTS idx_workflow_instances_status
    ON workflow_instances (status);

CREATE TABLE IF NOT EXISTS task_queue (
    task_id      UUID        PRIMARY KEY,
    instance_id  UUID        NOT NULL REFERENCES workflow_instances (instance_id) ON DELETE CASCADE,
    node_id      TEXT        NOT NULL,
    payload      JSONB       NOT NULL DEFAULT '{}'::jsonb,
    available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    claimed_by   TEXT,
    claimed_at   TIMESTAMPTZ,
    attempts     INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- Partial index serving the worker claim query (unclaimed, due tasks).
CREATE INDEX IF NOT EXISTS idx_task_queue_claimable
    ON task_queue (available_at) WHERE claimed_by IS NULL;
-- Serves the stale-claim reclaimer.
CREATE INDEX IF NOT EXISTS idx_task_queue_claimed
    ON task_queue (claimed_at) WHERE claimed_by IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_task_queue_instance
    ON task_queue (instance_id);

CREATE TABLE IF NOT EXISTS bookmarks (
    bookmark_id UUID        PRIMARY KEY,
    instance_id UUID        NOT NULL REFERENCES workflow_instances (instance_id) ON DELETE CASCADE,
    node_id     TEXT        NOT NULL,
    name        TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    resumed_at  TIMESTAMPTZ
);
-- Serves bookmark lookup by (instance, name) among unresumed bookmarks.
CREATE INDEX IF NOT EXISTS idx_bookmarks_pending
    ON bookmarks (instance_id, name) WHERE resumed_at IS NULL;

CREATE TABLE IF NOT EXISTS timers (
    timer_id    UUID        PRIMARY KEY,
    instance_id UUID        NOT NULL REFERENCES workflow_instances (instance_id) ON DELETE CASCADE,
    node_id     TEXT        NOT NULL,
    fire_at     TIMESTAMPTZ NOT NULL,
    resumed_at  TIMESTAMPTZ
);
-- Serves the timer poller (unresumed timers ordered by fire time).
CREATE INDEX IF NOT EXISTS idx_timers_pending
    ON timers (fire_at) WHERE resumed_at IS NULL;

CREATE TABLE IF NOT EXISTS instance_history (
    id          BIGSERIAL   PRIMARY KEY,
    instance_id UUID        NOT NULL REFERENCES workflow_instances (instance_id) ON DELETE CASCADE,
    node_id     TEXT,
    event_type  TEXT        NOT NULL,
    data        JSONB,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_instance_history_instance
    ON instance_history (instance_id, id);
