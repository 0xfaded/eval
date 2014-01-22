package eval

import (
	"reflect"

	"go/token"
)

func evalUnaryExpr(ctx *Ctx, unary *UnaryExpr, env *Env) (r reflect.Value, err error) {
	if unary.IsConst() {
		return unary.Const(), nil
	}

	var xx *[]reflect.Value
	if xx, _, err = EvalExpr(ctx, unary.X.(Expr), env); err != nil {
		return reflect.Value{}, err
	}
	x := (*xx)[0]

	// handle & and <- first
	if unary.Op == token.AND {
		return x.Addr(), nil
	}
	// TODO handle <-

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err = evalUnaryIntExpr(ctx, x, unary.Op)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err = evalUnaryUintExpr(ctx, x, unary.Op)
	case reflect.Float32, reflect.Float64:
		r, err = evalUnaryFloatExpr(ctx, x, unary.Op)
	case reflect.Complex64, reflect.Complex128:
		r, err = evalUnaryComplexExpr(ctx, x, unary.Op)
	case reflect.Bool:
		r, err = evalUnaryBoolExpr(ctx, x, unary.Op)
	default:
	        err = ErrInvalidOperand{x, unary.Op}
	}
	return
}

func evalUnaryIntExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var r int64
	var err error

	xx := x.Int()
	switch op {
	case token.ADD: r = +xx
	case token.SUB: r = -xx
	case token.XOR: r = ^xx
	default: err = ErrInvalidOperand{x, op}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

func evalUnaryUintExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var err error
	var r uint64

	xx := x.Uint()
	switch op {
	case token.ADD: r = +xx
	case token.SUB: r = -xx
	case token.XOR: r = ^xx
	default: err = ErrInvalidOperand{x, op}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

func evalUnaryFloatExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var err error
	var r float64

	xx := x.Float()
	switch op {
	case token.ADD: r = + xx
	case token.SUB: r = - xx
	default: err = ErrInvalidOperand{x, op}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

func evalUnaryComplexExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var err error
	var r complex128

	xx := x.Complex()
	switch op {
	case token.ADD: r = +xx
	case token.SUB: r = -xx
	default: err = ErrInvalidOperand{x, op}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

func evalUnaryBoolExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var err error
	var r bool

	xx := x.Bool()
	switch op {
	case token.NOT: r = !xx
	default: err = ErrInvalidOperand{x, op}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}
