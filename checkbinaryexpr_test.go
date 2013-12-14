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
	expectConst(t, "5 / 2", env, NewConstInt64(5 / 2), ConstInt)
	expectConst(t, "5 % 2", env, NewConstInt64(5 % 2), ConstInt)
	expectConst(t, "5 & 2", env, NewConstInt64(5 & 2), ConstInt)
	expectConst(t, "5 == 2", env, 5 == 2, ConstBool)
	expectConst(t, "5 <= 2", env, 5 <= 2, ConstBool)
}

func TestBasicCheckConstBinaryIntegerRune(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "5 - 'a'", env, NewConstInt64(5 - 'a'), ConstRune)
	expectConst(t, "5 % 'a'", env, NewConstInt64(5 % 'a'), ConstRune)
	expectConst(t, "5 | 'a'", env, NewConstInt64(5 | 'a'), ConstRune)
	expectConst(t, "5 != 'a'", env, 5 != 'a', ConstBool)
	expectConst(t, "5 >= 'a'", env, 5 >= 'a', ConstBool)
}

func TestBasicCheckConstBinaryIntegerFloating(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "5 / 2.0", env, NewConstFloat64(5 / 2.0), ConstFloat)
	expectConst(t, "5 != 1.5", env, 5 != 1.5, ConstBool)
	expectConst(t, "5 < 1.5", env, 5 < 1.5, ConstBool)

	// Invalid
	expectCheckError(t, "5 % 1.0", env, "illegal constant expression: floating-point % operation")
	expectCheckError(t, "5 | 1.0", env, "illegal constant expression: ideal | ideal")
}

func TestBasicCheckConstBinaryIntegerComplex(t *testing.T) {
	env := makeEnv()

	// Vaild
	expectConst(t, "5 / 1.25i", env, NewConstComplex128(5 / 1.25i), ConstComplex)
	expectConst(t, "5 != 1.5i", env, 5 != 1.5i, ConstBool)

	// Invalid
	expectCheckError(t, "5 > 1.5i", env, "illegal constant expression: ideal > ideal")
	expectCheckError(t, "5 % 2.0i", env, "illegal constant expression: ideal % ideal")
	expectCheckError(t, "5 & 1.5i", env, "illegal constant expression: ideal & ideal")
}

// rune op X tests
func TestBasicCheckConstBinaryRuneInteger(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "'a' / 2", env, NewConstRune('a' / 2), ConstRune)
	expectConst(t, "'a' % 2", env, NewConstRune('a' % 2), ConstRune)
	expectConst(t, "'a' & 2", env, NewConstRune('a' & 2), ConstRune)
	expectConst(t, "'a' != 2", env, 'a' != 2, ConstBool)
	expectConst(t, "'a' >= 2", env, 'a' >= 2, ConstBool)
}

func TestBasicCheckConstBinaryRuneRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, "'a' / '\\x04'", env, NewConstRune('a' / '\x04'), ConstRune)
	expectConst(t, "'a' % '\\x04'", env, NewConstRune('a' % '\x04'), ConstRune)
	expectConst(t, "'a' ^ '\\x04'", env, NewConstRune('a' ^ '\x04'), ConstRune)
	expectConst(t, "'a' != '\\x04'", env, 'a' != '\x04', ConstBool)
	expectConst(t, "'a' > '\\x04'", env, 'a' > '\x04', ConstBool)
}

func TestBasicCheckConstBinaryRuneFloating(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "'d' * 1.25", env, NewConstFloat64('d' * 1.25), ConstFloat)
	expectConst(t, "'d' < 1.25", env, 'd' < 1.25, ConstBool)
	expectConst(t, "'d' != 100.0", env, 'd' != 100.0, ConstBool)

	// Invalid
	expectCheckError(t, "'d' % 1.0", env, "illegal constant expression: floating-point % operation")
	expectCheckError(t, "'d' &^ 1.0", env, "illegal constant expression: ideal &^ ideal")
}

func TestBasicCheckConstBinaryRuneComplex(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "'d' / 4i", env, NewConstComplex128('d' / 4i), ConstComplex)
	expectConst(t, "'a' == 1i", env, 'a' == 1i, ConstBool)

	// Invalid
	expectCheckError(t, "'d' > 1.5i", env, "illegal constant expression: ideal > ideal")
	expectCheckError(t, "'d' % 1.0i", env, "illegal constant expression: ideal % ideal")
	expectCheckError(t, "'d' ^ 1.0i", env, "illegal constant expression: ideal ^ ideal")
}

// floating op X tests
func TestBasicCheckConstBinaryFloatingInteger(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "1.5 + 2", env, NewConstFloat64(1.5 + 2), ConstFloat)
	expectConst(t, "1.5 < 100", env, 1.5 < 100, ConstBool)
	expectConst(t, "1.5 == 5", env, 1.5 == 5, ConstBool)

	// Invalid
	expectCheckError(t, "1.5 % 2", env, "illegal constant expression: floating-point % operation")
	expectCheckError(t, "1.5 &^ 3", env, "illegal constant expression: ideal &^ ideal")
}

