package graph

import (
	"strings"
	"testing"
)

func TestValidate_ValidDefinition(t *testing.T) {
	if err := Validate(sampleDefinition()); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestValidate_ValidDefinitionWithLoopBackEdge(t *testing.T) {
	def := &Definition{
		ID: "loopy",
		Nodes: []Node{
			{ID: "start", Type: TypeStart, Next: []string{"loop"}},
			{
				ID:   "loop",
				Type: TypeLoop,
				// "true" branch loops back to itself (cycle); "false"
				// branch exits to stop. The cycle must not defeat
				// reachability BFS.
				Branches: map[string][]string{
					"true":  {"loop"},
					"false": {"end"},
				},
			},
			{ID: "end", Type: TypeStop},
		},
	}
	if err := Validate(def); err != nil {
		t.Fatalf("Validate() unexpected error for cyclic definition: %v", err)
	}
}

func TestValidate_EmptyID(t *testing.T) {
	def := &Definition{Nodes: []Node{{ID: "a", Type: TypeStart}}}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), "definition id is required") {
		t.Fatalf("Validate() = %v, want error about missing definition id", err)
	}
}

func TestValidate_NoNodes(t *testing.T) {
	def := &Definition{ID: "empty"}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), "has no nodes") {
		t.Fatalf("Validate() = %v, want error about no nodes", err)
	}
}

func TestValidate_NodeWithEmptyID(t *testing.T) {
	def := &Definition{
		ID: "bad",
		Nodes: []Node{
			{ID: "start", Type: TypeStart, Next: []string{""}},
			{ID: "", Type: TypeStop},
		},
	}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), "empty id") {
		t.Fatalf("Validate() = %v, want error about empty node id", err)
	}
}

func TestValidate_DuplicateNodeID(t *testing.T) {
	def := &Definition{
		ID: "dup",
		Nodes: []Node{
			{ID: "start", Type: TypeStart, Next: []string{"end"}},
			{ID: "end", Type: TypeStop},
			{ID: "end", Type: TypeStop},
		},
	}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), `duplicate node id "end"`) {
		t.Fatalf("Validate() = %v, want error about duplicate node id", err)
	}
}

func TestValidate_DanglingSuccessor(t *testing.T) {
	def := &Definition{
		ID: "dangling",
		Nodes: []Node{
			{ID: "start", Type: TypeStart, Next: []string{"nowhere"}},
			{ID: "end", Type: TypeStop},
		},
	}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), `unknown successor "nowhere"`) {
		t.Fatalf("Validate() = %v, want error about unknown successor", err)
	}
}

func TestValidate_DanglingBranchSuccessor(t *testing.T) {
	def := &Definition{
		ID: "dangling-branch",
		Nodes: []Node{
			{ID: "start", Type: TypeStart, Next: []string{"decide"}},
			{
				ID:   "decide",
				Type: TypeIf,
				Branches: map[string][]string{
					"true":  {"end"},
					"false": {"missing"},
				},
			},
			{ID: "end", Type: TypeStop},
		},
	}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), `unknown successor "missing"`) {
		t.Fatalf("Validate() = %v, want error about unknown branch successor", err)
	}
}

func TestValidate_NoStartNode(t *testing.T) {
	def := &Definition{
		ID: "no-start",
		Nodes: []Node{
			{ID: "end", Type: TypeStop},
		},
	}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), "must have exactly one builtin.start node, found 0") {
		t.Fatalf("Validate() = %v, want error about missing start node", err)
	}
}

func TestValidate_MultipleStartNodes(t *testing.T) {
	def := &Definition{
		ID: "two-starts",
		Nodes: []Node{
			{ID: "start1", Type: TypeStart, Next: []string{"end"}},
			{ID: "start2", Type: TypeStart, Next: []string{"end"}},
			{ID: "end", Type: TypeStop},
		},
	}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), "must have exactly one builtin.start node, found 2") {
		t.Fatalf("Validate() = %v, want error about multiple start nodes", err)
	}
}

func TestValidate_NoStopNode(t *testing.T) {
	def := &Definition{
		ID: "no-stop",
		Nodes: []Node{
			{ID: "start", Type: TypeStart},
		},
	}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), "must have at least one builtin.stop node") {
		t.Fatalf("Validate() = %v, want error about missing stop node", err)
	}
}

func TestValidate_StopNodeUnreachable(t *testing.T) {
	def := &Definition{
		ID: "unreachable-stop",
		Nodes: []Node{
			{ID: "start", Type: TypeStart, Next: []string{"loop"}},
			// loop only ever points back to itself, so "end" is never
			// reached from "start".
			{ID: "loop", Type: TypeLoop, Next: []string{"loop"}},
			{ID: "end", Type: TypeStop},
		},
	}
	err := Validate(def)
	if err == nil || !strings.Contains(err.Error(), "no builtin.stop node is reachable from builtin.start") {
		t.Fatalf("Validate() = %v, want error about unreachable stop node", err)
	}
}

func TestValidate_TerminateCountsAsStop(t *testing.T) {
	def := &Definition{
		ID: "terminate-only",
		Nodes: []Node{
			{ID: "start", Type: TypeStart, Next: []string{"end"}},
			{ID: "end", Type: TypeTerminate},
		},
	}
	if err := Validate(def); err != nil {
		t.Fatalf("Validate() unexpected error when only a terminate node exists: %v", err)
	}
}
