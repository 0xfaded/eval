package interactive

import (
	"testing"
	"reflect"

	"go/parser"
)

func TestCompositeStructValues(t *testing.T) {
	type Alice struct {
		Bob int
	}

	env := addBuiltinTypes(makeEnv())
	env.types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 10 }
	expr := "Alice{ 10 }"

	expectResult(t, expr, env, expected)
}

func expectVoid(t *testing.T, expr string, env *env) {
	expectResults(t, expr, env, []interface{}{})
}

func expectResult(t *testing.T, expr string, env *env, expected interface{}) {
	expectResults(t, expr, env, []interface{}{expected})
}

func expectResults(t *testing.T, expr string, env *env, expected []interface{}) {
	if e, err := parser.ParseExpr(expr); err != nil {
		t.Fatalf("Failed to parse expression '%s' (%v)", expr, err)
	} else if results, _, err := evalExpr(e, env); err != nil {
		t.Fatalf("Error evaluating expression '%s' (%v)", expr, err)
	} else {
		resultsi := make([]interface{}, len(results))
		for i, result := range results {
			resultsi[i] = result.Interface()
		}
		if !reflect.DeepEqual(resultsi, expected) {
			t.Fatalf("Expression '%s' yielded '%+v', expected '%+v'", expr, resultsi, expected)
		}
	}
}

func makeEnv() *env {
	return &env {
		vars: make(map[string] reflect.Value),
		consts: make(map[string] reflect.Value),
		funcs: make(map[string] reflect.Value),
		types: make(map[string] reflect.Type),
		pkgs: make(map[string] pkg),
	}
}
