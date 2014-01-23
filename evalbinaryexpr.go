package eval

import (
	"reflect"
	"go/token"
)

func evalBinaryExpr(ctx *Ctx, b *BinaryExpr, env *Env) (r reflect.Value, rtyped bool, err error) {

        if b.IsConst() {
                return b.Const(), true, nil
        }

        xexpr := b.X.(Expr)
        yexpr := b.Y.(Expr)

        // Compute the operand type
        // TODO[crc] I have decided that const nodes with inferred types
        // need to be retyped to avoid logic like below.
        var zt []reflect.Type
        if xexpr.IsConst() {
                zt = yexpr.KnownType()
        } else {
                zt = xexpr.KnownType()
        }

        var xs, ys []reflect.Value
        if xs, err = evalTypedExpr(ctx, xexpr, zt, env); err != nil {
                return reflect.Value{}, false, err
        } else if ys, err = evalTypedExpr(ctx, yexpr, zt, env); err != nil {
                return reflect.Value{}, false, err
        }
        x, y := xs[0], ys[0]

	switch zt[0].Kind() {
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
	case reflect.Bool:
		r, err = evalBinaryBoolExpr(ctx, x, b.Op, y)
	default:
                panic("eval: unimplemented binary ops")
	}
	return r, true, err
}

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
	case token.QUO:
                if yy == 0 {
                        return reflect.Value{}, PanicDivideByZero{}
                }
                r = xx / yy
	case token.REM:
                if yy == 0 {
                        return reflect.Value{}, PanicDivideByZero{}
                }
                r = xx % yy
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
	default:
		panic(dytc("bad binary op"))
	}
	if is_bool {
		return reflect.ValueOf(b), err
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), err
	}
}

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
	case token.QUO:
                if yy == 0 {
                        return reflect.Value{}, PanicDivideByZero{}
                }
                r = xx / yy
	case token.REM:
                if yy == 0 {
                        return reflect.Value{}, PanicDivideByZero{}
                }
                r = xx % yy
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
	default:
		panic(dytc("bad binary op"))
	}
	if is_bool {
		return reflect.ValueOf(b), err
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), err
	}
}

func evalBinaryFloatExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var r float64
        var is_bool bool
	var b bool

	xx, yy := x.Float(), y.Float()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO:
                if yy == 0 {
                        return reflect.Value{}, PanicDivideByZero{}
                }
                r = xx / yy
	case token.EQL: b = xx == yy; is_bool = true
	case token.NEQ: b = xx != yy; is_bool = true
	case token.LEQ: b = xx <= yy; is_bool = true
	case token.GEQ: b = xx >= yy; is_bool = true
	case token.LSS: b = xx < yy;  is_bool = true
	case token.GTR: b = xx > yy;  is_bool = true
	default:
		panic(dytc("bad binary op"))
	}
	if is_bool {
		return reflect.ValueOf(b), nil
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), nil
	}
}

func evalBinaryComplexExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var r complex128
        var is_bool bool
	var b bool

	xx, yy := x.Complex(), y.Complex()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO:
                if yy == 0 {
                        return reflect.Value{}, PanicDivideByZero{}
                }
                r = xx / yy
	case token.EQL: b = xx == yy; is_bool = true
	case token.NEQ: b = xx != yy; is_bool = true
	default:
		panic(dytc("bad binary op"))
	}
	if is_bool {
		return reflect.ValueOf(b), nil
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), nil
	}
}

func evalBinaryStringExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
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
	default:
		panic(dytc("bad binary op"))
	}
	if is_bool {
		return reflect.ValueOf(b), nil
	} else {
		return reflect.ValueOf(r).Convert(x.Type()), nil
	}
}

func evalBinaryBoolExpr(ctx *Ctx, x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	xx, yy := x.Bool(), y.Bool()
	var r bool
	switch op {
	case token.LAND: r = xx && yy
	case token.LOR: r = xx || yy
	case token.EQL: r = xx == yy
	case token.NEQ: r = xx != yy
	default:
		panic(dytc("bad binary op"))
	}
        return reflect.ValueOf(r), nil
}
