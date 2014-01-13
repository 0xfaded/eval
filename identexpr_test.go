package eval

import (
	"reflect"
	"testing"
)

func TestNilValue(t *testing.T) {
	env := makeEnv()
	expectNil(t, "nil", env)
}

func TestStringVar(t *testing.T) {
	env := makeEnv()
        s := "abc"
	env.Vars["arg0"] = reflect.ValueOf(&s)
	expectResult(t, "arg0", env, "abc")
}
