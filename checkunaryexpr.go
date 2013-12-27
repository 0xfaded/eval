package eval

import (
	"go/ast"
	"go/token"
)

func checkUnaryExpr(ctx *Ctx, unary *ast.UnaryExpr, env *Env) (aexpr *UnaryExpr, errs []error) {
	aexpr = &UnaryExpr{UnaryExpr: unary}

	var moreErrs []error
	if aexpr.X, moreErrs = CheckExpr(ctx, unary.X, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}

	if errs == nil {
		a := aexpr.X.(Expr)
		t := a.KnownType()

		// Check for multi valued expressions. Not much we can do if we find one
		// TODO check for single values

		// TODO shim
		if len(t) != 1 {
			return aexpr, errs
		}

		if a.IsConst() {
			if c, ok := t[0].(ConstType); ok {
				aexpr.constValue, moreErrs = evalConstUnaryExpr(ctx, aexpr, c)
				if moreErrs != nil {
					errs = append(errs, moreErrs...)
				} else {
					aexpr.knownType = t
				}
			}
		}
	}
	return aexpr, errs
}

// Evaluates a const binary Expr. May return a sensical constValue
// even if ErrTruncatedConst errors are present
func evalConstUnaryExpr(ctx *Ctx, constExpr *UnaryExpr, resultType ConstType) (constValue, []error) {
	x := constExpr.X.(Expr).Const()
	switch resultType.(type) {
	case ConstIntType, ConstRuneType, ConstFloatType, ConstComplexType:
		xx := x.Interface().(*ConstNumber)
		return evalConstUnaryNumericExpr(ctx, constExpr, xx)
	case ConstBoolType:
		xx := x.Bool()
		return evalConstUnaryBoolExpr(ctx, constExpr, xx)
	default:
		return constValue{}, []error{ErrInvalidUnaryOperation{at(ctx, constExpr)}}
	}
}

func evalConstUnaryNumericExpr(ctx *Ctx, constExpr *UnaryExpr, x *ConstNumber) (constValue, []error) {
	switch constExpr.Op {
	case token.ADD:
		return constValueOf(x), nil
	case token.SUB:
		zero := &ConstNumber{Type: x.Type}
		return constValueOf(zero.Sub(zero, x)), nil
	case token.XOR:
		if x.Type.IsIntegral() {
			minusOne := NewConstInt64(-1)
			return constValueOf(minusOne.Xor(minusOne, x)), nil
		}
	}
	return constValue{}, []error{ErrInvalidUnaryOperation{at(ctx, constExpr)}}
}

func evalConstUnaryBoolExpr(ctx *Ctx, constExpr *UnaryExpr, x bool) (constValue, []error) {
	switch constExpr.Op {
	case token.NOT:
		return constValueOf(!x), nil
	default:
		return constValue{}, []error{ErrInvalidUnaryOperation{at(ctx, constExpr)}}
	}
}
