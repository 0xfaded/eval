package eval

import (
	"testing"
	"reflect"
)

func TestSliceArray(t *testing.T) {
	a := [3]int{1, 2, 3}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expected := a[1:1]
	expr := "a[1:1]"

	expectResult(t, expr, env, expected)

	expected = a[0:2]
	expr = "a[0:2]"

	expectResult(t, expr, env, expected)
}

func TestSliceSlice(t *testing.T) {
	a := []int{1, 2, 3}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expected := a[0:2]
	expr := "a[0:2]"

	expectResult(t, expr, env, expected)

	expected = a[1:1]
	expr = "a[1:1]"

	expectResult(t, expr, env, expected)

}

func TestSliceString(t *testing.T) {
	a := "abc"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expected := a[1:1]
	expr := "a[1:1]"

	expectResult(t, expr, env, expected)

	expected = a[0:2]
	expr = "a[0:2]"

	expectResult(t, expr, env, expected)

}

func TestSliceArrayConstantOutOfBounds(t *testing.T) {
	a := [2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2:2]"
	expectError(t, expr, env, "invalid array index 2 (out of bounds for 2-element array)")
}

func TestSliceArrayNegativeIndex(t *testing.T) {
	a := [2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[-1:1]"
	expectError(t, expr, env, "invalid array index -1 (index must be non-negative)")
}

func TestSliceArrayNonIntIndex(t *testing.T) {
	a := [2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := `a["abc":1]`
	expectError(t, expr, env, `non-integer array index "abc"`)
}

func TestSliceArrayPtrConstantOutOfBounds(t *testing.T) {
	a := &[2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2:2]"
	expectError(t, expr, env, "invalid array index 2 (out of bounds for 2-element array)")
}

func TestSliceArrayPtrNegativeIndex(t *testing.T) {
	a := &[2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[-1:1]"
	expectError(t, expr, env, "invalid array index -1 (index must be non-negative)")
}

func TestSliceArrayPtrNonIntIndex(t *testing.T) {
	a := &[2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := `a["abc":1]`
	expectError(t, expr, env, `non-integer array index "abc"`)
}

func TestSliceSliceNegativeIndex(t *testing.T) {
	a := []int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[-1:1]"
	expectError(t, expr, env, "invalid slice index -1 (index must be non-negative)")
}

func TestSliceSliceNonIntIndex(t *testing.T) {
	a := []int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := `a["abc":1]`
	expectError(t, expr, env, `non-integer slice index "abc"`)
}

/* TODO string constants should error when out of bounds
func TestSliceStringConstantOutOfBounds(t *testing.T) {
	a := "ab"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2]"
	expectError(t, expr, env, "invalid array index 2 (out of bounds for 2-byte string)")
}
*/

func TestSliceStringNegativeIndex(t *testing.T) {
	a := "ab"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[-1:1]"
	expectError(t, expr, env, "invalid string index -1 (index must be non-negative)")
}

func TestSliceStringLowHigh(t *testing.T) {
	a := "ab"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2:1]"
	expectError(t, expr, env, "string index out of range")

	a = "abc"
	expectError(t, expr, env, "invalid slice index: 2 > 1")
}

func TestSliceStringNonIntIndex(t *testing.T) {
	a := "ab"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := `a["abc":2]`
	expectError(t, expr, env, `string index out of range`)
}

func TestSliceInvalidIndexInt(t *testing.T) {
	a := 1

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[1:1]"
	expectError(t, expr, env, `cannot slice (type int)`)
}

func TestSliceInvalidIndexSlicePtr(t *testing.T) {
	a := &[]int{1,2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[1:1]"
	expectError(t, expr, env, `cannot slice (type *[]int)`)
}

func TestSliceInvalidIndexArray(t *testing.T) {
	a := []int{1,2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2:2]"
	expectError(t, expr, env, `slice index out of range`)
}
