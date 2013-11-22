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

