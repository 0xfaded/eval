package interactive

import (
	"testing"
)

func TestBasicLiterals(t *testing.T) {
	env := makeEnv()
	expectNil(t, "nil",   env)
	/// expectResult(t, "5",   env, 5)
	/// expectResult(t, "007",   env, 7)
	expectResult(t, "\"abc\"",   env, "abc")
}
