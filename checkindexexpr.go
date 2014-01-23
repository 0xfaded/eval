package eval

import (
	"reflect"
	"go/ast"
)

func checkIndexExpr(ctx *Ctx, index *ast.IndexExpr, env *Env) (*IndexExpr, []error) {
	aexpr := &IndexExpr{IndexExpr: index}
	x, errs := CheckExpr(ctx, index.X, env)
	aexpr.X = x
	if errs != nil && !x.IsConst() {
		return aexpr, errs
	}

	t, err := expectSingleType(ctx, x.KnownType(), x)
	if err != nil {
		return aexpr, append(errs, err)
	}

	// index of array pointer is short hand for dereference and then index
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Array {
		t = t.Elem()
	}

	switch t.Kind() {
	// case reflect.Map:
	case reflect.String:
		aexpr.knownType = knownType{u8}
		i, moreErrs := checkIndexVectorExpr(ctx, x, index.Index, env)
		if moreErrs != nil {
			errs = append(errs, moreErrs...)
		}
		aexpr.Index = i
		return aexpr, errs
	case reflect.Array, reflect.Slice:
		aexpr.knownType = knownType{t.Elem()}
		i, moreErrs := checkIndexVectorExpr(ctx, x, index.Index, env)
		if moreErrs != nil {
			errs = append(errs, moreErrs...)
		}
		aexpr.Index = i
		return aexpr, errs
	default:
		return aexpr, append(errs, ErrInvalidIndexOperation{at(ctx, index)})
	}
}

func checkIndexVectorExpr(ctx *Ctx, x Expr, index ast.Expr, env *Env) (Expr, []error) {
	t := x.KnownType()[0]
	i, iint, ok, errs := checkInteger(ctx, index, env)
	if errs != nil && !i.IsConst() {
		// Type check of index failed
	} else if !ok {
		// Type check of index passed but this node is not an integer
		printableIndex := fakeCheckExpr(index, env)
		printableIndex.setKnownType(i.KnownType())
		errs = append(errs, ErrNonIntegerIndex{at(ctx, printableIndex)})
	} else if i.IsConst() {
		// If we know the index at compile time, we must assert it is in bounds.
		if iint < 0 {
			errs = append(errs, ErrIndexOutOfBounds{at(ctx, index), iint})
		} else if t.Kind() == reflect.Array {
			if iint >= t.Len() {
				errs = append(errs, ErrIndexOutOfBounds{at(ctx, index), iint})
			}
		} else if t.Kind() == reflect.String && x.IsConst() {
			str := x.Const()
			if iint >= str.Len() {
				errs = append(errs, ErrIndexOutOfBounds{at(ctx, index), iint})
			}
		}
	}
	return i, errs
}
