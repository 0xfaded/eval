package interactive

import (
	"testing"
)

func TestNilValues(t *testing.T) {
	env := makeEnv()
	expectNil(t, "nil", env)
}
