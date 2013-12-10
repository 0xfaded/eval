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
	expectResult(t, "\"a \\\"fixed\\\" bug\"", env, "a \"fixed\" bug")
	// There is a tab in the below expect string...
	expectResult(t, "\"a\tbc\"", env, "a	bc")
}
