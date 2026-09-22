package graph

import (
	"bytes"
	"encoding/json"
	"testing"
)

// mustMarshal marshals v with HTML-escaping disabled so structurally
// identical JSON (e.g. differing only in whether '>' was escaped) compares
// equal as a string.
func mustMarshal(t *testing.T, v any) string {
	t.Helper()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return buf.String()
}

func sampleDefinition() *Definition {
	return &Definition{
		ID:      "order-fulfillment",
		Version: 3,
		Name:    "Order Fulfillment",
		Variables: []VariableDecl{
			{Name: "orderId", Type: "string"},
			{Name: "stock", Type: "int"},
		},
		Nodes: []Node{
			{ID: "start", Type: TypeStart, Next: []string{"checkStock"}},
			{
				ID:     "checkStock",
				Type:   TypeIf,
				Config: json.RawMessage(`{"expr":"stock > 0"}`),
				Branches: map[string][]string{
					"true":  {"reserve"},
					"false": {"backorder"},
				},
			},
			{ID: "reserve", Type: TypePersist, Next: []string{"end"}},
			{ID: "backorder", Type: TypePersist, Next: []string{"end"}},
			{ID: "end", Type: TypeStop},
		},
	}
}

func TestDefinitionJSONRoundTrip(t *testing.T) {
	original := sampleDefinition()

	// Use the no-HTML-escaping encoder for the initial marshal too: the
	// standard json.Marshal HTML-escapes characters like '>' when writing
	// json.RawMessage bytes into the surrounding document, which would
	// otherwise get permanently baked into roundTripped.Config as escaped
	// text and make a later byte/string comparison spuriously fail even
	// though the JSON is semantically identical.
	originalJSON := mustMarshal(t, original)

	var roundTripped Definition
	if err := json.Unmarshal([]byte(originalJSON), &roundTripped); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	roundTrippedJSON := mustMarshal(t, &roundTripped)
	if originalJSON != roundTrippedJSON {
		t.Fatalf("round-tripped definition does not match original.\noriginal:  %s\nroundtrip: %s", originalJSON, roundTrippedJSON)
	}
}

func TestNodeSuccessors(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want []string
	}{
		{
			name: "next only",
			node: Node{ID: "a", Next: []string{"b", "c"}},
			want: []string{"b", "c"},
		},
		{
			name: "branches only",
			node: Node{ID: "a", Branches: map[string][]string{"true": {"b"}}},
			want: []string{"b"},
		},
		{
			name: "no successors",
			node: Node{ID: "a"},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.Successors()
			if len(got) != len(tt.want) {
				t.Fatalf("Successors() = %v, want %v", got, tt.want)
			}
			wantSet := make(map[string]bool, len(tt.want))
			for _, w := range tt.want {
				wantSet[w] = true
			}
			for _, g := range got {
				if !wantSet[g] {
					t.Fatalf("Successors() returned unexpected id %q; got %v want %v", g, got, tt.want)
				}
			}
		})
	}
}

func TestDefinitionNodeByID(t *testing.T) {
	def := sampleDefinition()

	if n, ok := def.NodeByID("checkStock"); !ok || n.Type != TypeIf {
		t.Fatalf("NodeByID(%q) = %+v, %v; want the checkStock node", "checkStock", n, ok)
	}

	if _, ok := def.NodeByID("does-not-exist"); ok {
		t.Fatalf("NodeByID(%q) unexpectedly found a node", "does-not-exist")
	}
}
