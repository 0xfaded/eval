package interactive

import (
	"testing"
	"reflect"

	"go/parser"
)

/*
	expr, err := parser.ParseExpr("uint8(2500) + int8(10)")
	expr, err := parser.ParseExpr("log.Printf(\"Help %v %v %d\", 1, 2, 3)")
	expr, err := parser.ParseExpr("log.Printf(\"Help\")")
	expr, err := parser.ParseExpr("log.New(stdout, \"Awesome> \", log.Ltime).Printf(\":)\")")
*/

func TestBuiltinComplex(t *testing.T) {
	env := makeEnv()

	expectResult(t, "complex(1, 2)", env, complex(1, 2))
	expectResult(t, "complex(float64(1), 2)", env, complex(float64(1), 2))
	expectResult(t, "complex(1, float64(2))", env, complex(1, float64(2)))
	expectResult(t, "complex(float64(1), float64(2))", env, complex(float64(1), float64(2)))
	expectResult(t, "complex(float32(1), 2)", env, complex(float32(1), 2))
	expectResult(t, "complex(1, float32(2))", env, complex(1, float32(2)))
	expectResult(t, "complex(float32(1), float32(2))", env, complex(float32(1), float32(2)))
}

func TestBuiltinReal(t *testing.T) {
	env := makeEnv()

	expectResult(t, "real(complex(1, 2))", env, real(complex(1, 2)))
	expectResult(t, "real(complex(float64(1), 2))", env, real(complex(float64(1), 2)))
	expectResult(t, "real(complex(1, float64(2)))", env, real(complex(1, float64(2))))
	expectResult(t, "real(complex(float64(1), float64(2)))", env, real(complex(float64(1), float64(2))))
	expectResult(t, "real(complex(float32(1), 2))", env, real(complex(float32(1), 2)))
	expectResult(t, "real(complex(1, float32(2)))", env, real(complex(1, float32(2))))
	expectResult(t, "real(complex(float32(1), float32(2)))", env, real(complex(float32(1), float32(2))))
}

func TestBuiltinImag(t *testing.T) {
	env := makeEnv()

	expectResult(t, "imag(complex(1, 2))", env, imag(complex(1, 2)))
	expectResult(t, "imag(complex(float64(1), 2))", env, imag(complex(float64(1), 2)))
	expectResult(t, "imag(complex(1, float64(2)))", env, imag(complex(1, float64(2))))
	expectResult(t, "imag(complex(float64(1), float64(2)))", env, imag(complex(float64(1), float64(2))))
	expectResult(t, "imag(complex(float32(1), 2))", env, imag(complex(float32(1), 2)))
	expectResult(t, "imag(complex(1, float32(2)))", env, imag(complex(1, float32(2))))
	expectResult(t, "imag(complex(float32(1), float32(2)))", env, imag(complex(float32(1), float32(2))))
}

func TestIntBinaryOps(t *testing.T) {
	env := makeEnv()

	expectResult(t, "1+2", env, int64(1)+2)
	expectResult(t, "1-2", env, int64(1)-2)
	expectResult(t, "2*3", env, int64(2)*3)
	expectResult(t, "5/2", env, int64(5)/2)
	expectResult(t, "5%2", env, int64(5)%2)
	expectResult(t, "3&1", env, int64(3)&1)
	expectResult(t, "2|1", env, int64(2)|1)
	expectResult(t, "3^1", env, int64(3)^1)
	expectResult(t, "3&^1", env, int64(3)&^1)
}

func TestUintBinaryOps(t *testing.T) {
	env := makeEnv()

	expectResult(t, "uint64(1)+2", env, uint64(1)+2)
	expectResult(t, "uint64(2)-1", env, uint64(2)-1)
	expectResult(t, "uint64(2)*3", env, uint64(2)*3)
	expectResult(t, "uint64(5)/2", env, uint64(5)/2)
	expectResult(t, "uint64(5)%2", env, uint64(5)%2)
	expectResult(t, "uint64(3)&1", env, uint64(3)&1)
	expectResult(t, "uint64(2)|1", env, uint64(2)|1)
	expectResult(t, "uint64(3)^1", env, uint64(3)^1)
	expectResult(t, "uint64(3)&^1", env, uint64(3)&^1)
}

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
}

// << and >> ops behave diffently 
func TestIntBinaryShiftOps(t *testing.T) {
	env := makeEnv()

	expectResult(t, "1+2", env, int64(1)+2)
	expectResult(t, "1-2", env, int64(1)-2)
	expectResult(t, "2*3", env, int64(2)*3)
	expectResult(t, "5/2", env, int64(5)/2)
	expectResult(t, "5%2", env, int64(5)%2)
	expectResult(t, "3&1", env, int64(3)&1)
	expectResult(t, "2|1", env, int64(2)|1)
	expectResult(t, "3^1", env, int64(3)^1)
	expectResult(t, "3&^1", env, int64(3)&^1)
}

func TestCompositeStructValues(t *testing.T) {
	type Alice struct {
		Bob int
	}

	env := makeEnv()
	env.types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 10 }
	expr := "Alice{ 10 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeStructKeyValues(t *testing.T) {
	type Alice struct {
		Bob int
	}

	env := makeEnv()
	env.types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { Bob: 10 }
	expr := "Alice{ Bob: 10 }"

	expectResult(t, expr, env, expected)
}

func expectVoid(t *testing.T, expr string, env *env) {
	expectResults(t, expr, env, []interface{}{})
}

func expectResult(t *testing.T, expr string, env *env, expected interface{}) {
	expectResults(t, expr, env, []interface{}{expected})
}

func expectResults(t *testing.T, expr string, env *env, expected []interface{}) {
	if e, err := parser.ParseExpr(expr); err != nil {
		t.Fatalf("Failed to parse expression '%s' (%v)", expr, err)
	} else if results, _, err := evalExpr(e, env); err != nil {
		t.Fatalf("Error evaluating expression '%s' (%v)", expr, err)
	} else {
		resultsi := make([]interface{}, len(results))
		for i, result := range results {
			resultsi[i] = result.Interface()
		}
		if !reflect.DeepEqual(resultsi, expected) {
			t.Fatalf("Expression '%s' yielded '%+v', expected '%+v'", expr, resultsi, expected)
		}
	}
}

func makeEnv() *env {
	return &env {
		vars: make(map[string] reflect.Value),
		consts: make(map[string] reflect.Value),
		funcs: make(map[string] reflect.Value),
		types: make(map[string] reflect.Type),
		pkgs: make(map[string] pkg),
	}
}
