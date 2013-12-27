package main

import (
	"fmt"
	"reflect"
	"go/parser"
	"github.com/0xfaded/eval"
)

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

func TestBinaryOps() {
	env := makeEnv()

	expectResult("\"a\" + \"b\" == \"ab\"", env, "a" + "b" == "ab")
	expectResult("\"a\" + \"b\"", env, "a" + "b")
}

func TestUintBinaryOps() {
	env := makeEnv()

	expectResult("uint64(1)+2", env, uint64(1)+2)
}

func TestUnaryOps() {
	env := makeEnv()

	expectResult("-1.0", env, float64(-1.0))
}

func TestComplexOps() {
	env := makeEnv()
	expectResult("complex(1, 2) + complex(3, 4)", env, complex(1, 2) + complex(3, 4))
}

func TestTypedBinaryOps() {

	type Foo int

	env := makeEnv()
	env.Types["Foo"] = reflect.TypeOf(Foo(0))

	expectResult("Foo(1)+Foo(2)", env, Foo(1)+Foo(2))
	expectResult("1-Foo(2)", env, 1-Foo(2))
	expectResult("Foo(1)|2", env, Foo(1)|2)
}

func TestNil() {
	env := makeEnv()

	expectResult("nil", env, nil)
}

func main() {
	TestBinaryOps()
	TestUnaryOps()
	TestTypedBinaryOps()
}
