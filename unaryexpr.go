package interactive

import (
	"reflect"

	"go/ast"
	"go/token"
)

func evalUnaryExpr(ctx *Ctx, b *ast.UnaryExpr, env *Env) (r reflect.Value, rtyped bool, err error) {
	var xx *[]reflect.Value
	var xtyped bool
	if xx, xtyped, err = EvalExpr(ctx, b.X, env); err != nil {
		return reflect.Value{}, false, err
	}
	rtyped = xtyped
	x := (*xx)[0]

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err = evalUnaryIntExpr(ctx, x, b.Op)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err = evalUnaryUintExpr(ctx, x, b.Op)
	case reflect.Float32, reflect.Float64:
		r, err = evalUnaryFloatExpr(ctx, x, b.Op)
	case reflect.Complex64, reflect.Complex128:
		r, err = evalUnaryComplexExpr(ctx, x, b.Op)
	case reflect.String:
		r, err = evalUnaryStringExpr(ctx, x, b.Op)
	default:
		err = ErrInvalidOperands{x, b.Op, x}
	}
	return
}

// Assumes y is assignable to x, panics otherwise
func evalUnaryIntExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var r int64
	var err error
	// var b bool
	is_bool := false

	var b bool = false

	xx := x.Int()
	switch op {
	case token.ADD: r = +xx
	case token.SUB: r = -xx
	default: err = ErrInvalidOperand{x, op}
	}
	if is_bool {
		return reflect.ValueOf(b), err
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), err
	}
}

// Assumes y is assignable to x, panics otherwise
func evalUnaryUintExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var err error
	var r uint64
	var b bool
	is_bool := false

	xx := x.Uint()
	switch op {
	case token.ADD: r = +xx
	// case token.SUB: r = -xx
	default: err = ErrInvalidOperand{x, op}
	}
	if is_bool {
		return reflect.ValueOf(b), err
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), err
	}
}

// Assumes y is assignable to x, panics otherwise
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

// Assumes y is assignable to x, panics otherwise
func evalUnaryComplexExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var err error
	var r complex128

	// xx := x.Complex()
	switch op {
	default: err = ErrInvalidOperand{x, op}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

// Assumes y is assignable to x, panics otherwise
func evalUnaryStringExpr(ctx *Ctx, x reflect.Value, op token.Token) (reflect.Value, error) {
	var err error
	var r string

	// xx := x.String()
	switch op {
	default: err = ErrInvalidOperand{x, op}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}
