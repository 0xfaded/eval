package eval

import (
	"reflect"
	"testing"
)

func TestIntBinaryOps(t *testing.T) {
	slice := []int {1, 2}
	env := makeEnv()
	env.Vars["slice"] = reflect.ValueOf(&slice)

	expectResult(t, "1+2",   env, int64(1)+2)
	expectResult(t, "1-2",   env, int64(1)-2)
	expectResult(t, "2*3",   env, int64(2)*3)
	expectResult(t, "5/2",   env, int64(5)/2)
	expectResult(t, "5%2",   env, int64(5)%2)
	expectResult(t, "3&1",   env, int64(3)&1)
	expectResult(t, "2|1",   env, int64(2)|1)
	expectResult(t, "3^1",   env, int64(3)^1)
	expectResult(t, "3&^1",  env, int64(3)&^1)

	expectResult(t, "3<1",   env, bool(3<1))
	expectResult(t, "-1<3",  env, bool(-1<3))
	expectResult(t, "3>1",   env, bool(3>1))
	expectResult(t, "1>3",   env, bool(1>3))
	expectResult(t, "-1==1", env, bool(-1==1))
	expectResult(t, "-1==3", env, bool(1==3))
	expectResult(t, "1!=1",  env, bool(1!=1))
	expectResult(t, "slice[0]!=3",  env, bool(slice[0]!=3))
	expectError(t, "slice[0]+int32(5)", env,
		"invalid operation <int Value> + <int32 Value> (mismatched types int and int32)")

	expectResult(t, "\"a\" + \"b\"",  env, "a" + "b")

}

// func TestUintBinaryOps(t *testing.T) {
// 	env := makeEnv()

// 	expectResult(t, "uint64(1)+2",  env, uint64(1)+2)
// 	expectResult(t, "uint64(2)-1",  env, uint64(2)-1)
// 	expectResult(t, "uint64(2)*3",  env, uint64(2)*3)
// 	expectResult(t, "uint64(5)/2",  env, uint64(5)/2)
// 	expectResult(t, "uint64(5)%2",  env, uint64(5)%2)
// 	expectResult(t, "uint64(3)&1",  env, uint64(3)&1)
// 	expectResult(t, "uint64(2)|1",  env, uint64(2)|1)
// 	expectResult(t, "uint64(3)^1",  env, uint64(3)^1)
// 	expectResult(t, "uint64(3)&^1", env, uint64(3)&^1)
// 	expectResult(t, "uint64(3)<2",  env, bool(uint64(3)<2))
// 	expectResult(t, "uint64(2)<3",  env, bool(uint64(2)<3))
// 	expectResult(t, "uint64(3)>2",  env, bool(uint64(3)>2))
// 	expectResult(t, "uint64(2)>3",  env, bool(uint64(2)>3))
// 	expectResult(t, "uint64(2)==2", env, bool(uint64(2)==2))
// 	expectResult(t, "uint64(2)==3", env, bool(uint64(2)==3))
// 	expectResult(t, "uint64(2)!=2", env, bool(uint64(2)!=2))
// 	expectResult(t, "uint64(2)!=3", env, bool(uint64(2)!=3))
// }

func TestFloatBinaryOps(t *testing.T) {
	env := makeEnv()

	expectResult(t, "1.0+2.0", env, float64(1)+2)
	expectResult(t, "1.0-2.0", env, float64(1)-2)
	expectResult(t, "2.0*3.0", env, float64(2)*3)
	expectResult(t, "5.0/2.0", env, float64(5)/2)
}

func TestComplexOps(t *testing.T) {
	env := makeEnv()

	expectResult(t, "complex(1, 2) + complex(3, 4)", env, complex(1, 2) + complex(3, 4))
	expectResult(t, "complex(1, 2) - complex(3, 4)", env, complex(1, 2) - complex(3, 4))
	expectResult(t, "complex(1, 2) * complex(3, 4)", env, complex(1, 2) * complex(3, 4))
	expectResult(t, "complex(1, 2) / complex(3, 4)", env, complex(1, 2) / complex(3, 4))

	expectResult(t, "\"a\" + \"b\"", env, "a" + "b")
	expectResult(t, "\"a\" + \"b\" == \"ab\"", env, "a" + "b" == "ab")
	expectResult(t, "\"a\" + \"b\" <= \"ab\"", env, "a" + "b" <= "ab")
	expectResult(t, "\"a\" + \"b\" >= \"ab\"", env, "a" + "b" >= "ab")
	expectResult(t, "\"a\" + \"b\" <  \"ab\"", env, "a" + "b" <  "ab")
	expectResult(t, "\"a\" + \"b\" >  \"ab\"", env, "a" + "b" >  "ab")
	expectResult(t, "\"a\" + \"b\" == \"ac\"", env, "a" + "b" == "ac")
	expectResult(t, "\"a\" + \"b\" != \"ab\"", env, "a" + "b" != "ab")
	expectResult(t, "\"a\" + \"b\" != \"ac\"", env, "a" + "b" != "ac")

}

// func TestTypedBinaryOps(t *testing.T) {

// 	type Foo int

// 	env := makeEnv()
// 	env.Types["Foo"] = reflect.TypeOf(Foo(0))

// 	expectResult(t, "Foo(1)+Foo(2)", env, Foo(1)+Foo(2))
// 	expectResult(t, "1-Foo(2)", env, 1-Foo(2))
// 	expectResult(t, "Foo(1)|2", env, Foo(1)|2)
// }
