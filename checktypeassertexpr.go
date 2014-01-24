package eval

import (
	"reflect"

	"go/ast"
)

func checkTypeAssertExpr(ctx *Ctx, assert *ast.TypeAssertExpr, env *Env) (*TypeAssertExpr, []error) {
	aexpr := &TypeAssertExpr{TypeAssertExpr: assert}
	x, errs := CheckExpr(ctx, aexpr.X, env)
	aexpr.X = x

	if errs != nil && !x.IsConst() {
		return aexpr, errs
	} else if xT, err := expectSingleType(ctx, x.KnownType(), x); err != nil {
		errs = append(errs, err)
	} else if xT == ConstNil {
		errs = append(errs, ErrUntypedNil{at(ctx, x)})
	} else if xT.Kind() != reflect.Interface {
		errs = append(errs, ErrInvalidTypeAssert{at(ctx, aexpr)})
	} else {
		typ, t, moreErrs := checkType(ctx, assert.Type, env)
		aexpr.Type = typ
		if moreErrs != nil {
			errs = append(errs, moreErrs...)
		} else {
			aexpr.knownType = knownType{t}
			if t.Kind() != reflect.Interface && !unhackType(t).Implements(xT) {
				errs = append(errs, ErrImpossibleTypeAssert{at(ctx, aexpr)})
			}
		}
	}
	return aexpr, errs
}
