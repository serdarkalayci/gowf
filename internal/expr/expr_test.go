package expr

import (
	"strings"
	"testing"
)

func TestEvalBool(t *testing.T) {
	vars := map[string]any{"stock": 5, "name": "widget", "active": true}

	tests := []struct {
		name string
		expr string
		want bool
	}{
		{"greater than true", "stock > 0", true},
		{"greater than false", "stock > 10", false},
		{"string equality", `name == "widget"`, true},
		{"bool variable", "active", true},
		{"logical and", "stock > 0 && active", true},
		{"logical not", "!active", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalBool(tt.expr, vars)
			if err != nil {
				t.Fatalf("EvalBool(%q) unexpected error: %v", tt.expr, err)
			}
			if got != tt.want {
				t.Fatalf("EvalBool(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestEvalBool_NonBoolResult(t *testing.T) {
	_, err := EvalBool("stock + 1", map[string]any{"stock": 5})
	if err == nil || !strings.Contains(err.Error(), "did not evaluate to a bool") {
		t.Fatalf("EvalBool() = %v, want non-bool error", err)
	}
}

func TestEvalBool_CompileError(t *testing.T) {
	_, err := EvalBool("stock >", map[string]any{"stock": 5})
	if err == nil || !strings.Contains(err.Error(), "compiling expression") {
		t.Fatalf("EvalBool() = %v, want compile error", err)
	}
}

func TestEvalBool_UnknownVariable(t *testing.T) {
	_, err := EvalBool("missing > 0", map[string]any{"stock": 5})
	if err == nil {
		t.Fatal("EvalBool() expected error for unknown variable, got nil")
	}
}

func TestEval_ReturnsRawValue(t *testing.T) {
	got, err := Eval("stock * 2", map[string]any{"stock": 5})
	if err != nil {
		t.Fatalf("Eval() unexpected error: %v", err)
	}
	if got != 10 {
		t.Fatalf("Eval() = %v (%T), want 10", got, got)
	}
}

func TestEval_RuntimeError(t *testing.T) {
	_, err := Eval("items[5]", map[string]any{"items": []int{1}})
	if err == nil {
		t.Fatal("Eval() expected runtime error for out-of-range index, got nil")
	}
	if !strings.Contains(err.Error(), "evaluating expression") {
		t.Fatalf("Eval() = %v, want runtime evaluation error", err)
	}
}

func TestValidate(t *testing.T) {
	if err := Validate("stock > 0 && name == \"x\""); err != nil {
		t.Fatalf("Validate() unexpected error for valid expression: %v", err)
	}

	err := Validate("stock >")
	if err == nil || !strings.Contains(err.Error(), "invalid expression") {
		t.Fatalf("Validate() = %v, want invalid expression error", err)
	}
}
