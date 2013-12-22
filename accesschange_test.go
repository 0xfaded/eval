// Tests replacing the default identifier selection lookup value mechanism with
// our own custom versions.

package interactive

import (
	"testing"
	"reflect"
)

// Here's our custom ident lookup.
func MyEvalIdentExpr(ctx *Ctx, ident *Ident, env *Env) (
	*reflect.Value, bool, error) {
	name := ident.Name
	if name == "nil" {
		return nil, false, nil
	} else if name[0] == 'v' {
		val := reflect.ValueOf(5)
		return &val, true, nil
	} else if name[0] == 'c' {
		val := reflect.ValueOf("constant")
		return &val, true, nil
	} else if name[0] == 'c' {
		val := reflect.ValueOf(true)
		return &val, true, nil
	} else {
		val := reflect.ValueOf('x')
		return &val, true, nil
	}
}


// Here's our custom selector lookup.
func MyEvalSelectorExpr(ctx *Ctx, selector *SelectorExpr, env *Env) (
	*reflect.Value, bool, error) {
	var err error
	var x *[]reflect.Value
	if x, _, err = EvalExpr(ctx, selector.X.(Expr), env); err != nil {
		return nil, true, err
	}
	sel   := selector.Sel.Name
	x0    := (*x)[0]

	if x0.Kind() == reflect.Ptr {
		// Special case for handling packages
		if x0.Type() == reflect.TypeOf(Pkg(nil)) {
			val := reflect.ValueOf("bogus")
			return &val, true, nil
		} else if !x0.IsNil() && x0.Elem().Kind() == reflect.Struct {
			x0 = x0.Elem()
		}
	}

	switch x0.Type().Kind() {
	case reflect.Struct:
		if v := x0.FieldByName(sel); v.IsValid() {
			return &v, true, nil
		} else if x0.CanAddr() {
			if v := x0.Addr().MethodByName(sel); v.IsValid() {
				return &v, true, nil
			}
		}
		return nil, true, nil
	case reflect.Interface:
		if v := x0.MethodByName(sel); !v.IsValid() {
			return &v, true, nil
		} else {
			return &v, true, nil
		}
	default:
		return nil, true, nil
	}
}


func TestReplaceIdentLookup(t *testing.T) {
	defer SetEvalIdentExprCallback(EvalIdentExpr)
	env := makeEnv()
	SetEvalIdentExprCallback(MyEvalIdentExpr)
	expectResult(t, "fdafdsa", env, 'x')
	expectResult(t, "c + \" value\"", env, "constant value")

}


func TestReplaceSelectorLookup(t *testing.T) {
	defer SetEvalSelectorExprCallback(EvalSelectorExpr)
	env  := makeEnv()
	pkg := makeEnv()
	env.Pkgs["bogusPackage"] = pkg
	SetEvalSelectorExprCallback(MyEvalSelectorExpr)
	expectResult(t, "bogusPackage.something", env, "bogus")

}
