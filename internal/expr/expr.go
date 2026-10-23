package expr

import (
	"fmt"

	"github.com/expr-lang/expr"
)

// Eval compiles and runs exprStr against vars, returning the raw result.
func Eval(exprStr string, vars map[string]any) (any, error) {
	program, err := expr.Compile(exprStr, expr.Env(vars))
	if err != nil {
		return nil, fmt.Errorf("compiling expression %q: %w", exprStr, err)
	}
	out, err := expr.Run(program, vars)
	if err != nil {
		return nil, fmt.Errorf("evaluating expression %q: %w", exprStr, err)
	}
	return out, nil
}

// EvalBool evaluates exprStr and requires the result to be a bool.
func EvalBool(exprStr string, vars map[string]any) (bool, error) {
	out, err := Eval(exprStr, vars)
	if err != nil {
		return false, err
	}
	b, ok := out.(bool)
	if !ok {
		return false, fmt.Errorf("expression %q did not evaluate to a bool (got %T)", exprStr, out)
	}
	return b, nil
}

// Validate checks that exprStr is syntactically valid without running it,
// so errors can be caught at definition-registration time. Variable names
// are not checked since the runtime variable set is unknown here.
func Validate(exprStr string) error {
	if _, err := expr.Compile(exprStr); err != nil {
		return fmt.Errorf("invalid expression %q: %w", exprStr, err)
	}
	return nil
}
