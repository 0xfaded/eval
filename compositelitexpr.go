package interactive

import (
	"errors"
	"fmt"
	"reflect"

	"go/ast"
)

func evalCompositeLit(lit *ast.CompositeLit, env *Env) (reflect.Value, bool, error) {
	t, err := evalType(lit.Type, env)
	if err != nil {
		return reflect.Value{}, true, err
	}

	switch t.Kind() {
	//case reflect.Array:
	//case reflect.Map:
	//case reflect.Slice:
	case reflect.Struct:
		vp := reflect.New(t)
		v := vp.Elem()

		if len(lit.Elts) == 0 {
			return v, true, nil
		}

		_, pairs := lit.Elts[0].(*ast.KeyValueExpr)
		for i, elt := range lit.Elts {
			var field, value reflect.Value
			var typed bool
			var fname string
			if kv, ok := elt.(*ast.KeyValueExpr); ok != pairs {
				return v, true, errors.New("Elements are either all key value pairs or not")
			} else if pairs {
				if k, ok := kv.Key.(*ast.Ident); !ok {
					return v, true, errors.New(fmt.Sprintf("Invalid key node %v %T", kv.Key, kv.Key))
				} else if f := v.FieldByName(k.Name); !f.IsValid() {
					return v, true, errors.New(t.Name() + " has no field " + k.Name)
				} else if fv, ft, err := EvalExpr(kv.Value, env); err != nil {
					return v, true, err
				} else {
					field, value = f, fv[0]
					typed = ft
					tfield, _ := v.Type().FieldByName(k.Name)
					fname = tfield.Name
				}
			} else {
				if i >= v.NumField() {
					return v, true, errors.New("Too many elements for struct " + t.Name())
				} else if _, ok := elt.(*ast.KeyValueExpr); ok {
					return v, true, errors.New("Elements are either all key value pairs or not")
				} else if fv, ft, err := EvalExpr(elt, env); err != nil {
					return v, true, err
				} else {
					field, value = v.Field(i), fv[0]
					typed = ft
					fname = v.Type().Field(i).Name
				}
			}

			if err := setTypedValue(field, value, typed); err != nil {
				return v, true, errors.New(fmt.Sprintf("cannot set %s.%s : %v",
					t.Name(), fname, value.String()))
			}
		}
		return v, true, nil

	default:
		return reflect.Value{}, true, errors.New(fmt.Sprintf("invalid type for composite literal %s", t.Name()))

	}
}
