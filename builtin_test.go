package interactive


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

func TestBuiltinLen(t *testing.T) {
	env := makeEnv()
	slice := []int {1, 2}
	env.Vars["slice"] = reflect.ValueOf(&slice)

	expectResult(t, "len(\"abc\")", env, len("abc"))
	expectResult(t, "len(slice)", env, len(slice))
	// FIXME: add tests for map, array and channel
}
