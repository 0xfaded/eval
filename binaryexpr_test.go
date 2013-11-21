package interactive

import (
	"testing"
)

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

