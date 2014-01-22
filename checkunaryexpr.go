package eval

import (
	"reflect"

	"go/ast"
	"go/token"
)

func checkUnaryExpr(ctx *Ctx, unary *ast.UnaryExpr, env *Env) (*UnaryExpr, []error) {
	aexpr := &UnaryExpr{UnaryExpr: unary}

	x, errs := CheckExpr(ctx, unary.X, env)
	if errs == nil || x.IsConst() {
		if t, err := expectSingleType(ctx, x.KnownType(), x); err != nil {
			errs = append(errs, err)
		} else if unary.Op == token.AND { // address off
			if !isAddressable(x) {
				printableX := fakeCheckExpr(unary.X, env)
				printableX.setKnownType(knownType{t})
				errs = append(errs, ErrInvalidAddressOf{at(ctx, printableX)})
			}
			t := x.KnownType()[0]
			if ct, ok := t.(ConstType); ok {
				if ct == ConstNil {
					errs = append(errs, ErrUntypedNil{at(ctx, x)})
				} else {
					ptrT := reflect.PtrTo(unhackType(ct.DefaultPromotion()))
					aexpr.knownType = knownType{ptrT}
				}
			} else {
				ptrT := reflect.PtrTo(unhackType(t))
				aexpr.knownType = knownType{ptrT}
			}
			aexpr.X = x
		// TODO handle <-
		} else {
			aexpr.X = x
			// All numeric and bool unary expressions do not change type
			aexpr.knownType = knownType(x.KnownType())
			if x.IsConst() {
				if ct, ok := t.(ConstType); ok {
					c, moreErrs := evalConstUnaryExpr(ctx, aexpr, x, ct)
					if moreErrs != nil {
						errs = append(errs, moreErrs...)
					}
					aexpr.constValue = c

				} else {
					c, moreErrs := evalConstTypedUnaryExpr(ctx, aexpr, x)
					if moreErrs != nil {
						errs = append(errs, moreErrs...)
					}
					aexpr.constValue = c
				}
			} else {
				if !isUnaryOpDefinedOn(unary.Op, t) {
					errs = append(errs, ErrInvalidUnaryOperation{at(ctx, unary)})
				}
			}
		}
        }
	aexpr.X = x
	return aexpr, errs
}

func evalConstUnaryExpr(ctx *Ctx, unary *UnaryExpr, x Expr, xT ConstType) (constValue, []error) {
	switch xT.(type) {
	case ConstIntType, ConstRuneType, ConstFloatType, ConstComplexType:
		xx := x.Const().Interface().(*ConstNumber)
		return evalConstUnaryNumericExpr(ctx, unary, xx)
	case ConstBoolType:
		xx := x.Const().Bool()
		return evalConstUnaryBoolExpr(ctx, unary, xx)
	default:
		return constValue{}, []error{ErrInvalidUnaryOperation{at(ctx, unary)}}
	}
}

func evalConstTypedUnaryExpr(ctx *Ctx, unary *UnaryExpr, x Expr) (constValue, []error) {
	t := x.KnownType()[0]
	if xx, ok := convertTypedToConstNumber(x.Const()); ok {
		zz, errs := evalConstUnaryNumericExpr(ctx, unary, xx)
		if !reflect.Value(zz).IsValid() {
			return constValue{}, errs
		}
		rr, moreErrs := promoteConstToTyped(ctx, xx.Type, zz, t, x)
		return rr, append(errs, moreErrs...)
	} else if t.Kind() == reflect.Bool {
		return evalConstUnaryBoolExpr(ctx, unary, x.Const().Bool())
	}
	return constValue{}, []error{ErrInvalidUnaryOperation{at(ctx, unary)}}
}

func evalConstUnaryNumericExpr(ctx *Ctx, unary *UnaryExpr, x *ConstNumber) (constValue, []error) {
	switch unary.Op {
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
	return constValue{}, []error{ErrInvalidUnaryOperation{at(ctx, unary)}}
}

func evalConstUnaryBoolExpr(ctx *Ctx, unary *UnaryExpr, x bool) (constValue, []error) {
	switch unary.Op {
	case token.NOT:
		return constValueOf(!x), nil
	default:
		return constValue{}, []error{ErrInvalidUnaryOperation{at(ctx, unary)}}
	}
}
