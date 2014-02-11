package eval

import (
	"reflect"
)

func evalIndexExpr(ctx *Ctx, index *IndexExpr, env *Env) ([]reflect.Value, error) {
	xs, _, err := EvalExpr(ctx, index.X.(Expr), env)
	if err != nil {
		return []reflect.Value{}, err
	}
	x := (*xs)[0]

	t := index.X.(Expr).KnownType()[0]
	switch t.Kind() {
	case reflect.Map:
		k, err := evalTypedExpr(ctx, index.Index.(Expr), knownType{t.Key()}, env)
		if err != nil {
			return []reflect.Value{}, err
		}
		v := x.MapIndex(k[0])
		ok := v.IsValid()
		if !ok {
			v = reflect.New(t.Key()).Elem()
		}
		// TODO[crc] return ok as well when assignment support is added
		return []reflect.Value{v}, nil
	case reflect.Ptr:
		// Short hand for array pointers
		x = x.Elem()
		fallthrough
	default:
		i, err := evalInteger(ctx, index.Index.(Expr), env)
		if err != nil {
			return []reflect.Value{}, err
		}
		if !(0 <= i && i < x.Len()) {
			return []reflect.Value{}, PanicIndexOutOfBounds{}
		}
		return []reflect.Value{x.Index(i)}, nil
	}
}
