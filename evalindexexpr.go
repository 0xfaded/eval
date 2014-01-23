package eval

import (
	"reflect"
)

func evalIndexExpr(ctx *Ctx, index *IndexExpr, env *Env) (reflect.Value, error) {
	xs, _, err := EvalExpr(ctx, index.X.(Expr), env)
	if err != nil {
		return reflect.Value{}, err
	}
	x := (*xs)[0]

	i, err := evalInteger(ctx, index.Index.(Expr), env)
	if err != nil {
		return reflect.Value{}, err
	}

	switch x.Type().Kind() {
	// case reflect.Map:
	case reflect.Ptr:
		// Short hand for array pointers
		x := x.Elem()
		fallthrough
	default:
		if !(0 <= i && i < x.Len()) {
			return reflect.Value{}, PanicIndexOutOfBounds{}
		}
		return x.Index(i), nil
	}
}
