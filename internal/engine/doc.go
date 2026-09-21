// Package engine implements the workflow scheduler: a poll/queue driven
// worker pool that claims task_queue rows, executes the corresponding
// activity, and persists the resulting state transition back to Postgres.
//
// Implemented in Epic 7 of the Phase 1 task list.
package engine
