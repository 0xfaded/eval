package eval

import (
	"errors"
	"fmt"
	"reflect"
)

func evalCompositeLit(ctx *Ctx, lit *CompositeLit, env *Env) (reflect.Value, error) {
	t := lit.KnownType()[0]

	switch t.Kind() {
	case reflect.Map:
		return evalCompositeLitMap(ctx, t, lit, env)
	case reflect.Array, reflect.Slice:
		return evalCompositeLitArrayOrSlice(ctx, t, lit, env)
	case reflect.Struct:
		return evalCompositeLitStruct(ctx, t, lit, env)
	default:
		return reflect.Value{}, errors.New(fmt.Sprintf("eval: unimplemented type for composite literal %s", t.Name()))
	}
}

func evalCompositeLitMap(ctx *Ctx, t reflect.Type, lit *CompositeLit, env *Env) (reflect.Value, error) {

	m := reflect.MakeMap(t)

	kT := knownType{t.Key()}
	vT := knownType{t.Elem()}
	for _, elt := range lit.Elts {
		kv := elt.(*KeyValueExpr)
		k, err := evalTypedExpr(ctx, kv.Key.(Expr), kT, env)
		if err != nil {
			return reflect.Value{}, err
		}
		if kT[0].Kind() == reflect.Interface {
			dynamicT := k[0].Elem().Type()
			if !isStaticTypeComparable(dynamicT) {
				return reflect.Value{}, PanicUnhashableType{dynamicT}
			}
		}
		v, err := evalTypedExpr(ctx, kv.Value.(Expr), vT, env)
		if err != nil {
			return reflect.Value{}, err
		}
		m.SetMapIndex(k[0], v[0])
	}
	return m, nil
}

func evalCompositeLitArrayOrSlice(ctx *Ctx, t reflect.Type, lit *CompositeLit, env *Env) (reflect.Value, error) {

	var v reflect.Value
	if t.Kind() == reflect.Slice {
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
		if elem, err := evalTypedExpr(ctx, elt, eT, env); err != nil {
			return reflect.Value{}, err
		} else {
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
		field := v.Field(f)
		if elem, err := evalTypedExpr(ctx, elt, knownType{field.Type()}, env); err != nil {
			return reflect.Value{}, err
		} else {
			field.Set(elem[0])
		}
	}
	return v, nil
}
