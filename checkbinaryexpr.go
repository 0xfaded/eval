package interactive

import (
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

	if errs == nil {
		ax := aexpr.X.(Expr)
		ay := aexpr.Y.(Expr)

		tx := ax.KnownType()
		ty := ay.KnownType()

		// Check for multi valued expressions. Not much we can do if we find one
		// TODO check for single values

		// Check for compatible types

		// TODO tx and ty will always have a known type once checker is complete
		//      This if() is a shim
		if len(tx) != 1 || len(ty) != 1 {
			return aexpr, errs
		}

		cx, okx := tx[0].(ConstType)
		cy, oky := ty[0].(ConstType)
		if ax.IsConst() && ay.IsConst() {
			if okx && oky {
				if promoted, err := promoteConsts(ctx, cx, cy, ay); err != nil {
					errs = append(errs, err)
					errs = append(errs, ErrInvalidBinaryOperation{at(ctx, binary)})
				} else {
					aexpr.knownType = knownType{promoted}
					aexpr.constValue, moreErrs = evalConstBinaryExpr(ctx, aexpr, promoted)
					if moreErrs != nil {
						errs = append(errs, moreErrs...)
					}
				}
			}
		}
	}
	return aexpr, errs
}

// Evaluates a const binary Expr. May return a sensical constValue
// even if ErrTruncatedConst errors are present
func evalConstBinaryExpr(ctx *Ctx, constExpr *BinaryExpr, promotedType ConstType) (constValue, []error) {
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
		return constValueOf(NewBigRune(0).Add(x, y)), nil
	case token.SUB:
		return constValueOf(NewBigRune(0).Sub(x, y)), nil
	case token.MUL:
		return constValueOf(NewBigRune(0).Mul(x, y)), nil
	case token.QUO:
		if y.IsZero() {
			return constValue{}, []error{ErrDivideByZero{at(ctx, constExpr.Y)}}
		}
		return constValueOf(NewBigRune(0).Quo(x, y)), nil
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
			errs = append(errs, ErrTruncatedConstant{at(ctx, constExpr.X), x})
		}
		if yy, trunc = y.Integer(); trunc {
			errs = append(errs, ErrTruncatedConstant{at(ctx, constExpr.Y), y})
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
		return constValueOf(x.Rat.Cmp(y.Rat) == 0 && x.Imag.Cmp(y.Imag) == 0), nil
	case token.NEQ:
		return constValueOf(x.Rat.Cmp(y.Rat) != 0 || x.Imag.Cmp(y.Imag) != 0), nil

	case token.LEQ, token.GEQ, token.LSS, token.GTR:
		var b bool
		var trunc bool
		if xx, trunc = x.Real(); trunc {
			errs = append(errs, ErrTruncatedConstant{at(ctx, constExpr.X), x})
		}
		if yy, trunc = y.Real(); trunc {
			errs = append(errs, ErrTruncatedConstant{at(ctx, constExpr.Y), y})
		}
		cmp := xx.Rat.Cmp(yy.Rat)
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

