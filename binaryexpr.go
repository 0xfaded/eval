package eval

import (
	"reflect"
	"go/token"
)

func evalBinaryExpr(ctx *Ctx, b *BinaryExpr, env *Env) (r reflect.Value, rtyped bool, err error) {
	var xx, yy *[]reflect.Value
	var xtyped, ytyped bool
	if xx, xtyped, err = EvalExpr(ctx, b.X.(Expr), env); err != nil {
		return reflect.Value{}, false, err
	}

	if yy, ytyped, err = EvalExpr(ctx, b.Y.(Expr), env); err != nil {
		return reflect.Value{}, false, err
	}
	rtyped = xtyped || ytyped
	x, y := (*xx)[0], (*yy)[0]

	if userConversion != nil {
		x, xtyped, err = userConversion(x, xtyped)
		y, ytyped, err = userConversion(y, xtyped)
	}

	// Rearrange x and y such that y is assignable to x, if possible
	if xtyped && ytyped {
		if x.Type().AssignableTo(y.Type()) {
			x = x.Convert(y.Type())
		} else if !y.Type().AssignableTo(x.Type()) {
			return r, rtyped, ErrInvalidOperands{x, b.Op, y}
		}
	} else if xtyped {
		if !y.Type().ConvertibleTo(x.Type()) {
			return r, rtyped, ErrInvalidOperands{x, b.Op, y}
		}
		y = y.Convert(x.Type())
	} else if ytyped {
		if !x.Type().ConvertibleTo(y.Type()) {
			return r, rtyped, ErrInvalidOperands{x, b.Op, y}
		}
		x = x.Convert(y.Type())
	} else if isUntypedNumeral(x) && isUntypedNumeral(y) {
		x, y = promoteUntypedNumerals(x, y)
	} else {
		return r, rtyped, ErrInvalidOperands{x, b.Op, y}
	}

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err = evalBinaryIntExpr(ctx, x, b.Op, y)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err = evalBinaryUintExpr(ctx, x, b.Op, y)
	case reflect.Float32, reflect.Float64:
		r, err = evalBinaryFloatExpr(ctx, x, b.Op, y)
	case reflect.Complex64, reflect.Complex128:
		r, err = evalBinaryComplexExpr(ctx, x, b.Op, y)
	case reflect.String:
		r, err = evalBinaryStringExpr(ctx, x, b.Op, y)
	default:
		err = ErrInvalidOperands{x, b.Op, y}
	}
	return
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryIntExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var r int64
	var err error
	var b bool
	is_bool := false

	xx, yy := x.Int(), y.Int()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO: r = xx / yy
	case token.REM: r = xx % yy
	case token.AND: r = xx & yy
	case token.OR:  r = xx | yy
	case token.XOR: r = xx ^ yy
	case token.AND_NOT: r = xx &^ yy
	case token.EQL: b = xx == yy; is_bool = true
	case token.NEQ: b = xx != yy; is_bool = true
	case token.LEQ: b = xx <= yy; is_bool = true
	case token.GEQ: b = xx >= yy; is_bool = true
	case token.LSS: b = xx < yy;  is_bool = true
	case token.GTR: b = xx > yy;  is_bool = true
	default: err = ErrInvalidOperands{x, op, y}
	}
	if is_bool {
		return reflect.ValueOf(b), err
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), err
	}
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryUintExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var err error
	var r uint64
	var b bool
	is_bool := false

	xx, yy := x.Uint(), y.Uint()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO: r = xx / yy
	case token.REM: r = xx % yy
	case token.AND: r = xx & yy
	case token.OR:  r = xx | yy
	case token.XOR: r = xx ^ yy
	case token.AND_NOT: r = xx &^ yy
	case token.EQL: b = xx == yy; is_bool = true
	case token.NEQ: b = xx != yy; is_bool = true
	case token.LEQ: b = xx <= yy; is_bool = true
	case token.GEQ: b = xx >= yy; is_bool = true
	case token.LSS: b = xx < yy;  is_bool = true
	case token.GTR: b = xx > yy;  is_bool = true
	default: err = ErrInvalidOperands{x, op, y}
	}
	if is_bool {
		return reflect.ValueOf(b), err
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), err
	}
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryFloatExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var err error
	var r float64

	xx, yy := x.Float(), y.Float()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO: r = xx / yy
	// case token.EQL: b = xx == yy
	// case token.LSS: b = xx < yy
	// case token.GTR: b = xx > yy
	default: err = ErrInvalidOperands{x, op, y}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryComplexExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var err error
	var r complex128

	xx, yy := x.Complex(), y.Complex()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO: r = xx / yy
	default: err = ErrInvalidOperands{x, op, y}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryStringExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var err error
	var r string
	var b bool
	is_bool := false

	xx, yy := x.String(), y.String()
	switch op {
	case token.ADD:	r = xx + yy
	case token.EQL: b = xx == yy; is_bool = true
	case token.NEQ: b = xx != yy; is_bool = true
	case token.LEQ: b = xx <= yy; is_bool = true
	case token.GEQ: b = xx >= yy; is_bool = true
	case token.LSS: b = xx < yy;  is_bool = true
	case token.GTR: b = xx > yy;  is_bool = true
	default: err = ErrInvalidOperands{x, op, y}
	}
	if is_bool {
		return reflect.ValueOf(b), err
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), err
	}
}
