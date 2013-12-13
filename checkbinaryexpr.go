package interactive

import (
	"reflect"

	"go/ast"
	"go/token"
)

func checkBinaryExpr(ctx *Ctx, binary *ast.BinaryExpr, env *Env) (aexpr *BinaryExpr, errs []error) {
	aexpr = &BinaryExpr{BinaryExpr: binary}

	var moreErrs []error
	if aexpr.X, moreErrs = checkExpr(ctx, binary.X, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	if aexpr.Y, moreErrs = checkExpr(ctx, binary.Y, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}

	if errs != nil {
		return aexpr, errs
	}

	xa := aexpr.X.(Expr)
	ya := aexpr.Y.(Expr)

	xt := xa.KnownType()
	yt := ya.KnownType()

	// Check for multi valued expressions. Not much we can do if we find one
	// TODO check for single values

	// Check for compatible types

	// TODO tx and ty will always have a known type once checker is complete
	//      This if() is a shim
	if len(xt) != 1 || len(yt) != 1 {
		return aexpr, errs
	}

	xc, xuntyped := xt[0].(ConstType)
	yc, yuntyped := yt[0].(ConstType)
	if xa.IsConst() && ya.IsConst() {
		if xuntyped && yuntyped {
			yv := ya.Const()
			if promoted, err := promoteConsts(ctx, xc, yc, ya, yv); err != nil {
				errs = append(errs, err)
				errs = append(errs, ErrInvalidBinaryOperation{at(ctx, binary)})
			} else {
				if isBooleanOp(binary.Op) {
					aexpr.knownType = []reflect.Type{ConstBool}
				} else {
					aexpr.knownType = knownType{promoted}
				}
				aexpr.constValue, moreErrs = evalConstUntypedBinaryExpr(ctx, aexpr, promoted)
				if moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
			}
		} else if yuntyped {
			if z, moreErrs := evalConstTypedUntypedBinaryExpr(ctx, aexpr, xa, ya); moreErrs!= nil {
				errs = append(errs, moreErrs...)
			} else {
				if isBooleanOp(binary.Op) {
					aexpr.knownType = []reflect.Type{ConstBool}
				} else {
					aexpr.knownType = xt
				}
				aexpr.constValue = z
			}
		} else if xuntyped {
			if z, moreErrs := evalConstTypedUntypedBinaryExpr(ctx, aexpr, ya, xa); moreErrs!= nil {
				errs = append(errs, moreErrs...)
			} else {
				if isBooleanOp(binary.Op) {
					aexpr.knownType = []reflect.Type{ConstBool}
				} else {
					aexpr.knownType = yt
				}
				aexpr.constValue = z
			}
		} else {
			panic("Unimplemented")
		}
	}
	return aexpr, errs
}

// Evaluates a const binary Expr. May return a sensical constValue
// even if ErrTruncatedConst errors are present
func evalConstUntypedBinaryExpr(ctx *Ctx, constExpr *BinaryExpr, promotedType ConstType) (constValue, []error) {
	x := constExpr.X.(Expr).Const()
	y := constExpr.Y.(Expr).Const()
	switch promotedType.(type) {
	case ConstIntType, ConstRuneType, ConstFloatType, ConstComplexType:
		xx := x.Interface().(*BigComplex)
		yy := y.Interface().(*BigComplex)
		return evalConstBinaryNumericExpr(ctx, constExpr, xx, yy)
	case ConstStringType:
		xx := x.String()
		yy := y.String()
		return evalConstBinaryStringExpr(ctx, constExpr, xx, yy)
	case ConstBoolType:
		xx := x.Bool()
		yy := y.Bool()
		return evalConstBinaryBoolExpr(ctx, constExpr, xx, yy)
	default:
		// It is possible that both x and y are ConstNil, however no operator is defined, not even ==
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
	}

}

func evalConstBinaryNumericExpr(ctx *Ctx, constExpr *BinaryExpr, x, y *BigComplex) (constValue, []error) {
	var errs []error
	var xx *BigComplex
	var yy *BigComplex

	switch constExpr.Op {
	case token.ADD:
		return constValueOf(new(BigComplex).Add(x, y)), nil
	case token.SUB:
		return constValueOf(new(BigComplex).Sub(x, y)), nil
	case token.MUL:
		return constValueOf(new(BigComplex).Mul(x, y)), nil
	case token.QUO:
		if y.IsZero() {
			return constValue{}, []error{ErrDivideByZero{at(ctx, constExpr.Y)}}
		}
		return constValueOf(new(BigComplex).Quo(x, y)), nil
	case token.REM:
		if y.IsZero() {
			return constValue{}, []error{ErrDivideByZero{at(ctx, constExpr.Y)}}
		} else if !(x.IsInteger() && y.IsInteger()) {
			return constValue{}, []error{ErrIllegalConstantExpr{at(ctx, constExpr)}}
		} else {
			z := NewBigRune(1)
			z.Rat.Num().Rem(x.Num(), y.Num())
			return constValueOf(z), nil
		}
	case token.AND, token.OR, token.XOR, token.AND_NOT:
		var trunc bool
		if xx, trunc = x.Integer(); trunc {
			errs = append(errs, ErrTruncatedConstant{at(ctx, constExpr.X), ConstFloat, x})
		}
		if yy, trunc = y.Integer(); trunc {
			errs = append(errs, ErrTruncatedConstant{at(ctx, constExpr.Y), ConstFloat, y})
		}

		z := NewBigRune(1)
		switch constExpr.Op {
		case token.AND:
			z.Num().And(xx.Num(), yy.Num())
		case token.OR:
			z.Num().Or(xx.Num(), yy.Num())
		case token.XOR:
			z.Num().Xor(xx.Num(), yy.Num())
		case token.AND_NOT:
			z.Num().AndNot(xx.Num(), yy.Num())
		}
		return constValueOf(z), errs

	case token.EQL:
		return constValueOf(x.Rat.Cmp(&y.Rat) == 0 && x.Imag.Cmp(&y.Imag) == 0), nil
	case token.NEQ:
		return constValueOf(x.Rat.Cmp(&y.Rat) != 0 || x.Imag.Cmp(&y.Imag) != 0), nil

	case token.LEQ, token.GEQ, token.LSS, token.GTR:
		var b bool
		var trunc bool
		if xx, trunc = x.Real(); trunc {
			errs = append(errs, ErrTruncatedConstant{at(ctx, constExpr.X), ConstFloat, x})
		}
		if yy, trunc = y.Real(); trunc {
			errs = append(errs, ErrTruncatedConstant{at(ctx, constExpr.Y), ConstFloat, y})
		}
		cmp := xx.Rat.Cmp(&yy.Rat)
		switch constExpr.Op {
		case token.NEQ:
			b = cmp != 0
		case token.LEQ:
			b = cmp <= 0
		case token.GEQ:
			b = cmp >= 0
		case token.LSS:
			b = cmp < 0
		case token.GTR:
			b = cmp > 0
		}
		return constValueOf(b), errs
	default:
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
	}
}

func evalConstBinaryStringExpr(ctx *Ctx, constExpr *BinaryExpr, x, y string) (constValue, []error) {
	switch constExpr.Op {
	case token.ADD:
		return constValueOf(x + y), nil
	case token.NEQ:
		return constValueOf(x != y), nil
	case token.LEQ:
		return constValueOf(x <= y), nil
	case token.GEQ:
		return constValueOf(x >= y), nil
	case token.LSS:
		return constValueOf(x < y), nil
	case token.GTR:
		return constValueOf(x > y), nil
	default:
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
	}
}

func evalConstBinaryBoolExpr(ctx *Ctx, constExpr *BinaryExpr, x, y bool) (constValue, []error) {
	switch constExpr.Op {
	case token.EQL:
		return constValueOf(x == y), nil
	case token.NEQ:
		return constValueOf(x != y), nil
	case token. LAND:
		return constValueOf(x && y), nil
	case token.LOR:
		return constValueOf(x || y), nil
	default:
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
	}
}

// Evaluate x op y
func evalConstTypedUntypedBinaryExpr(ctx *Ctx, expr *BinaryExpr, typedExpr, untypedExpr Expr) (
	constValue, []error) {

	xt := untypedExpr.KnownType()[0]
	yt := typedExpr.KnownType()[0]

	switch xt.(type) {
	case ConstIntType, ConstRuneType, ConstFloatType, ConstComplexType:
		x := untypedExpr.Const().Interface().(*BigComplex)
		var y *BigComplex
		switch yt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			y = NewBigInt64(typedExpr.Const().Int())

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			y = NewBigUint64(typedExpr.Const().Uint())

		case reflect.Float32, reflect.Float64:
			y = NewBigFloat64(typedExpr.Const().Float())

		case reflect.Complex64, reflect.Complex128:
			y = NewBigComplex128(typedExpr.Const().Complex())

		default:
			// This will result in a bad conversion error
			_, errs := convertConstToTyped(ctx, xt, constValueOf(x), yt, untypedExpr)
			return constValue{}, errs
		}

		z, errs := evalConstBinaryNumericExpr(ctx, expr, x, y)
		r, moreErrs := convertConstToTyped(ctx, xt, z, yt, expr)
		errs = append(errs, moreErrs...)
		return constValue(r), errs

	case ConstStringType:
		if yt.Kind() == reflect.String {
			xstring := untypedExpr.Const().String()
			ystring := typedExpr.Const().String()
			z, errs := evalConstBinaryStringExpr(ctx, expr, xstring, ystring)
			r, moreErrs := convertConstToTyped(ctx, xt, z, yt, expr)
			errs = append(errs, moreErrs...)
			return constValue(r), errs
		}

	case ConstBoolType:
		if yt.Kind() == reflect.Bool {
			xbool := untypedExpr.Const().Bool()
			ybool := typedExpr.Const().Bool()
			z, errs := evalConstBinaryBoolExpr(ctx, expr, xbool, ybool)
			r, moreErrs := convertConstToTyped(ctx, xt, z, yt, expr)
			errs = append(errs, moreErrs...)
			return constValue(r), errs
		}

	case ConstNilType:
		panic("func")
	}
	panic("func")

}
