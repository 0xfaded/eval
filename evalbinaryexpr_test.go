package eval

import (
	"reflect"
	"testing"
)

func TestIntBinaryOps(t *testing.T) {
	x := int32(5)
	env := makeEnv()
	env.Vars["x"] = reflect.ValueOf(&x)

	expectResult(t, "x+2",   env, x+2)
	expectResult(t, "x-2",   env, x-2)
	expectResult(t, "x*3.0",   env, x*3.0)
	expectResult(t, "x/2",   env, x/2)
	expectResult(t, "x%2",   env, x%2)
	expectResult(t, "x&1",   env, x&1)
	expectResult(t, "x|1",   env, x|1)
	expectResult(t, "x^1",   env, x^1)
	expectResult(t, "x&^1",  env, x&^1)
	expectResult(t, "x+x",  env, x+x)

	expectResult(t, "x<1",   env, bool(x<1))
	expectResult(t, "-x<=3",  env, bool(-x<=3))
	expectResult(t, "x>1",   env, bool(x>1))
	expectResult(t, "x>=3",   env, bool(x>=3))
	expectResult(t, "-x==-5", env, bool(-x==-5))
	expectResult(t, "-x==3", env, bool(x==3))
	expectResult(t, "x!=1",  env, bool(x!=1))
}

func TestUintBinaryOps(t *testing.T) {
	x := uint32(5)
	env := makeEnv()
	env.Vars["x"] = reflect.ValueOf(&x)

	expectResult(t, "x+2",   env, x+2)
	expectResult(t, "x-2",   env, x-2)
	expectResult(t, "x*3.0",   env, x*3.0)
	expectResult(t, "x/2",   env, x/2)
	expectResult(t, "x%2",   env, x%2)
	expectResult(t, "x&1",   env, x&1)
	expectResult(t, "x|1",   env, x|1)
	expectResult(t, "x^1",   env, x^1)
	expectResult(t, "x&^1",  env, x&^1)
	expectResult(t, "x+x",  env, x+x)

	expectResult(t, "x<1",   env, bool(x<1))
	expectResult(t, "-x<=3",  env, bool(-x<=3))
	expectResult(t, "x>1",   env, bool(x>1))
	expectResult(t, "x>=3",   env, bool(x>=3))
	expectResult(t, "-x==3", env, bool(x==3))
	expectResult(t, "x!=1",  env, bool(x!=1))
}

func TestFloatBinaryOps(t *testing.T) {
	x := float32(2.25)
	env := makeEnv()
	env.Vars["x"] = reflect.ValueOf(&x)

	expectResult(t, "x+2",   env, x+2)
	expectResult(t, "x-2",   env, x-2)
	expectResult(t, "x*3.4",   env, x*3.4)
	expectResult(t, "x/2",   env, x/2)
	expectResult(t, "x+x",  env, x+x)

	expectResult(t, "x<1",   env, bool(x<1))
	expectResult(t, "-x<=3",  env, bool(-x<=3))
	expectResult(t, "x>1",   env, bool(x>1))
	expectResult(t, "x>=3",   env, bool(x>=3))
	expectResult(t, "-x==-2.25", env, bool(-x==-2.25))
	expectResult(t, "x!=1",  env, bool(x!=1))
}

func TestComplexBinaryOps(t *testing.T) {
	x := complex(1, 2)
	env := makeEnv()
	env.Vars["x"] = reflect.ValueOf(&x)

	expectResult(t, "x + complex(3, 4)", env, x + complex(3, 4))
	expectResult(t, "x - complex(3, 4)", env, x - complex(3, 4))
	expectResult(t, "x * complex(3, 4)", env, x * complex(3, 4))
	expectResult(t, "x / complex(3, 4)", env, x / complex(3, 4))
        expectResult(t, `x + x`, env, x + x)

	expectResult(t, "x == complex(1, 2)", env, bool(x == complex(1, 2)))
	expectResult(t, "x != complex(1, 2)", env, bool(x != complex(1, 2)))
}

func TestStringBinaryOps(t *testing.T) {
	x := "a"
	env := makeEnv()
	env.Vars["x"] = reflect.ValueOf(&x)

        expectResult(t, `x + "b"`, env, "a" + "b")
        expectResult(t, `x == "b"`, env, x == "b")
        expectResult(t, `x != "b"`, env, x != "b")
        expectResult(t, `x < "b"`, env, x < "b")
        expectResult(t, `x <= "b"`, env, x <= "b")
        expectResult(t, `x > "a"`, env, x > "a")
        expectResult(t, `x >= "a"`, env, x >= "a")
        expectResult(t, `x + x`, env, x + x)
}

func TestBoolBinaryOps(t *testing.T) {
	x := true
	env := makeEnv()
	env.Vars["x"] = reflect.ValueOf(&x)

        expectResult(t, "x == true", env, x == true)
        expectResult(t, "x != true", env, x != true)
        expectResult(t, "x && false", env, x && false)
        expectResult(t, "x || false", env, x || false)
        expectResult(t, `x && x`, env, x && x)
}

func TestTypedBinaryOps(t *testing.T) {

	type Foo int

	env := makeEnv()
	env.Types["Foo"] = reflect.TypeOf(Foo(0))

	expectResult(t, "Foo(1)+Foo(2)", env, Foo(1)+Foo(2))
	expectResult(t, "1-Foo(2)", env, 1-Foo(2))
	expectResult(t, "Foo(1)|2", env, Foo(1)|2)
}