func TestBasicCheckConstBinaryFloatingRune(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "1.5 - 'a'", env, NewConstFloat64(1.5 - 'a'), ConstFloat)
	expectConst(t, "1.5 < 'a'", env, 1.5 < 'a', ConstBool)
	expectConst(t, "1.5 == 'a'", env, 1.5 == 'a', ConstBool)

	// Invalid
	expectCheckError(t, "1.5 % 'a'", env, "illegal constant expression: floating-point % operation")
	expectCheckError(t, "1.5 | 'a'", env, "illegal constant expression: ideal | ideal")
}

func TestBasicCheckConstBinaryFloatingFloating(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "2.5 * 1.25", env, NewConstFloat64(2.5 * 1.25), ConstFloat)
	expectConst(t, "1.5 < 1.25", env, 1.5 < 1.25, ConstBool)
	expectConst(t, "1.5 == 1.25", env, 1.5 == 1.25, ConstBool)

	// Invalid
	expectCheckError(t, "1.5 % 1.5", env, "illegal constant expression: floating-point % operation")
	expectCheckError(t, "1.5 | 1.5", env, "illegal constant expression: ideal | ideal")
}

func TestBasicCheckConstBinaryFloatingComplex(t *testing.T) {
	env := makeEnv()

	// Valid
	expectConst(t, "2.5 / 4i", env, NewConstComplex128(2.5 / 4i), ConstComplex)
	expectConst(t, "1.5 == 1.25i", env, 1.5 == 1.25i, ConstBool)

	// Invalid
	expectCheckError(t, "1.5 < 1.25i", env, "illegal constant expression: ideal < ideal")
	expectCheckError(t, "1.5 % 1.5i", env, "illegal constant expression: ideal % ideal")
	expectCheckError(t, "1.5 | 1.5i", env, "illegal constant expression: ideal | ideal")
}

// floating op X tests
func TestBasicCheckConstBinaryComplexInteger(t *testing.T) {
	env := makeEnv()

	// Invalid
	expectConst(t, "2.5i + 2", env, NewConstComplex128(2.5i + 2), ConstComplex)
	expectConst(t, "2.5i == 2", env, 2.5i == 2, ConstBool)

	// Invalid
	expectCheckError(t, "2.5i < 2", env, "illegal constant expression: ideal < ideal")
	expectCheckError(t, "2.5i % 2", env, "illegal constant expression: ideal % ideal")
	expectCheckError(t, "2.5i | 2", env, "illegal constant expression: ideal | ideal")
}

func TestBasicCheckConstBinaryComplexRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, "2.5i - 'a'", env, NewConstComplex128(2.5i - 'a'), ConstComplex)
	expectConst(t, "2.5i == 'a'", env, 2.5i == 'a', ConstBool)

	// Invalid
	expectCheckError(t, "2.5i < 'a'", env, "illegal constant expression: ideal < ideal")
	expectCheckError(t, "2.5i % 'a'", env, "illegal constant expression: ideal % ideal")
	expectCheckError(t, "2.5i | 'a'", env, "illegal constant expression: ideal | ideal")
}

func TestBasicCheckConstBinaryComplexFloating(t *testing.T) {
	env := makeEnv()

	expectConst(t, "2.5i * 1.25", env, NewConstComplex128(2.5i * 1.25), ConstComplex)
	expectConst(t, "2.5i == 2.0", env, 2.5i == 2.0, ConstBool)

	// Invalid
	expectCheckError(t, "2.5i < 2.0", env, "illegal constant expression: ideal < ideal")
	expectCheckError(t, "2.5i % 2.0", env, "illegal constant expression: ideal % ideal")
	expectCheckError(t, "2.5i | 2.0", env, "illegal constant expression: ideal | ideal")
}

func TestBasicCheckConstBinaryComplexComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, "3i / 4i", env, NewConstComplex128(3i / 4i), ConstComplex)
	expectConst(t, "2.5i == 2i", env, 2.5i == 2i, ConstBool)

	// Invalid
	expectCheckError(t, "2.5i < 2i", env, "illegal constant expression: ideal < ideal")
	expectCheckError(t, "2.5i % 2i", env, "illegal constant expression: ideal % ideal")
	expectCheckError(t, "2.5i | 2i", env, "illegal constant expression: ideal | ideal")
}

// bool op X tests
func TestBasicCheckConstBinaryBoolBool(t *testing.T) {
	env := makeEnv()

	expectConst(t, "true == true", env, true, ConstBool)
	expectConst(t, "true != true", env, false, ConstBool)
	expectConst(t, "true == false", env, false, ConstBool)
	expectConst(t, "true != false", env, true, ConstBool)
}

