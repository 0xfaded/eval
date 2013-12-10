package interactive

import (
	"os"
	"log"
	"testing"
	"reflect"
)

func TestFuncCallWithConst(t *testing.T) {
	env := makeEnv()
	env.Consts["X"] = reflect.ValueOf(int64(10))
	env.Funcs["Foo"] = reflect.ValueOf(func (int) int { return 1; })

	expectResult(t, "Foo(X)", env, 1)
}

func TestFuncCallWithWrongArgs(t *testing.T) {
	env := makeEnv()
	env.Funcs["Foo"] = reflect.ValueOf(func (string) int { return 1; })

	expectFail(t, "Foo(1.5)", env)
}

func TestFuncCallLogNewWithWrongArgs(t *testing.T) {
	env := makeEnv()
	logpkg := makeEnv()
	ospkg := makeEnv()

	env.Pkgs["log"] = logpkg
	env.Pkgs["os"] = ospkg

	logpkg.Funcs["New"] = reflect.ValueOf(log.New)
	ospkg.Vars["Stdout"] = reflect.ValueOf(&os.Stdout)

	expectFail(t, "log.New(\"Bob\"), os.Stdout, 0)", env)
}

func TestFuncCallWithSplatOne(t *testing.T) {
	env := makeEnv()

	f := func() int { return 1 }
	g := func(a int) int { return a }

	env.Funcs["f"] = reflect.ValueOf(f)
	env.Funcs["g"] = reflect.ValueOf(g)

	expr := "g(f())"
	expected := g(f())

	expectResult(t, expr, env, expected)
}

func TestFuncCallWithSplatTwo(t *testing.T) {
	env := makeEnv()

	f := func() (int, int) { return 1, 2 }
	g := func(a int, b int) int { return a + b }

	env.Funcs["f"] = reflect.ValueOf(f)
	env.Funcs["g"] = reflect.ValueOf(g)

	expr := "g(f())"
	expected := g(f())

	expectResult(t, expr, env, expected)
}

// This test hits a specific case in the implementation where
// f(g()) is evaluated as args := g(); f(args)
func TestFuncCallWithMissingValueSplat(t *testing.T) {
	env := makeEnv()

	env.Funcs["f"] = reflect.ValueOf(func() {})
	env.Funcs["g"] = reflect.ValueOf(func(int) {})

	expectError(t, "g(f())", env, "f() used as value")
}

func TestFuncCallWithMissingValue(t *testing.T) {
	env := makeEnv()

	env.Funcs["f"] = reflect.ValueOf(func() {})
	env.Funcs["g"] = reflect.ValueOf(func(int, int) {})

	expectError(t, "g(1, f())", env, "f() used as value")
}

