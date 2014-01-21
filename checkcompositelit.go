package eval

import (
	"reflect"
	"go/ast"
)

/*
func checkCompositeLit(ctx *Ctx, lit *ast.CompositeLit, env *Env) (aexpr *CompositeLit, errs []error) {
	aexpr = &CompositeLit{CompositeLit: lit}

	var moreErrs []error
        // TODO confirm the type actually exists
        aexpr.Type = lit.Type
	t, _ := evalType(ctx, aexpr.Type, env)

	if t.Kind() == reflect.Struct {
		for i := range lit.Elts {
			// Don't typecheck struct keys.
			if kv, ok := aexpr.Elts[i].(*ast.KeyValueExpr); ok {
				if kv.Value, moreErrs = CheckExpr(ctx, kv.Value, env); moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
				aexpr.Elts[i] = &KeyValueExpr{KeyValueExpr: kv}
			} else {
				if aexpr.Elts[i], moreErrs = CheckExpr(ctx, lit.Elts[i], env); moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
			}
		}
	} else {
		for i := range lit.Elts {
			if aexpr.Elts[i], moreErrs = CheckExpr(ctx, lit.Elts[i], env); moreErrs != nil {
				errs = append(errs, moreErrs...)
			}
		}
	}
	return aexpr, errs
}

*/

func checkCompositeLit(ctx *Ctx, lit *ast.CompositeLit, env *Env) (*CompositeLit, []error) {
	return checkCompositeLitR(ctx, lit, nil, env)
}

// Recursively check composite literals, where a child composite lit's type depends the
// parent's type For example, the expression [][]int{{1,2},{3,4}} contains two
// slice lits, {1,2} and {3,4}, but their types are inferenced from the parent [][]int{}.
func checkCompositeLitR(ctx *Ctx, lit *ast.CompositeLit, t reflect.Type, env *Env) (*CompositeLit, []error) {
	alit := &CompositeLit{CompositeLit: lit}

	// We won't generate any errors here if the given type does not match lit.Type.
	// The caller will need to detect the type incompatibility.
	if lit.Type != nil {
		var err error
		t, err = evalType(ctx, lit.Type, env)
		if err != nil {
			return alit, []error{err}
		}
	}

	alit.knownType = knownType{t}

	switch t.Kind() {
	//case reflect.Map:
	case reflect.Array, reflect.Slice:
		return checkCompositeLitArrayOrSlice(ctx, alit, t, env)
	case reflect.Struct:
		return checkCompositeLitStruct(ctx, alit, t, env)
	default:
		panic("eval: unimplemented composite lit " + t.Kind().String())
	}
}

func checkCompositeLitArrayOrSlice(ctx *Ctx, lit *CompositeLit, t reflect.Type, env *Env) (*CompositeLit, []error) {
	var errs, moreErrs []error
	eltT := t.Elem()
	maxIndex, curIndex := 0, 0
	outOfBounds := false
	length := -1
	if t.Kind() == reflect.Array {
		length = t.Len()
	}
	used := make(map[int] bool, len(lit.Elts))
	// Check all keys are valid and calculate array or slice length.
	// Elements with key are placed at the keyed position.
	// Elements without are placed in the next position.
	// For example, []int{1, 2:1, 1} -> [1, 0, 1, 1]
	for i := range lit.Elts {
		var value *ast.Expr
		kv, ok := lit.Elts[i].(*ast.KeyValueExpr)
		if !ok {
			value = &lit.Elts[i]
		} else {
			lit.Elts[i] = &KeyValueExpr{KeyValueExpr: kv}
			value = &kv.Value
			// Check the array key
			var index int
			kv.Key, index, ok, moreErrs = checkArrayIndex(ctx, kv.Key, env);
			if !ok || moreErrs != nil {
				// NOTE[crc] Haven't checked the gc implementation, but
				// from experimentation it seems that only undefined
				// idents are reported. This filter should perhaps be part
				// of checkArrayIndex
				for _, err := range moreErrs {
					if _, ok := err.(ErrUndefined); ok {
						errs = append(errs, err)
					}
				}
				errs = append(errs, ErrBadArrayKey{at(ctx, kv.Key)})
				// Don't include this element in index calculations
				curIndex -= 1
				goto check
			}
			lit.indices = append(lit.indices, struct{pos, index int}{i, index})
			curIndex = index
		}
		if maxIndex < curIndex {
			maxIndex = curIndex
		}
		if !outOfBounds && length != -1 && curIndex >= length {
			outOfBounds = true
			errs = append(errs, ErrArrayKeyOutOfBounds{at(ctx, lit.Elts[i]), t, curIndex})
		}
		// has this index been used already
		if used[curIndex] {
			errs = append(errs, ErrDuplicateArrayKey{at(ctx, kv.Key), curIndex})
		}
		used[curIndex] = true

check:
		// finally check the value
		*value, moreErrs = checkArrayValue(ctx, *value, eltT, env)
		if moreErrs != nil {
			errs = append(errs, moreErrs...)
		}

		curIndex += 1
	}
	lit.indices = append(lit.indices, struct{pos, index int}{-1, -1})
	lit.length = maxIndex
	return lit, errs
}

