package eval


import (
	"reflect"
	"testing"
)

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

func TestBuiltinCap(t *testing.T) {
	env := makeEnv()
	slice := []int {1, 2}
	env.Vars["slice"] = reflect.ValueOf(&slice)

	expectResult(t, "cap(slice)", env, cap(slice))
	// FIXME: this is wrong. The type should be int, not int64 and
	// reflects something wrong in the eval's type system.
	expectError(t, "cap(5)", env, "invalid argument 5 (type int64) for cap")
}

func TestBuiltinLen(t *testing.T) {
	env := makeEnv()
	slice := []int {1, 2}
	env.Vars["slice"] = reflect.ValueOf(&slice)

	expectResult(t, "len(\"abc\")", env, len("abc"))
	expectResult(t, "len(slice)", env, len(slice))
	// FIXME: add tests for map, array and channel
}

func TestBuiltinNew(t *testing.T) {
	env := makeEnv()
	expr := "new(int)"
	results := getResults(t, expr, env)
	returnKind := (*results)[0].Kind().String()
	if  returnKind != "ptr" {
		t.Fatalf("Error Expecting `%s' return Kind to be `ptr' is `%s`", expr, returnKind)
	}
	expectError(t, "new(5)", env, "new parameter is not a type")
}
