package interactive

import (
	"testing"
	"reflect"
)

func TestCompositeStructValues(t *testing.T) {
	type Alice struct {
		Bob int
	}

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 10 }
	expr := "Alice{ 10 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeStructKeyValues(t *testing.T) {
	type Alice struct {
		Bob int
	}

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { Bob: 10 }
	expr := "Alice{ Bob: 10 }"

	expectResult(t, expr, env, expected)
}

