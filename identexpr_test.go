package interactive

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
	env.Vars["arg0"] = reflect.ValueOf("abc")
	expectResult(t, "arg0", env, "abc")
}
