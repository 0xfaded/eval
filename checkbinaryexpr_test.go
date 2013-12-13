package interactive

import (
	"testing"
)

// This file contains three groups of tests.
// Key: ck. check, ov. overflow, tr. truncation, bo. bad operation
// 1. Untyped op Untyped. These tests are divided accordingly
//
//	    | integer |  rune   |floating | complex |  bool   | string  |   nil   |
//          +---------+---------+---------|---------+---------+---------+---------+
// integer  |         |         |         |         |         |         |         |
// rune     |         |         |         |         |         |         |         |
// floating |         |         |         |         |         |         |         |
// complex  |         |         |         |         |         |         |         |
// bool     |         |         |         |         |         |         |         |
// string   |         |         |         |         |         |         |         |
// nil      |         |         |         |         |         |         |         |

// integer op X tests
func TestBasicCheckConstBinaryIntegerInteger(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "5 + 2", env, NewBigInt64(5 + 2), ConstInt)
	expectConst(t, "5 % 2", env, NewBigInt64(5 % 2), ConstInt)
	expectConst(t, "5 & 2", env, NewBigInt64(5 & 2), ConstInt)
	expectConst(t, "5 == 2", env, 5 == 2, ConstBool)
	expectConst(t, "5 <= 2", env, 5 <= 2, ConstBool)
}

func TestBasicCheckConstBinaryIntegerRune(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "5 - 'a'", env, NewBigInt64(5 - 'a'), ConstRune)
	expectConst(t, "5 % 'a'", env, NewBigInt64(5 % 'a'), ConstRune)
	expectConst(t, "5 | 'a'", env, NewBigInt64(5 | 'a'), ConstRune)
	expectConst(t, "5 != 'a'", env, 5 != 'a', ConstBool)
	expectConst(t, "5 >= 'a'", env, 5 >= 'a', ConstBool)
}

func TestBasicCheckConstBinaryIntegerFloating(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "5 / 1.25", env, NewBigFloat64(5 / 1.25), ConstFloat)
	expectConst(t, "5 != 1.5", env, 5 != 1.5, ConstBool)
	expectConst(t, "5 < 1.5", env, 5 < 1.5, ConstBool)

	// Invalid
	expectCheckError(t, "5 % 1.0", env, "illegal constant expression: floating-point % operation")
	expectCheckError(t, "5 | 1.0", env, "illegal constant expression: ideal | ideal")
}

func TestBasicCheckConstBinaryIntegerComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, "5 * 1.25i", env, NewBigComplex128(5 * 1.25i), ConstComplex)
	expectConst(t, "5 != 1.5i", env, 5 != 1.5i, ConstBool)
}

// rune op X tests
func TestBasicCheckConstBinaryRuneInteger(t *testing.T) {
	env := makeEnv()

	expectConst(t, "'a' + 2", env, NewBigRune('a' + 2), ConstRune)
	expectConst(t, "'a' == 5", env, 'a' == 5, ConstBool)
}

func TestBasicCheckConstBinaryRuneRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, "'a' - 'a'", env, NewBigRune('a' - 'a'), ConstRune)
	expectConst(t, "'a' != 'a'", env, 'a' != 'a', ConstBool)
}

func TestBasicCheckConstBinaryRuneFloating(t *testing.T) {
	env := makeEnv()

	expectConst(t, "'d' * 1.25", env, NewBigFloat64('d' * 1.25), ConstFloat)
	expectConst(t, "'d' < 1.25", env, 'd' < 1.25, ConstBool)
}

func TestBasicCheckConstBinaryRuneComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, "'d' / 4i", env, NewBigComplex128('d' / 4i), ConstComplex)
	expectConst(t, "'a' == 1i", env, 'a' == 1i, ConstBool)
}

// floating op X tests
func TestBasicCheckConstBinaryFloatingInteger(t *testing.T) {
	env := makeEnv()

	expectConst(t, "1.5 + 2", env, NewBigFloat64(1.5 + 2), ConstFloat)
	expectConst(t, "1.5 == 5", env, 1.5 == 5, ConstBool)
}

func TestBasicCheckConstBinaryFloatingRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, "1.5 - 'a'", env, NewBigFloat64(1.5 - 'a'), ConstFloat)
	expectConst(t, "1.5 == 'a'", env, 1.5 == 'a', ConstBool)
}

func TestBasicCheckConstBinaryFloatingFloating(t *testing.T) {
	env := makeEnv()

	expectConst(t, "2.5 * 1.25", env, NewBigFloat64(2.5 * 1.25), ConstFloat)
	expectConst(t, "1.5 == 'a'", env, 1.5 == 'a', ConstBool)
}

func TestBasicCheckConstBinaryFloatingComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, "2.5 / 4i", env, NewBigComplex128(2.5 / 4i), ConstComplex)
}

// floating op X tests
func TestBasicCheckConstBinaryComplexInteger(t *testing.T) {
	env := makeEnv()

	expectConst(t, "2.5i + 2", env, NewBigComplex128(2.5i + 2), ConstComplex)
}

func TestBasicCheckConstBinaryComplexRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, "2.5i - 'a'", env, NewBigComplex128(2.5i - 'a'), ConstComplex)
}

func TestBasicCheckConstBinaryComplexComplexing(t *testing.T) {
	env := makeEnv()

	expectConst(t, "2.5i * 1.25", env, NewBigComplex128(2.5i * 1.25), ConstComplex)
}

func TestBasicCheckConstBinaryComplexComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, "3i / 4i", env, NewBigComplex128(3i / 4i), ConstComplex)
}

// bool op X tests
func TestBasicCheckConstBinaryBoolBool(t *testing.T) {
	env := makeEnv()

	expectConst(t, "true == true", env, true, ConstBool)
	expectConst(t, "true != true", env, false, ConstBool)
	expectConst(t, "true == false", env, false, ConstBool)
	expectConst(t, "true != false", env, true, ConstBool)
}

