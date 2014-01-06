package eval

import (
	"errors"
	"reflect"
	"fmt"
)

func evalStarExpr(ctx *Ctx, starExpr *StarExpr, env *Env) (*reflect.Value, bool, error) {
	// return nil, false, errors.New("Star expressions not done yet")
	var cexpr Expr
	var errs []error
	if cexpr, errs = CheckExpr(ctx, starExpr.X, env); len(errs) != 0 {
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
		return nil, false, errors.New("Cannot dereferece nil type")
	}

	var x reflect.Value
	if x, err = expectSingleValue(ctx, *xs, starExpr.X); err != nil {
		return nil, false, err
	}

	switch x.Type().Kind() {
	case reflect.Interface, reflect.Ptr:
		val := x.Elem()
		return &val, true, nil
	default:
		return nil, true, ErrInvalidIndirect{x.Type()}
	}
}
