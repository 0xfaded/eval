package eval

import (
	"reflect"
)

func evalTypeAssertExpr(ctx *Ctx, assert *TypeAssertExpr, env *Env) (reflect.Value, error) {
	x := assert.X.(Expr)
	if vs, _, err := EvalExpr(ctx, x, env); err != nil {
		return reflect.Value{}, err
	} else {
		v := (*vs)[0]
		xT := x.KnownType()[0]
		aT := assert.KnownType()[0]
		if v.IsNil() {
			return reflect.Value{}, PanicInterfaceConversion{aT: aT}
		}
		dynamic := v.Elem()
		dT := dynamic.Type()
		if aT.Kind() == reflect.Interface {
			if !dT.Implements(aT) {
				return reflect.Value{}, PanicInterfaceConversion{xT, aT, nil}
			}
		} else {
			if dT != aT {
				return reflect.Value{}, PanicInterfaceConversion{xT, aT, dT}
			}
		}
		r := reflect.New(aT).Elem()
		r.Set(dynamic)
		return r, nil
	}
}
