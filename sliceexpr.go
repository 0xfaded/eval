package eval

import (
	"fmt"
	"errors"
	"reflect"
)

// Note: we also need slice3...

func evalSliceExpr(ctx *Ctx, e *SliceExpr, env *Env) (*reflect.Value, bool, error) {

	var cexpr Expr
	var errs []error
	if cexpr, errs = CheckExpr(ctx, e.X, env); len(errs) != 0 {
		for _, cerr := range errs {
			fmt.Printf("%v\n", cerr)
		}
		return nil, false, errors.New("Something wrong checking * expression")
	}

	xs, _, err := EvalExpr(ctx, cexpr, env)
	if err != nil {
		return nil, false, err
	} else if xs == nil {
		// XXX temporary error until typed evaluation of nil
		return nil, false, errors.New("Cannot index nil type")
	}

	var x reflect.Value
	if x, err = expectSingleValue(ctx, *xs, e.X); err != nil {
		return nil, false, err
	}

	var low, high int

	t := x.Type()
	// Special short hand for array pointers
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Array {
		x = x.Elem()
	}

	typ := x.Type()
	switch typ.Kind() {
	case reflect.Slice, reflect.Array, reflect.String:
	default:
		return nil, true, ErrInvalidSliceType{at(ctx, e), x.Type()}
	}

	if e.High != nil {
		if cexpr, errs = CheckExpr(ctx, e.High, env); len(errs) != 0 {
			for _, cerr := range errs {
				fmt.Printf("%v\n", cerr)
			}
			return nil, false, errors.New("Something wrong checking slice ub")
		}
		if high, err = evalIntIndex(ctx, cexpr, env, x.Type()); err != nil {
			return nil, false, err
		}
		if err := CannotIndex(x, high); err != nil {
			return nil, false, err
		}

	}

	if e.Low != nil {
		if cexpr, errs = CheckExpr(ctx, e.Low, env); len(errs) != 0 {
			for _, cerr := range errs {
				fmt.Printf("%v\n", cerr)
			}
			return nil, false, errors.New("Something wrong checking slice lb")
		}
		if low, err = evalIntIndex(ctx, cexpr, env, x.Type()); err != nil {
			return nil, false, err
		}
		if err := CannotIndex(x, low); err != nil {
			return nil, false, err
		}
	}

	if low > high {
		return nil, false, errors.New(fmt.Sprintf("invalid slice index: %d > %d", low, high))
	}
	slice := x.Slice(low, high)
	return &slice, true, nil
}
