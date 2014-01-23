package eval

import (
	"reflect"
)

// TODO[crc] support slice[::] syntax after go1.2 upgrade
func evalSliceExpr(ctx *Ctx, slice *SliceExpr, env *Env) (reflect.Value, error) {
	xs, _, err := EvalExpr(ctx, slice.X.(Expr), env)
	if err != nil {
		return reflect.Value{}, err
	}
	x := (*xs)[0]

	var l, h int
	if slice.Low != nil {
		if l, err = evalInteger(ctx, slice.Low.(Expr), env); err != nil {
			return reflect.Value{}, err
		}
	}
	if slice.High != nil {
		if h, err = evalInteger(ctx, slice.High.(Expr), env); err != nil {
			return reflect.Value{}, err
		}
	} else {
		h = x.Len()
	}

	t := slice.KnownType()[0]
	switch t.Kind() {
	case reflect.Ptr:
		// Short hand for array pointers
		x = x.Elem()
		fallthrough
	case reflect.Array, reflect.String:
		if l < 0 || h > x.Len() || h < l {
			return reflect.Value{}, PanicSliceOutOfBounds{}
		}
	case reflect.Slice:
		if l < 0 || h > x.Cap() || h < l {
			return reflect.Value{}, PanicSliceOutOfBounds{}
		}
	}
	return x.Slice(l, h), nil
}
