package eval

import (
	"reflect"

	"go/token"
)

func evalUnaryExpr(ctx *Ctx, b *UnaryExpr, env *Env) (r reflect.Value, rtyped bool, err error) {
	var xx *[]reflect.Value
	var xtyped bool
	if xx, xtyped, err = EvalExpr(ctx, b.X.(Expr), env); err != nil {
		return reflect.Value{}, false, err
	}
	rtyped = xtyped
	x := (*xx)[0]

	if userConversion != nil {
		x, xtyped, err = userConversion(x, xtyped)
	}

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err = evalUnaryIntExpr(ctx, x, b.Op)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err = evalUnaryUintExpr(ctx, x, b.Op)
	case reflect.Float32, reflect.Float64:
		r, err = evalUnaryFloatExpr(ctx, x, b.Op)
	case reflect.Complex64, reflect.Complex128:
		r, err = evalUnaryComplexExpr(ctx, x, b.Op)
	case reflect.Bool:
		r, err = evalUnaryBoolExpr(ctx, x, b.Op)
	default:
	        err = ErrInvalidOperand{x, b.Op}
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
