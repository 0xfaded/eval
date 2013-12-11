package interactive

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"go/ast"
	"go/token"
)

func evalCompositeLit(ctx *Ctx, lit *CompositeLit, env *Env) (*reflect.Value, bool, error) {
	t, err := evalType(ctx, lit.Type, env)
	if err != nil {
		return nil, true, err
	}

	switch t.Kind() {
	//case reflect.Map:
	case reflect.Array, reflect.Slice:
		return evalCompositeLitArrayOrSlice(ctx, t, lit, env)
	case reflect.Struct:
		return evalCompositeLitStruct(ctx, t, lit, env)
	default:
		return nil, true, errors.New(fmt.Sprintf("invalid type for composite literal %s", t.Name()))
	}
}

func evalCompositeLitArrayOrSlice(ctx *Ctx, t reflect.Type, lit *CompositeLit, env *Env) (*reflect.Value, bool, error) {

	v := reflect.New(t).Elem()

	var curKey uint64 = 0
	var size uint64 = 0
	if t.Kind() == reflect.Array {
		size = uint64(t.Len())
	} else {
		// Check all keys are valid and calculate slice size.
		// Elements with key are placed at the keyed position.
		// Elements without are placed one after the previous.
		// For example, []int{1, 2:1, 1} -> [1, 0, 1, 1]
		for _, elt := range lit.Elts {
			if kv, ok := elt.(*KeyValueExpr); !ok {
				size += 1
			} else if k, ok := kv.Key.(*BasicLit); !ok || k.Kind != token.INT {
				return nil, false, ErrArrayKey

			// The limit of 2^31 elements is infered from the go implementation
			// The actual limit is "largest value representable by an int"
			} else if i, err := strconv.ParseUint(k.Value, 0, 31); err != nil {
				return nil, false, ErrArrayKey
			} else if !(i < size) {
				size = i + 1
			}
		}
		// Allocate the slice
		v.Set(reflect.MakeSlice(t, int(size), int(size)))
	}

	// Fill the array or slice, the reflect interface is identical for both
	for _, elt := range lit.Elts {
		var expr ast.Expr
		if kv, ok := elt.(*KeyValueExpr); !ok {
			expr = elt
		} else {
			// We know this expression to be valid from above.
			curKey, _ = strconv.ParseUint(kv.Key.(*BasicLit).Value, 0, 31)
			expr = kv.Value
		}

		if !(curKey < size) {
			return nil, false, ErrArrayIndexOutOfBounds{t, curKey}
		}

		// Evaluate and set the element
		elem := v.Index(int(curKey))
		if values, typed, err := EvalExpr(ctx, expr, env); err != nil {
			return nil, false, err
		} else if value, err := expectSingleValue(ctx, *values, elt); err != nil {
			return nil, false, err
		} else if err := setTypedValue(elem, value, typed); err != nil {
			return nil, false, err
		}
		curKey += 1
	}
	return &v, true, nil
}

func evalCompositeLitStruct(ctx *Ctx, t reflect.Type, lit *CompositeLit, env *Env) (*reflect.Value, bool, error) {
	vp := reflect.New(t)
	v := vp.Elem()

	if len(lit.Elts) == 0 {
		return &v, true, nil
	}

	_, pairs := lit.Elts[0].(*KeyValueExpr)
	for i, elt := range lit.Elts {
		var field, value *reflect.Value
		var typed bool
		var fname string
		if kv, ok := elt.(*KeyValueExpr); ok != pairs {
			return &v, true, errors.New("Elements are either all key value pairs or not")
		} else if pairs {
			if k, ok := kv.Key.(*Ident); !ok {
				return &v, true, errors.New(fmt.Sprintf("Invalid key node %v %T", kv.Key, kv.Key))
			} else if f := v.FieldByName(k.Name); !f.IsValid() {
				return &v, true, errors.New(t.Name() + " has no field " + k.Name)
			} else {
				fv, ft, err := EvalExpr(ctx, kv.Value, env)
				if err != nil {
					return &v, true, err
				} else if fv == nil {
					return nil, false, nil
				} else {
					field = &f
					value = &(*fv)[0]
					typed = ft
					tfield, _ := v.Type().FieldByName(k.Name)
					fname = tfield.Name
				}
			}
		} else {
			if i >= v.NumField() {
				return &v, true, errors.New("Too many elements for struct " + t.Name())
			} else if _, ok := elt.(*KeyValueExpr); ok {
				return &v, true, errors.New("Elements are either all key value pairs or not")
			} else {
				fv, ft, err := EvalExpr(ctx, elt, env)
				if err != nil {
					return &v, true, err
				} else if fv == nil {
					return &v, false, nil
				} else {
					fieldi := v.Field(i)
					field, value = &fieldi, &(*fv)[0]
					typed = ft
					fname = v.Type().Field(i).Name
				}
			}
		}

		if err := setTypedValue(*field, *value, typed); err != nil {
			return nil, true, errors.New(fmt.Sprintf("cannot set %s.%s : %v",
				t.Name(), fname, value.String()))
		}
	}
	return &v, true, nil
}
