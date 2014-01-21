package eval

import (
	"errors"
	"fmt"
	"reflect"
)

func evalCompositeLit(ctx *Ctx, lit *CompositeLit, env *Env) (reflect.Value, error) {
	t := lit.KnownType()[0]

	switch t.Kind() {
	//case reflect.Map:
	case reflect.Array, reflect.Slice:
		return evalCompositeLitArrayOrSlice(ctx, t, lit, env)
	case reflect.Struct:
		return evalCompositeLitStruct(ctx, t, lit, env)
	default:
		return reflect.Value{}, errors.New(fmt.Sprintf("eval: unimplemented type for composite literal %s", t.Name()))
	}
}

func evalCompositeLitArrayOrSlice(ctx *Ctx, t reflect.Type, lit *CompositeLit, env *Env) (reflect.Value, error) {

	var v reflect.Value
	if t.Kind() == reflect.Slice {
		fmt.Printf("%v\n", lit.length)
		v = reflect.MakeSlice(t, lit.length, lit.length)
	} else {
		v = reflect.New(t).Elem()
	}

	eT := knownType{t.Elem()}
	for src, dst, i := 0, 0, 0; src < len(lit.Elts); src, dst = src + 1, dst + 1 {
		var elt Expr
		if lit.indices[i].pos == src {
			elt = lit.Elts[src].(*KeyValueExpr).Value.(Expr)
			dst = lit.indices[i].index
			i += 1
		} else {
			elt = lit.Elts[src].(Expr)
		}
		fmt.Printf("%v %v %v %v\n", src, dst, lit.indices[i].pos, lit.indices[i].index)
		if elem, err := evalTypedExpr(ctx, elt, eT, env); err != nil {
			return reflect.Value{}, err
		} else {
			fmt.Printf("elem> %v %v\n", elem[0], elem[0].Interface())
			v.Index(dst).Set(elem[0])
		}
	}
	return v, nil
}

func evalCompositeLitStruct(ctx *Ctx, t reflect.Type, lit *CompositeLit, env *Env) (reflect.Value, error) {
	v := reflect.New(t).Elem()
	for i, f := range lit.fields {
		var elt Expr
		if kv, ok := lit.Elts[i].(*KeyValueExpr); ok {
			elt = kv.Value.(Expr)
		} else {
			elt = lit.Elts[i].(Expr)
		}
		field := v.Index(f)
		if elem, err := evalTypedExpr(ctx, elt, knownType{field.Type()}, env); err != nil {
			return reflect.Value{}, err
		} else {
			field.Set(elem[0])
		}
	}
	return v, nil
}