func checkCompositeLitStruct(ctx *Ctx, lit *CompositeLit, t reflect.Type, env *Env) (*CompositeLit, []error) {
	var errs, moreErrs []error

	// X{} is treated as if it has zero KeyValue'd elements, i.e. unspecified
	// elements are set to zero. This is always valid
	if len(lit.Elts) == 0 {
		return lit, nil
	}

	// gc first checks if there are ANY keys present, and then decides how
	// to process the initialisers.
	keysPresent := false
	for _, elt := range lit.Elts {
		_, ok := elt.(*ast.KeyValueExpr)
		keysPresent = keysPresent || ok
	}

	if keysPresent {
		seen := make(map[string] bool, len(lit.Elts))
		mixed := false
		for i := 0; i < len(lit.Elts); i += 1 {
			kv, ok := lit.Elts[i].(*ast.KeyValueExpr)
			if !ok {
				if !mixed {
					// This error only gets reported once
					mixed = true
					errs = append(errs, ErrMixedStructValues{at(ctx, lit.Elts[i])})
				}
				continue
			}

			lit.Elts[i] = &KeyValueExpr{KeyValueExpr: kv}
			// Check the key is a struct member
			if ident, ok := kv.Key.(*ast.Ident); !ok {
				// This check is a hack for making kv.Key printable.
				// field identifiers should not usually be type checked.
				kv.Key = fakeCheckExpr(kv.Key, env)
				errs = append(errs, ErrInvalidStructField{at(ctx, kv.Key)})
			} else if name := ident.Name; false {
			} else if field, ok := t.FieldByName(name); !ok {
				errs = append(errs, ErrUnknownStructField{at(ctx, kv.Key), t, name})
			} else {
				if seen[name] {
					errs = append(errs, ErrDuplicateStructField{at(ctx, kv.Key), name})
				}
				seen[name] = true
				lit.fields = append(lit.fields, field.Index[0])
				kv.Value, moreErrs = checkStructField(ctx, kv.Value, field, env)
				if moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
			}
		}
	} else {
		numFields := t.NumField()
		var i int
		for i = 0; i < numFields && i < len(lit.Elts); i += 1 {
			field := t.Field(i)
			lit.Elts[i], moreErrs = checkStructField(ctx, lit.Elts[i], field, env)
			if moreErrs != nil {
				errs = append(errs, moreErrs...)
			}
		}
		if numFields != len(lit.Elts) {
			errs = append(errs, ErrWrongNumberOfStructValues{at(ctx, lit)})
		}
		// Remaining fields are type checked reguardless of use
		for ; i < len(lit.Elts); i += 1 {
			lit.Elts[i], moreErrs = CheckExpr(ctx, lit.Elts[i], env)
			if moreErrs != nil {
				errs = append(errs, moreErrs...)
			}
		}
	}
	return lit, errs
}

func checkArrayValue(ctx *Ctx, expr ast.Expr, eltT reflect.Type, env *Env) (Expr, []error) {
	aexpr, conversionFailed, errs := checkExprAssignableTo(ctx, expr, eltT, env)
	if conversionFailed {
		// NOTE[crc] this hack removes conversion errors from consts other
		// than strings and nil to match the output of gc.
		if ccerr, ok := errs[0].(ErrBadConstConversion); ok {
			if ccerr.from == ConstNil {
				// No ErrBadArrayValue for nil
				return aexpr, errs
			} else if ccerr.from != ConstString {
				// gc implementation only displays string conversion errors
				errs = nil
			}
		}
		errs = append(errs, ErrBadArrayValue{at(ctx, aexpr), eltT})
	}
	return aexpr, errs
}

func checkStructField(ctx *Ctx, expr ast.Expr, field reflect.StructField, env *Env) (Expr, []error) {
	aexpr, conversionFailed, errs := checkExprAssignableTo(ctx, expr, field.Type, env)
	if conversionFailed {
		errs = append([]error{}, ErrBadStructValue{at(ctx, aexpr), field.Type})
	}
	return aexpr, errs
}
