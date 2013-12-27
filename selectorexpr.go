package eval

import (
	"errors"
	"fmt"
	"reflect"
)

type EvalSelectorExprFunc func(ctx *Ctx, selector *SelectorExpr, env *Env)  (
	*reflect.Value, bool, error)

func EvalSelectorExpr(ctx *Ctx, selector *SelectorExpr, env *Env) (*reflect.Value, bool, error) {
	var err error
	var x *[]reflect.Value
	if x, _, err = EvalExpr(ctx, selector.X.(Expr), env); err != nil {
		return nil, true, err
	}
	sel   := selector.Sel.Name
	x0    := (*x)[0]
	xname := x0.Type().Name()

	if x0.Kind() == reflect.Ptr {
		// Special case for handling packages
		if x0.Type() == reflect.TypeOf(Pkg(nil)) {
			sel := &Ident{ Ident: selector.Sel }
			return evalIdentExprCallback(ctx, sel, x0.Interface().(Pkg))
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
		return nil, true, errors.New(fmt.Sprintf("%s has no field or method %s", xname, sel))
	case reflect.Interface:
		if v := x0.MethodByName(sel); !v.IsValid() {
			return &v, true, errors.New(fmt.Sprintf("%s has no method %s", xname, sel))
		} else {
			return &v, true, nil
		}
	default:
		err = errors.New(fmt.Sprintf("%s.%s undefined (%s has no field or method %s)",
			xname, sel, xname, sel))
		return nil, true, err
	}
}

var evalSelectorExprCallback EvalSelectorExprFunc

func init() {
	evalSelectorExprCallback = EvalSelectorExpr
}

func SetEvalSelectorExprCallback(callback EvalSelectorExprFunc) {
	evalSelectorExprCallback = callback
}

func GetEvalSelectorExprCallback() EvalSelectorExprFunc {
	return evalSelectorExprCallback
}
