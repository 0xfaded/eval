package interactive

import (
	"testing"
)

// Test + Int
func TestCheckUnaryExprAddInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `+ 4`, env, NewConstInt64(+ 4), ConstInt)
}

// Test + Rune
func TestCheckUnaryExprAddRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `+ '@'`, env, NewConstRune(+ '@'), ConstRune)
}

// Test + Float
func TestCheckUnaryExprAddFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `+ 2.0`, env, NewConstFloat64(+ 2.0), ConstFloat)
}

// Test + Complex
func TestCheckUnaryExprAddComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, `+ 8.0i`, env, NewConstComplex128(+ 8.0i), ConstComplex)
}

// Test + Bool
func TestCheckUnaryExprAddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `+ true`, env,
		`invalid operation: + ideal bool`,
	)

}

// Test + String
func TestCheckUnaryExprAddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `+ "abc"`, env,
		`invalid operation: + ideal string`,
	)

}

// Test + Nil
func TestCheckUnaryExprAddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `+ nil`, env,
		`invalid operation: + nil`,
	)

}

// Test - Int
func TestCheckUnaryExprSubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `- 4`, env, NewConstInt64(- 4), ConstInt)
}

// Test - Rune
func TestCheckUnaryExprSubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `- '@'`, env, NewConstRune(- '@'), ConstRune)
}

// Test - Float
func TestCheckUnaryExprSubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `- 2.0`, env, NewConstFloat64(- 2.0), ConstFloat)
}

// Test - Complex
func TestCheckUnaryExprSubComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, `- 8.0i`, env, NewConstComplex128(- 8.0i), ConstComplex)
}

// Test - Bool
func TestCheckUnaryExprSubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `- true`, env,
		`invalid operation: - ideal bool`,
	)

}

// Test - String
func TestCheckUnaryExprSubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `- "abc"`, env,
		`invalid operation: - ideal string`,
	)

}

// Test - Nil
func TestCheckUnaryExprSubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `- nil`, env,
		`invalid operation: - nil`,
	)

}

// Test ^ Int
func TestCheckUnaryExprXorInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `^ 4`, env, NewConstInt64(^ 4), ConstInt)
}

// Test ^ Rune
func TestCheckUnaryExprXorRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `^ '@'`, env, NewConstRune(^ '@'), ConstRune)
}

// Test ^ Float
func TestCheckUnaryExprXorFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `^ 2.0`, env,
		`illegal constant expression ^ ideal`,
	)

}

// Test ^ Complex
func TestCheckUnaryExprXorComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `^ 8.0i`, env,
		`illegal constant expression ^ ideal`,
	)

}

// Test ^ Bool
func TestCheckUnaryExprXorBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `^ true`, env,
		`invalid operation: ^ ideal bool`,
	)

}

// Test ^ String
func TestCheckUnaryExprXorString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `^ "abc"`, env,
		`invalid operation: ^ ideal string`,
	)

}

// Test ^ Nil
func TestCheckUnaryExprXorNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `^ nil`, env,
		`invalid operation: ^ nil`,
	)

}
