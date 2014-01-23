package eval

import (
	"reflect"
	"go/ast"
)

// TODO[crc] support [::] syntax after go1.2 upgrade
func checkSliceExpr(ctx *Ctx, slice *ast.SliceExpr, env *Env) (*SliceExpr, []error) {
	aexpr := &SliceExpr{SliceExpr: slice}
	x, errs := CheckExpr(ctx, slice.X, env)
	aexpr.X = x
	if errs != nil && !x.IsConst() {
		return aexpr, errs
	}

	t, err := expectSingleType(ctx, x.KnownType(), x)
	if err != nil {
		return aexpr, append(errs, err)
	}

	// arrays must be addressable
	if t.Kind() == reflect.Array && !isAddressable(x) {
		return aexpr, append(errs, ErrUnaddressableSliceOperand{at(ctx, aexpr)})
	}
	// slice of array pointer is short hand for dereference and then slice
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Array {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		if t == ConstString {
			// spec: ConstString[:] fields string
			aexpr.knownType = knownType{stringType}
		} else {
			aexpr.knownType = knownType(x.KnownType())
		}
		if slice.Low != nil {
			low, moreErrs := checkIndexVectorExpr(ctx, x, slice.Low, env)
			aexpr.Low = low
			if moreErrs != nil {
				errs = append(errs, moreErrs...)
				if !low.IsConst() {
					return aexpr, errs
				}
			}
		}
		if slice.High != nil {
			high, moreErrs := checkIndexVectorExpr(ctx, x, slice.Low, env)
			aexpr.High = high
			if moreErrs != nil {
				errs = append(errs, moreErrs...)
				if !high.IsConst() {
					return aexpr, errs
				}
			}
		}
		return aexpr, errs
	default:
		return aexpr, append(errs, ErrInvalidSliceOperation{at(ctx, index)})
	}
}

