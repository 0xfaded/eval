// Demos replacing the default identifier lookup value mechanism with
// our own custom version.

package main

import (
	"fmt"
	"reflect"
	"go/parser"
	"github.com/0xfaded/eval"
)

// Here's our custom ident lookup.
func MyEvalIdentExpr(ctx *eval.Ctx, ident *eval.Ident, env *eval.Env) (
	*reflect.Value, bool, error) {
	name := ident.Name
	if name == "nil" {
		return nil, false, nil
	} else if name[0] == 'v' {
		val := reflect.ValueOf(5)
		return &val, true, nil
	} else if name[0] == 'c' {
		val := reflect.ValueOf("constant")
		return &val, true, nil
	} else if name[0] == 'c' {
		val := reflect.ValueOf(true)
		return &val, true, nil
	} else {
		val := reflect.ValueOf('x')
		return &val, true, nil
	}
}


func expectResult(expr string, env *eval.Env, expected interface{}) {
	ctx := &eval.Ctx{expr}
	if e, err := parser.ParseExpr(expr); err != nil {
		fmt.Printf("Failed to parse expression '%s' (%v)\n", expr, err)
		return
	} else if cexpr, errs := eval.CheckExpr(ctx, e, env); len(errs) != 0 {
		fmt.Printf("Error checking expression '%s' (%v)\n", expr, errs)
	} else if results, _, err := eval.EvalExpr(ctx, cexpr, env); err != nil {
		fmt.Printf("Error evaluating expression '%s' (%v)\n", expr, err)
		return
	} else {
		fmt.Printf("Expression '%s' yielded '%+v', expected '%+v'\n",
			expr, (*results)[0].Interface(), expected)
	}
}

func makeEnv() *eval.Env {
	return &eval.Env {
		Vars: make(map[string] reflect.Value),
		Consts: make(map[string] reflect.Value),
		Funcs: make(map[string] reflect.Value),
		Types: make(map[string] reflect.Type),
		Pkgs: make(map[string] eval.Pkg),
	}
}

func main() {
	env := makeEnv()
	eval.SetEvalIdentExprCallback(MyEvalIdentExpr)
	expectResult("v + 1", env, "6")
	expectResult("c + \" value\"", env, "constant value")

}
