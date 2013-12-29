package eval

import (
	"errors"
	"fmt"
	"reflect"
	"go/ast"
)

// CannotIndex returns error if i is a not valid index in v
// nil is returned if it is valid.
// It would be nicer if the reflect package provided a "CanIndex" method.
func CannotIndex(v reflect.Value, i int) (err error) {
	err = nil
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(fmt.Sprintf("%v", x))
		}
	}()
	v.Index(i)
	return err
}

func evalIndexExpr(ctx *Ctx, index *IndexExpr, env *Env) (*reflect.Value, bool, error) {
	xs, _, err := EvalExpr(ctx, index.X.(Expr), env)
	if err != nil {
		return nil, false, err
	} else if xs == nil {
		// XXX temporary error until typed evaluation of nil
		return nil, false, errors.New("Cannot index nil type")
	}

	var x reflect.Value
	if x, err = expectSingleValue(ctx, *xs, index.X); err != nil {
		return nil, false, err
	}

	t := x.Type()
	// Special short hand for array pointers
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Array {
		x = x.Elem()
	}

	switch x.Type().Kind() {
	// case reflect.Map:
	case reflect.Array, reflect.Slice, reflect.String:
		return evalIndexExprInt(ctx, x, index.Index, env)
	default:
		return nil, true, ErrInvalidIndexOperation{at(ctx, index), x.Type()}
	}
}

// For arrays, slices and strings
func evalIndexExprInt(ctx *Ctx, x reflect.Value, intExpr ast.Expr, env *Env) (*reflect.Value, bool, error) {
	if i, err := evalIntIndex(ctx, intExpr, env, x.Type()); err != nil {
		return nil, false, err
	} else {
		if err := CannotIndex(x, i); err != nil {
			return nil, false, err
		}
		v := x.Index(i)
		return &v, true, nil
	}
}

func evalIntIndex(ctx *Ctx, intExpr ast.Expr, env *Env, containerType reflect.Type) (int, error) {
	if is, typed, err := EvalExpr(ctx, intExpr.(Expr), env); err != nil {
		return -1, err
	} else if is == nil {
		// XXX temporary error until typed evaluation of nil
		return -1, errors.New("Cannot index nil type")
	} else if i, err := expectSingleValue(ctx, *is, intExpr); err != nil {
		return -1, err

	// XXX This untyped constant conversion is not correct, the index must evaluate exactly
	// to an integer.
	// a[2*0.5] is legal, a[2*0.4] is not.
	//
	// There is also the constraint that constant expressions such as "abc"[10] must be
	// in range. Fix both when constant expressions are correctly handled
	} else if !typed && i.Type().ConvertibleTo(reflect.TypeOf(int(0))) {
		result := int(i.Convert(reflect.TypeOf(int(0))).Int())
		if 0 <= result && (containerType.Kind() != reflect.Array || result < containerType.Len()) {
			return result, nil
		}
		return -1, ErrInvalidIndex{at(ctx, intExpr), reflect.ValueOf(result), containerType}
	} else {
		var result int
		switch i.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result = int(i.Int())
			if result >= 0 {
				return result, nil
			}
			return -1, ErrInvalidIndex{at(ctx, intExpr), reflect.ValueOf(result), containerType}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = int(i.Uint())
			if result >= 0 {
				return result, nil
			}
			return -1, ErrInvalidIndex{at(ctx, intExpr), reflect.ValueOf(result), containerType}
		default:
			return -1, ErrInvalidIndex{at(ctx, intExpr), i, containerType}
		}
	}
}
