package interactive

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"go/ast"
	"go/token"
)

func evalCompositeLit(lit *ast.CompositeLit, env *Env) (reflect.Value, bool, error) {
	t, err := evalType(lit.Type, env)
	if err != nil {
		return reflect.Value{}, true, err
	}

	switch t.Kind() {
	//case reflect.Map:
	case reflect.Array, reflect.Slice:
		return evalCompositeLitArrayOrSlice(t, lit, env)
	case reflect.Struct:
		return evalCompositeLitStruct(t, lit, env)
	default:
		return reflect.Value{}, true, errors.New(
			fmt.Sprintf("invalid type for composite literal %s", t.Name()))

	}
}

func evalCompositeLitArrayOrSlice(t reflect.Type, lit *ast.CompositeLit, env *Env) (reflect.Value, bool, error) {

	v := reflect.New(t).Elem()

	var curKey uint64 = 0
	var size uint64 = 0
	if t.Kind() == reflect.Array {
		size = uint64(t.Len())
	} else {
		// Check all keys are valid and calculate slice size
		for _, elt := range lit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); !ok {
				size += 1
			} else if k, ok := kv.Key.(*ast.BasicLit); !ok || k.Kind != token.INT {
				return reflect.Value{}, false, ErrArrayKey

			// The limit of 2^31 elements is infered from the go implementation
			} else if i, err := strconv.ParseUint(k.Value, 0, 31); err != nil {
				return reflect.Value{}, false, ErrArrayKey
			} else if !(i < size) {
				size = i + 1
			}
		}
		v.Set(reflect.MakeSlice(t, int(size), int(size)))
	}

	for _, elt := range lit.Elts {
		var expr ast.Expr
		if kv, ok := elt.(*ast.KeyValueExpr); !ok {
			expr = elt
		} else {
			// We know this expression to be valid from above.
			curKey, _ = strconv.ParseUint(kv.Key.(*ast.BasicLit).Value, 0, 31)
			expr = kv.Value
		}

		if !(curKey < size) {
			return reflect.Value{}, false, ErrArrayIndexOutOfBounds{t, curKey}
		}

		elem := v.Index(int(curKey))
		if value, typed, err := evalExpr(expr, env); err != nil {
			return reflect.Value{}, false, err
		} else if len(value) == 0 {
			return reflect.Value{}, false, ErrMissingValue
		} else if len(value) > 1 {
			return reflect.Value{}, false, ErrMultiInSingleContext{value}
		} else if err := setTypedValue(elem, value[0], typed); err != nil {
			return reflect.Value{}, false, err
		}
		curKey += 1
	}
	return v, true, nil
}

func evalCompositeLitStruct(t reflect.Type, lit *ast.CompositeLit, env *Env) (reflect.Value, bool, error) {
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
			} else if fv, ft, err := evalExpr(kv.Value, env); err != nil {
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
			} else if fv, ft, err := evalExpr(elt, env); err != nil {
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
}
