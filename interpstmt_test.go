package eval

import (
	"reflect"
	"testing"
)

func TestAssignTypedInt(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, "x := int(1)", env)
	if x, ok := env.Vars["x"]; !ok {
		t.Fatalf("x not in env")
	} else if x.Elem().Type() != reflect.TypeOf(int(1)) {
		t.Fatalf("x has type %v, expected int")
	}
}

func TestAssignUntypeInt(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, "x := 1", env)
	if x, ok := env.Vars["x"]; !ok {
		t.Fatalf("x not in env")
	} else if x.Elem().Type() != reflect.TypeOf(int(1)) {
		t.Fatalf("x has type %v, expected int", x.Elem().Type())
	}
}

func TestAssignMultiNew(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, "x, y := 1, float32(1.5)", env)
	if x, ok := env.Vars["x"]; !ok {
		t.Fatalf("x not in env")
	} else if x.Elem().Type() != reflect.TypeOf(int(1)) {
		t.Fatalf("x has type %v, expected int", x.Elem().Type())
	}
	if y, ok := env.Vars["y"]; !ok {
		t.Fatalf("y not in env")
	} else if y.Elem().Type() != reflect.TypeOf(float32(1)) {
		t.Fatalf("y has type %v, expected int", y.Elem().Type())
	}
}

func TestAssignMulti(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, "x := 1", env)
	expectInterp(t, "x, y := 2, 3", env)
	if x, ok := env.Vars["x"]; !ok {
		t.Fatalf("x not in env")
	} else if x.Elem().Type() != reflect.TypeOf(int(1)) {
		t.Fatalf("x has type %v, expected int", x.Elem().Type())
	} else if x.Elem().Int() != 2 {
		t.Fatalf("x has value %v, expected 2", x.Int())
	}
	if y, ok := env.Vars["y"]; !ok {
		t.Fatalf("y not in env")
	} else if y.Elem().Type() != reflect.TypeOf(int(1)) {
		t.Fatalf("y has type %v, expected int", y.Elem().Type())
	} else if y.Elem().Int() != 3 {
		t.Fatalf("y has value %v, expected 3", y.Int())
	}
}

func TestAssignMapIndexAbsent(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, "x, ok := map[int]int{}[1]", env)
	if x, ok := env.Vars["x"]; !ok {
		t.Fatalf("x not in env")
	} else if x.Elem().Type() != reflect.TypeOf(int(1)) {
		t.Fatalf("x has type %v, expected int", x.Elem().Type())
	} else if x.Elem().Int() != 0 {
		t.Fatalf("x has value %v, expected 0", x.Int())
	}
	if y, ok := env.Vars["ok"]; !ok {
		t.Fatalf("ok not in env")
	} else if y.Elem().Type() != reflect.TypeOf(true) {
		t.Fatalf("ok has type %v, expected int", y.Elem().Type())
	} else if y.Elem().Bool() != false {
		t.Fatalf("ok has value %v, expected false", y.Bool())
	}
}

func TestAssignMapIndexPresent(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, "x, ok := map[int]int{1:1}[1]", env)
	if x, ok := env.Vars["x"]; !ok {
		t.Fatalf("x not in env")
	} else if x.Elem().Type() != reflect.TypeOf(int(1)) {
		t.Fatalf("x has type %v, expected int", x.Elem().Type())
	} else if x.Elem().Int() != 1 {
		t.Fatalf("x has value %v, expected 1", x.Int())
	}
	if y, ok := env.Vars["ok"]; !ok {
		t.Fatalf("ok not in env")
	} else if y.Elem().Type() != reflect.TypeOf(true) {
		t.Fatalf("ok has type %v, expected int", y.Elem().Type())
	} else if y.Elem().Bool() != true {
		t.Fatalf("ok has value %v, expected true", y.Bool())
	}
}

func TestAssignBadTypeAssert(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, "x, ok := interface{}(float32(1)).(int)", env)
	if x, ok := env.Vars["x"]; !ok {
		t.Fatalf("x not in env")
	} else if x.Elem().Type() != reflect.TypeOf(int(1)) {
		t.Fatalf("x has type %v, expected int", x.Elem().Type())
	} else if x.Elem().Int() != 0 {
		t.Fatalf("x has value %v, expected 0", x.Int())
	}
	if y, ok := env.Vars["ok"]; !ok {
		t.Fatalf("ok not in env")
	} else if y.Elem().Type() != reflect.TypeOf(true) {
		t.Fatalf("ok has type %v, expected int", y.Elem().Type())
	} else if y.Elem().Bool() != false {
		t.Fatalf("ok has value %v, expected false", y.Bool())
	}
}

func TestAssignUnderscore(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, "x, _ := map[int]int{1:1}[1]", env)
	expectInterp(t, "_, _ = map[int]int{1:1}[1]", env)
}

func TestCheckAssignStmtParenUnderscore(t *testing.T) {
	env := MakeSimpleEnv()
	f := func() (int, int) { return 1, 1 }
	env.Vars["f"] = reflect.ValueOf(&f)
	expectInterp(t, `(_), f = 1, nil`, env)
}

// Test DefNil
func TestCheckAssignStmtDefNil(t *testing.T) {
	env := MakeSimpleEnv()
	expectInterp(t, `nil := 1`, env)
}

