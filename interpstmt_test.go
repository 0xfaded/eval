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

func TestAssignToMap(t *testing.T) {
	env := MakeSimpleEnv()
	env.Vars["x"] = reflect.ValueOf(&map[int]int{})
	expectInterp(t, "x[0] = 1", env)
	expectResults(t, "x[0]", env, []interface{}{1, true})
	expectInterp(t, "x[0] = 2", env)
	expectResults(t, "x[0]", env, []interface{}{2, true})
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

func TestIfTrue(t *testing.T) {
	env := MakeSimpleEnv()
	x := 0
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "if true { x = 1 } else { x = 2 }", env)
	expectResult(t, "x", env, 1)
}

func TestIfFalse(t *testing.T) {
	env := MakeSimpleEnv()
	x := 0
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "if false { x = 1 } else { x = 2 }", env)
	expectResult(t, "x", env, 2)
}

func TestIfElseIf(t *testing.T) {
	env := MakeSimpleEnv()
	x := 0
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "if false { } else if true { x = 2 }", env)
	expectResult(t, "x", env, 2)
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

func TestFor(t *testing.T) {
	env := MakeSimpleEnv()
	x := 1
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "for i := 1; i < 4; i = i + 1 { x = x * i }", env)
	expectResult(t, "x", env, 6)
}

func TestBinaryAssign(t *testing.T) {
	env := MakeSimpleEnv()
	x := 1
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "x += 2", env)
	expectResult(t, "x", env, 3)
}

func TestSwitch(t *testing.T) {
	env := MakeSimpleEnv()
	x := -1
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "switch 1 { case 0: x = 0; case 1: x = 1; }", env)
	expectResult(t, "x", env, 1)
}

func TestSwitchDefault(t *testing.T) {
	env := MakeSimpleEnv()
	x := -1
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "switch 1 { case 0: x = 0; default: x = 1; }", env)
	expectResult(t, "x", env, 1)
}

func TestSwitchMulti(t *testing.T) {
	env := MakeSimpleEnv()
	x := -1
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "switch 1 { case 0, 1, 2: x = 1; }", env)
	expectResult(t, "x", env, 1)
}

func TestSwitchNoMatch(t *testing.T) {
	env := MakeSimpleEnv()
	x := -1
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "switch 1 { case 0, 2: x = 1; }", env)
	expectResult(t, "x", env, -1)
}

func TestSwitchEmpty(t *testing.T) {
	env := MakeSimpleEnv()
	x := -1
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "switch { case true: x = 1; }", env)
	expectResult(t, "x", env, 1)
}

func TestSwitchInit(t *testing.T) {
	env := MakeSimpleEnv()
	x := -1
	env.Vars["x"] = reflect.ValueOf(&x)
	expectInterp(t, "switch y := 1; { case true: x = y; }", env)
	expectResult(t, "x", env, 1)
}

