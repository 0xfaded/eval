package eval

import (
	"testing"
	"reflect"
)

func TestIndexArray(t *testing.T) {
	a := [2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expected := a[1]
	expr := "a[1]"

	expectResult(t, expr, env, expected)
}

func TestIndexArrayPtr(t *testing.T) {
	a := &[2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expected := a[1]
	expr := "a[1]"

	expectResult(t, expr, env, expected)
}

func TestIndexSlice(t *testing.T) {
	a := []int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expected := a[1]
	expr := "a[1]"

	expectResult(t, expr, env, expected)
}

func TestIndexString(t *testing.T) {
	a := "ab"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expected := a[1]
	expr := "a[1]"

	expectResult(t, expr, env, expected)
}

func TestIndexArrayConstantOutOfBounds(t *testing.T) {
	a := [2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2]"
	expectError(t, expr, env, "invalid array index 2 (out of bounds for 2-element array)")
}

func TestIndexArrayNegativeIndex(t *testing.T) {
	a := [2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[-1]"
	expectError(t, expr, env, "invalid array index -1 (index must be non-negative)")
}

func TestIndexArrayNonIntIndex(t *testing.T) {
	a := [2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := `a["abc"]`
	expectError(t, expr, env, `non-integer array index "abc"`)
}

func TestIndexArrayPtrConstantOutOfBounds(t *testing.T) {
	a := &[2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2]"
	expectError(t, expr, env, "invalid array index 2 (out of bounds for 2-element array)")
}

func TestIndexArrayPtrNegativeIndex(t *testing.T) {
	a := &[2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[-1]"
	expectError(t, expr, env, "invalid array index -1 (index must be non-negative)")
}

func TestIndexArrayPtrNonIntIndex(t *testing.T) {
	a := &[2]int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := `a["abc"]`
	expectError(t, expr, env, `non-integer array index "abc"`)
}

func TestIndexSliceNegativeIndex(t *testing.T) {
	a := []int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[-1]"
	expectError(t, expr, env, "invalid slice index -1 (index must be non-negative)")
}

func TestIndexSliceNonIntIndex(t *testing.T) {
	a := []int{1, 2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := `a["abc"]`
	expectError(t, expr, env, `non-integer slice index "abc"`)
}

/* TODO string constants should error when out of bounds
func TestIndexStringConstantOutOfBounds(t *testing.T) {
	a := "ab"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2]"
	expectError(t, expr, env, "invalid array index 2 (out of bounds for 2-byte string)")
}
*/

func TestIndexStringNegativeIndex(t *testing.T) {
	a := "ab"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[-1]"
	expectError(t, expr, env, "invalid string index -1 (index must be non-negative)")
}

func TestIndexStringNonIntIndex(t *testing.T) {
	a := "ab"

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := `a["abc"]`
	expectError(t, expr, env, `non-integer string index "abc"`)
}

func TestInvalidIndexInt(t *testing.T) {
	a := 1

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[1]"
	expectError(t, expr, env, `invalid operation: a[1] (index of type int)`)
}

func TestInvalidIndexSlicePtr(t *testing.T) {
	a := &[]int{1,2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[1]"
	expectError(t, expr, env, `invalid operation: a[1] (index of type *[]int)`)
}

func TestInvalidIndexArray(t *testing.T) {
	a := []int{1,2}

	env := makeEnv()
	env.Vars["a"] = reflect.ValueOf(&a)

	expr := "a[2]"
	expectError(t, expr, env, `reflect: slice index out of range`)
}
