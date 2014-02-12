package eval

import (
	"reflect"
)

func evalIdent(ident *Ident, env Env) (reflect.Value, error) {
	if ident.IsConst() {
		return ident.Const(), nil
	}

	name := ident.Name
	switch ident.source {
	case envVar:
		for searchEnv := env; searchEnv != nil; searchEnv = searchEnv.PopScope() {
			if v := searchEnv.Var(name); v.IsValid() {
				return v.Elem(), nil
			}
		}
	case envFunc:
		for searchEnv := env; searchEnv != nil; searchEnv = searchEnv.PopScope() {
			if v := searchEnv.Func(name); v.IsValid() {
				return v, nil
			}
		}
	}
        panic(dytc("missing identifier"))
}

// TODO[crc] Everything below goes with Env interface refactor
func EvalIdentExpr(ident *Ident, env Env) (*reflect.Value, bool, error) {
	v, err := evalIdent(ident, env)
	return &v, true, err
}

type EvalIdentExprFunc func(ident *Ident, env Env)  (
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
