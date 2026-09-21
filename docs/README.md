# docs

## What is gowf?

`gowf` is a Go-native workflow platform inspired by .NET's Windows Workflow
Foundation 4 (WF4). It provides:

- **An engine** that runs, schedules, and persists workflow instances. The
  engine is poll/queue-driven (workers claim rows from a Postgres-backed
  task queue) rather than goroutine-per-instance, which means an instance's
  state is never held only in memory — "dehydration" and "rehydration" are
  implicit, and a crashed or restarted engine process can resume every
  in-flight instance from durable storage.
- **A built-in activity library** covering the essential control-flow and
  reliability primitives every workflow needs: start/stop, if/then/else,
  loop, parallel/join, retry, delay, and explicit persist checkpoints.
- **A plugin system** (planned, Phase 2) for user-defined custom
  activities, distributed as OCI images containing a native binary plus a
  `.proto` describing the activities it exposes. At runtime the engine
  extracts and runs the binary as a local subprocess and talks to it over
  gRPC, without depending on a container runtime.
- **A web-based drag & drop designer** (planned, Phase 3) for building
  workflows visually. Workflow definitions are a JSON flow graph (nodes +
  edges) — the same shape the designer would naturally produce — rather
  than a nested activity tree.

A single engine instance supports multiple workflow *definitions*
(types) simultaneously, as well as many concurrent *instances* of each
definition, all coordinated through the same Postgres-backed task queue,
bookmark, and timer tables.

## Status

The project is in **Phase 1**: building the core engine, built-in
activities, and persistence layer (no plugin system or designer UI yet).
See the full project plan below for the current breakdown of epics and
issues.

## Project Plan

The project is delivered in three phases. Phase 1 is broken down into
epics and issues small enough to finish, test, and commit individually;
Phases 2 and 3 are scoped at a design level for now and will get their own
detailed breakdowns when started.

### Phase 1 — Core engine, built-in activities, persistence

**Epic 1: Project scaffolding**
- 1.1 Initialize Go module and directory layout

**Epic 2: Workflow definition model (graph package)**
- 2.1 Define Node/Definition JSON types
- 2.2 Implement structural definition validation

**Epic 3: Expression evaluation**
- 3.1 Wrap expr-lang for condition evaluation

**Epic 4: Persistence layer (Postgres)**
- 4.1 Write Postgres schema migration SQL
- 4.2 Store connection + migration runner
- 4.3 Definitions repository
- 4.4 Instances repository
- 4.5 Task queue repository
- 4.6 Bookmarks and timers repository
- 4.7 History repository

**Epic 5: Activity interface and registry**
- 5.1 Define Activity, Context, and Result types
- 5.2 Implement the Activity registry

**Epic 6: Built-in activity library**
- 6.1 start/stop activities
- 6.2 if/then/else activity
- 6.3 loop activity
- 6.4 parallel/join activities
- 6.5 retry activity
- 6.6 delay activity
- 6.7 persist (checkpoint) activity
- 6.8 terminate/cancel activities

**Epic 7: Engine core (scheduler)**
- 7.1 Implement engine execContext (activity.Context)
- 7.2 Implement single-task execution (runTask)
- 7.3 Implement worker pool loop
- 7.4 Implement timer poller loop
- 7.5 Implement stale task reclaimer
- 7.6 Implement StartInstance / ResumeBookmark / CancelInstance
- 7.7 End-to-end engine scenario tests (Phase 1 acceptance gate)

**Epic 8: REST + WebSocket/SSE API**
- 8.1 API server bootstrap
- 8.2 Definitions endpoints
- 8.3 Instance lifecycle endpoints
- 8.4 Bookmark resume endpoint
- 8.5 SSE/WS event streaming endpoint

### Phase 2 — Plugin system (OCI-based custom activities)

Design-level scope, detailed breakdown to follow once Phase 1 is complete:

- Define the plugin `.proto` contract (an Activity gRPC service exposing at
  least `Describe` and `Execute`, with `Cancel` under consideration).
- Define the OCI image layout convention bundling a native binary plus its
  `.proto` file.
- `gowfctl plugin pull/inspect`: pull an OCI image, extract the binary and
  proto, and catalog the activities it exposes (name, inputs, outputs) into
  the activity registry.
- Process lifecycle management for extracted plugin binaries: spawn as a
  local subprocess (go-plugin-style handshake over a local socket/stdio, no
  container runtime at execution time), health-check, restart, and version
  pinning.
- Async/bookmark-based completion for long-running plugin activities,
  matching the built-in bookmark mechanism from Phase 1.
- Sandboxing/resource limits for plugin subprocesses without depending on a
  container runtime.

### Phase 3 — Web-based drag & drop designer

Design-level scope, detailed breakdown to follow once Phase 2 is complete:

- A web application (likely React + a graph/flow library) that edits the
  same JSON flow graph the engine consumes as a workflow definition.
- A "node palette" driven by the activity registry, so both built-in
  activities and registered plugin activities (from Phase 2) show up as
  draggable nodes automatically.
- Visual editing of composite node configuration (if/loop/retry/parallel
  branches, conditions, and nested subgraphs).
- Integration with the Phase 1 REST/SSE API to register definitions, start
  instances, and visualize live instance status/history.

### Explicitly deferred (not yet designed)

- Compensation/saga-style rollback (WF4 had this).
- Multi-node/distributed engine execution (clustering/HA). The Phase 1
  design keeps all state in Postgres so this can be layered on later
  without a rewrite; it will need a leader-elected timer poller to avoid
  duplicate firing across engine instances.

## Contents

Project documentation (architecture notes, API references, etc.) will be
added here as Phase 1 progresses.
