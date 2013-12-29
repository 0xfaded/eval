package eval_test

import (
	"fmt"
	"reflect"
	"go/parser"
	"github.com/0xfaded/eval"
)

const constant1 = "A constant"

//MakeEnv creates an environment to use in eval
func makeEnv() *eval.Env {
	env := &eval.Env {
		Vars: make(map[string] reflect.Value),
		Consts: make(map[string] reflect.Value),
		Funcs: make(map[string] reflect.Value),
		Types: make(map[string] reflect.Type),
		Pkgs: make(map[string] eval.Pkg),
	}
	env.Consts["constant1"] = reflect.ValueOf(constant1)
	var1 := 1
	env.Vars["var1"] = reflect.ValueOf(&var1)
	env.Funcs["expectResult"] = reflect.ValueOf(ExpectResult)
	return env

}

// ExpectResult check the evaluation of a string with an expected result.
// More importantly though, does the steps to evaluate a string:
//   0. Create an evaluation enviroment
//   1. Parse expression using parser.ParseExpr (go/parser)
//   2. Type check expression using evalCheckExpr (0xfaded/eval)
//   3. run eval.EvalExpr (0xfaded/eval)
func ExpectResult(expr string, expected interface{}) {
	env := makeEnv() // Create evaluation environment
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

func Example() {
	ExpectResult("\"A constant!\"!", constant1 + "!")
}
