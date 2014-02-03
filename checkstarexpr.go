package eval

import (
	"reflect"

	"fmt"
	"go/ast"
)

func checkStarExpr(ctx *Ctx, star *ast.StarExpr, env *Env) (*StarExpr, []error) {
	aexpr := &StarExpr{StarExpr: star}
	fmt.Printf("%v %T\n", aexpr.X, aexpr.X)
	x, errs := CheckExpr(ctx, aexpr.X, env)

	if errs != nil && !x.IsConst() {
		return aexpr, errs
	} else if t, err := expectSingleType(ctx, x.KnownType(), x); err != nil {
		errs = append(errs, err)
	} else if t == ConstNil {
		errs = append(errs, ErrInvalidIndirect{at(ctx, x)})
	} else if t.Kind() != reflect.Ptr {
		printableX := fakeCheckExpr(star.X, env)
		printableX.setKnownType(x.KnownType())
		errs = append(errs, ErrInvalidIndirect{at(ctx, printableX)})
	} else {
		aexpr.knownType = knownType{t.Elem()}
	}
	aexpr.X = x
	return aexpr, errs
}
