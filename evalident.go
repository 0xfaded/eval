package eval

import (
	"reflect"
)

func evalIdent(ctx *Ctx, ident *Ident, env *Env) (reflect.Value, error) {
	if ident.IsConst() {
		return ident.Const(), nil
	}

	name := ident.Name
	switch ident.source {
	case envVar:
		return env.Vars[name].Elem(), nil
	case envFunc:
		return env.Funcs[name], nil
	default:
                panic(dytc("missing identifier"))
	}
}

// TODO[crc] Everything below goes with Env interface refactor
func EvalIdentExpr(ctx *Ctx, ident *Ident, env *Env) (*reflect.Value, bool, error) {
	v, err := evalIdent(ctx, ident, env)
	return &v, true, err
}

type EvalIdentExprFunc func(ctx *Ctx, ident *Ident, env *Env)  (
	*reflect.Value, bool, error)

func DerefValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		return v.Elem()
	default:
		return v
	}
}

var evalIdentExprCallback EvalIdentExprFunc = EvalIdentExpr

func SetEvalIdentExprCallback(callback EvalIdentExprFunc) {
	evalIdentExprCallback = callback
}

func GetEvalIdentExprCallback() EvalIdentExprFunc {
	return evalIdentExprCallback
}
