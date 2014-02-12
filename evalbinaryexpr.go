package eval

import (
	"reflect"
	"go/token"

	"errors"
)

func evalBinaryExpr(binary *BinaryExpr, env Env) (r reflect.Value, err error) {

        if binary.IsConst() {
                return binary.Const(), nil
        }

        xexpr := binary.X.(Expr)
        yexpr := binary.Y.(Expr)

        var zt []reflect.Type
        if xexpr.IsConst() && xexpr.KnownType()[0].Kind() != reflect.Interface {
                zt = yexpr.KnownType()
        } else {
                zt = xexpr.KnownType()
        }

        var xs, ys []reflect.Value
        if xs, err = evalTypedExpr(xexpr, zt, env); err != nil {
                return reflect.Value{}, err
        } else if ys, err = evalTypedExpr(yexpr, zt, env); err != nil {
                return reflect.Value{}, err
        }
        x, y := xs[0], ys[0]

	var b bool
	switch zt[0].Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err = evalBinaryIntExpr(x, binary.Op, y)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err = evalBinaryUintExpr(x, binary.Op, y)
	case reflect.Float32, reflect.Float64:
		r, err = evalBinaryFloatExpr(x, binary.Op, y)
	case reflect.Complex64, reflect.Complex128:
		r, err = evalBinaryComplexExpr(x, binary.Op, y)
	case reflect.String:
		r, err = evalBinaryStringExpr(x, binary.Op, y)
	case reflect.Bool:
		r, err = evalBinaryBoolExpr(x, binary.Op, y)
	case reflect.Interface, reflect.Ptr:
		if xexpr.KnownType()[0] == ConstNil {
			b = y.IsNil()
		} else if yexpr.KnownType()[0] == ConstNil {
			b = x.IsNil()
		} else if t := areDynamicTypesComparable(x, y); t != nil {
			return reflect.Value{}, PanicUncomparableType{t}
		} else {
			b = x.Interface() == y.Interface()
		}
		if binary.Op == token.NEQ {
			b = !b
		}
		r = reflect.ValueOf(b)
	case reflect.Struct, reflect.Array:
		if t := areDynamicTypesComparable(x, y); t != nil {
			return reflect.Value{}, PanicUncomparableType{t}
		}
		b = x.Interface() == y.Interface()
		if binary.Op == token.NEQ {
			b = !b
		}
		r = reflect.ValueOf(b)
	case reflect.Map, reflect.Slice, reflect.Func:
		if xexpr.KnownType()[0] == ConstNil {
			b = y.IsNil()
		} else {
			b = x.IsNil()
		}
		if binary.Op == token.NEQ {
			b = !b
		}
		r = reflect.ValueOf(b)
	default:
                return reflect.Value{}, errors.New("eval: unimplemented binary ops :(")
	}
	return r, err
}

func evalBinaryIntExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
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

func evalBinaryUintExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
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

func evalBinaryFloatExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
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

func evalBinaryComplexExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
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

func evalBinaryStringExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
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

func evalBinaryBoolExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
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

func areDynamicTypesComparable(x, y reflect.Value) reflect.Type {
	if x.Type() != y.Type() {
		return nil
	}
	switch x.Type().Kind() {
	case reflect.Interface:
		return areDynamicTypesComparable(x.Elem(), y.Elem())
	case reflect.Struct:
		numField := x.NumField()
		for i := 0; i < numField; i += 1 {
			if t := areDynamicTypesComparable(x.Field(i), y.Field(i)); t != nil {
				if t.Kind() == reflect.Struct {
					return t
				} else {
					return x.Type()
				}
			}
		}
	case reflect.Array:
		length := x.Len()
		for i := 0; i < length; i += 1 {
			if t := areDynamicTypesComparable(x.Index(i), y.Index(i)); t != nil {
				return t
			}
		}
	case reflect.Map, reflect.Func, reflect.Slice:
		return x.Type()
	}
	return nil
}

