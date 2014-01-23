package eval

import (
	"reflect"
)

func evalSelectorExpr(ctx *Ctx, selector *SelectorExpr, env *Env) (reflect.Value, error) {

	if selector.pkgName != "" {
		vs, _, err := evalIdentExprCallback(ctx, selector.Sel, env.Pkgs[selector.pkgName])
		return *vs, err
	}

	vs, _, err := EvalExpr(ctx, selector.X.(Expr), env)
	if err != nil {
		return reflect.Value{}, err
	}
	v := (*vs)[0]
	t := v.Type()
	if selector.field != nil {
		if t.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		return v.FieldByIndex(selector.field), nil
	}

	if selector.isPtrReceiver {
		v = v.Addr()
	}
	return v.Method(selector.method), nil
}

// TODO[crc] Everything below here goes with the Env interface{} refactor
func EvalSelectorExpr(ctx *Ctx, selector *SelectorExpr, env *Env) (*reflect.Value, bool, error) {
	v, err := evalSelectorExpr(ctx, selector, env)
	return &v, true, err
}

type EvalSelectorExprFunc func(ctx *Ctx, selector *SelectorExpr, env *Env)  (
	*reflect.Value, bool, error)

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
