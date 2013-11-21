package interactive

// Utilities for other tests live here

import (
	"testing"
	"reflect"

	"go/parser"
)

func expectVoid(t *testing.T, expr string, env *Env) {
	expectResults(t, expr, env, []interface{}{})
}

func expectResult(t *testing.T, expr string, env *Env, expected interface{}) {
	expectResults(t, expr, env, []interface{}{expected})
}

func expectResults(t *testing.T, expr string, env *Env, expected []interface{}) {
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

func makeEnv() *Env {
	return &Env {
		Vars: make(map[string] reflect.Value),
		Consts: make(map[string] reflect.Value),
		Funcs: make(map[string] reflect.Value),
		Types: make(map[string] reflect.Type),
		Pkgs: make(map[string] Pkg),
	}
}
