package eval

import (
	"testing"
	"github.com/0xfaded/reflectext"
)

func TestEmptyFunc(t *testing.T) {
	if !reflectext.Available {
		return
	}
	env := MakeSimpleEnv()
	fn := "func() {}()"
	expectResults(t, fn,  env)
}

func TestSimpleFunc(t *testing.T) {
	if !reflectext.Available {
		return
	}
	env := MakeSimpleEnv()
	fn := "func() int { return 1 }()"
	expectResult(t, fn,  env, int(1))
}

func TestMultiFunc(t *testing.T) {
	if !reflectext.Available {
		return
	}
	env := MakeSimpleEnv()
	fn := "func() (int, int) { return 1, 2 }()"
	expectResults(t, fn,  env, int(1), int(2))
}

func TestScopingFunc(t *testing.T) {
	if !reflectext.Available {
		return
	}
	env := MakeSimpleEnv()
	fn := "func() int { { return 1 } }()"
	expectResults(t, fn,  env, int(1))
}

func TestMultiUnpackFunc(t *testing.T) {
	if !reflectext.Available {
		return
	}
	env := MakeSimpleEnv()
	expectInterp(t, "a := func() (int, int) { return 1, 2 }", env)
	fn := "func() (int, int) { return a() }()"
	expectResults(t, fn,  env, int(1), int(2))
}

