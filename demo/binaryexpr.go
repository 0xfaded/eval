package main

import (
	"fmt"
	"reflect"
	"go/parser"
	"github.com/rocky/go-interactive"
)

func expectResult(expr string, env *interactive.Env, expected interface{}) {
	if e, err := parser.ParseExpr(expr); err != nil {
		fmt.Printf("Failed to parse expression '%s' (%v)\n", expr, err)
		return
	} else if results, _, err := interactive.EvalExpr(e, env); err != nil {
		fmt.Printf("Error evaluating expression '%s' (%v)\n", expr, err)
		return
	} else {
		fmt.Printf("Expression '%s' yielded '%+v', expected '%+v'\n",
			expr, (*results)[0].Interface(), expected)
	}
}

func makeEnv() *interactive.Env {
	return &interactive.Env {
		Vars: make(map[string] reflect.Value),
		Consts: make(map[string] reflect.Value),
		Funcs: make(map[string] reflect.Value),
		Types: make(map[string] reflect.Type),
		Pkgs: make(map[string] interactive.Pkg),
	}
}

func TestIntBinaryOps() {
	env := makeEnv()

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
	TestIntBinaryOps()
	TestUnaryOps()
	TestTypedBinaryOps()
}
