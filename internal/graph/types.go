package graph

import "encoding/json"

// VariableDecl declares a workflow-level variable available in the
// instance's execution context.
type VariableDecl struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

// Node is a single step in the workflow graph.
type Node struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config,omitempty"`

	// Next lists unconditional successor node ids. For composite nodes such
	// as builtin.if, Next is unused in favor of Branches; for builtin.loop it
	// represents the "after loop" (condition-false) successor(s).
	Next []string `json:"next,omitempty"`

	// Branches maps a named outcome (e.g. "true"/"false") to successor node
	// ids. Used by conditional/composite nodes.
	Branches map[string][]string `json:"branches,omitempty"`
}

// Successors returns every node id that n can transition to, across both
// Next and Branches.
func (n Node) Successors() []string {
	out := make([]string, 0, len(n.Next))
	out = append(out, n.Next...)
	for _, ids := range n.Branches {
		out = append(out, ids...)
	}
	return out
}

// Definition is a versioned, named workflow graph.
type Definition struct {
	ID        string         `json:"id"`
	Version   int            `json:"version"`
	Name      string         `json:"name,omitempty"`
	Variables []VariableDecl `json:"variables,omitempty"`
	Nodes     []Node         `json:"nodes"`
}

// NodeByID returns the node with the given id, or false if not found.
func (d *Definition) NodeByID(id string) (Node, bool) {
	for _, n := range d.Nodes {
		if n.ID == id {
			return n, true
		}
	}
	return Node{}, false
}

// Builtin activity type names usable in a Node's Type field.
const (
	TypeStart     = "builtin.start"
	TypeStop      = "builtin.stop"
	TypeIf        = "builtin.if"
	TypeLoop      = "builtin.loop"
	TypeParallel  = "builtin.parallel"
	TypeJoin      = "builtin.join"
	TypeRetry     = "builtin.retry"
	TypeDelay     = "builtin.delay"
	TypePersist   = "builtin.persist"
	TypeTerminate = "builtin.terminate"
	TypeCancel    = "builtin.cancel"
)
