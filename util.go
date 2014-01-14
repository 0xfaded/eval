package eval

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"go/ast"
	"go/token"
)

// Equivalent of reflect.New, but unwraps internal Types into their original reflect.Type
func hackedNew(t reflect.Type) reflect.Value {
	switch tt := t.(type) {
	case Rune:
		return reflect.New(tt.Type)
	default:
		return reflect.New(t)
	}
}

// Get the underlying reflect.Type a hacked type
func unhackType(t reflect.Type) reflect.Type {
	switch tt := t.(type) {
	case Rune:
		return tt.Type
	default:
		return t
	}
}


// Determine if type from is assignable to type to
func typeAssignableTo(from, to reflect.Type) bool {
	// Handle the rune alias
	if r, ok := from.(Rune); ok {
		from = r.Type
	}

	return from.AssignableTo(to)
}

// Determine if the result of from expr is assignable to type to. to must be a vanilla reflect.Type.
// from must have a KnownType() of length 1. Const types that raise overflow and truncation
// errors will still return true, but the errors will be reflected in the []error slice.
func exprAssignableTo(ctx *Ctx, from Expr, to reflect.Type) (bool, []error) {
	if len(from.KnownType()) != 1 {
		panic("go-eval: assignableTo called with from.KnownType() != 1")
	}
	fromType := from.KnownType()[0]

	// Check that consts can be converted
	if c, ok := fromType.(ConstType); ok && from.IsConst() {
		// If cv is a valid value, then the types are assignable even if
		// other conversion errors, such as overflows, are present.
		cv, errs := promoteConstToTyped(ctx, c, constValue(from.Const()), to, from)
		if reflect.Value(cv).IsValid() {
			return true, errs
		} else {
			// If the conversion was invalid, the caller usually does not
			// want the ErrBadConstConversion error. Instead it will produce
			// its own error specific to the cause
			return false, nil
		}
	}

	return typeAssignableTo(fromType, to), nil
}

func assignableValue(x reflect.Value, to reflect.Type, xTyped bool) (reflect.Value, error) {
	var err error
	if xTyped {
		if x.Type().AssignableTo(to) {
			return x, nil
		}
	} else {
		if x, err = promoteUntypedNumeral(x, to); err == nil {
			return x, nil
		}
	}
	return x, errors.New(fmt.Sprintf("Cannot convert %v to type %v", x, to))
}

func setTypedValue(dst, src reflect.Value, srcTyped bool) error {
	if assignable, err := assignableValue(src, dst.Type(), srcTyped); err != nil {
		return errors.New(fmt.Sprintf("Cannot assign %v = %v", dst, src))
	} else {
		dst.Set(assignable)
		return nil
	}
}

func makeSliceWithValues(elts []reflect.Value, sliceType reflect.Type) (reflect.Value, error) {
	slice := reflect.MakeSlice(sliceType, len(elts), len(elts))
	for i := 0; i < slice.Len(); i += 1 {
		if err := setTypedValue(slice.Index(i), elts[i], true); err != nil {
			return reflect.Value{}, nil
		}
	}
	return slice, nil
}


// Only considers untyped kinds produced by our runtime. Assumes input type is unnamed
func isUntypedNumeral(x reflect.Value) bool {
	switch x.Kind() {
	case reflect.Int64, reflect.Float64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func promoteUntypedNumeral(untyped reflect.Value, to reflect.Type) (reflect.Value, error) {
	// The only valid promotion that cannot be directly converted is int|float -> complex
	if untyped.Type().ConvertibleTo(to) {
		return untyped.Convert(to), nil
	} else if to.Kind() == reflect.Complex64 || to.Kind() == reflect.Complex128 {
		floatType := reflect.TypeOf(float64(0))
		if untyped.Type().ConvertibleTo(floatType) {
			return reflect.ValueOf(complex(untyped.Convert(floatType).Float(), 0)), nil
		}
	}
	return reflect.Value{}, errors.New(fmt.Sprintf("cannot convert %v to %v", untyped, to))
}

// Only considers untyped kinds produced by our runtime. Assumes input type is unnamed
func promoteUntypedNumerals(x, y reflect.Value) (reflect.Value, reflect.Value) {
	switch x.Kind() {
	case reflect.Int64:
		switch y.Kind() {
		case reflect.Int64:
			return x, y
		case reflect.Float64:
			return x.Convert(y.Type()), y
		case reflect.Complex128:
			return reflect.ValueOf(complex(float64(x.Int()), 0)), y
		}
	case reflect.Float64:
		switch y.Kind() {
		case reflect.Int64:
			return x, y.Convert(x.Type())
		case reflect.Float64:
			return x, y
		case reflect.Complex128:
			return reflect.ValueOf(complex(x.Float(), 0)), y
		}
	case reflect.Complex128:
		switch y.Kind() {
		case reflect.Int64:
			return x, reflect.ValueOf(complex(float64(y.Int()), 0))
		case reflect.Float64:
			return x, reflect.ValueOf(complex(y.Float(), 0))
		case reflect.Complex128:
			return x, y
		}
	}
	panic(fmt.Sprintf("runtime: bad untyped numeras %v and %v", x, y))
}

// TODO remove this when type checker is complete
func expectSingleValue(ctx *Ctx, values []reflect.Value, srcExpr ast.Expr) (reflect.Value, error) {
	if len(values) == 0 {
		return reflect.Value{}, ErrMissingValue{at(ctx, srcExpr)}
	} else if len(values) != 1 {
		return reflect.Value{}, ErrMultiInSingleContext{at(ctx, srcExpr)}
	} else {
		return values[0], nil
	}
}

func expectSingleType(ctx *Ctx, types []reflect.Type, srcExpr ast.Expr) (reflect.Type, error) {
	if len(types) == 0 {
		return nil, ErrMissingValue{at(ctx, srcExpr)}
	} else if len(types) != 1 {
		return nil, ErrMultiInSingleContext{at(ctx, srcExpr)}
	} else {
		return types[0], nil
	}
}

func isBooleanOp(op token.Token) bool {
	switch op {
	case token.EQL, token.NEQ, token.LEQ, token.GEQ, token.LSS, token.GTR, token.LAND, token.LOR:
		return true
	default:
		return false
	}
}

func isOpDefinedOn(op token.Token, t reflect.Type) bool {
	if _, ok := t.(ConstNilType); ok {
		return false
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch op {
		case token.ADD, token.SUB, token.MUL, token.QUO,
			token.REM, token.AND, token.OR, token.XOR, token.AND_NOT,
			token.EQL, token.NEQ,
			token.LEQ, token.GEQ, token.LSS, token.GTR:
			return true
		}

	case reflect.Float32, reflect.Float64:
		switch op {
		case token.ADD, token.SUB, token.MUL, token.QUO,
			token.EQL, token.NEQ,
			token.LEQ, token.GEQ, token.LSS, token.GTR:
			return true
		}

	case reflect.Complex64, reflect.Complex128:
		switch op {
		case token.ADD, token.SUB, token.MUL, token.QUO,
			token.EQL, token.NEQ:
			return true
		}

	case reflect.Bool:
		switch op {
		case token.LAND, token.LOR, token.EQL, token.NEQ:
			return true
		}

	case reflect.String:
		switch op {
		case token.ADD, token.EQL, token.NEQ, token.LEQ, token.GEQ, token.LSS, token.GTR:
			return true
		}
	}
	return false
}

// FIXME: should also match and handle just a line and no column
var parseError = regexp.MustCompile(`^([0-9]+):([0-9]+): `)

// FormatErrorPos formats source to show the position that a (parse)
// error occurs. When this works, it returns a slice of one or two
// strings: the source line with the error and if it can find a column
// position under that, a line indicating the position where the error
// occurred.
//
// For example, if we have:
//		source := `split(os.Args ", )")`
//		errmsg := "1:15: expected ')'"
// then PrintErrPos(source, errmsg) returns:
//  {
//		`split(os.Args ", )")`,
//		`-------------^`
//  }
//
// If something is wrong parsing the error message or matching it with
// the source, an empty slice is returned.
func FormatErrorPos(source, errmsg string) (cursored [] string) {
	matches := parseError.FindStringSubmatch(errmsg)
	if len(matches) == 3 {
		var err error
		var line, column int
		if line, err = strconv.Atoi(matches[1]); err != nil {
			return cursored
		}
		if column, err = strconv.Atoi(matches[2]); err != nil {
			return cursored
		}
		sourceLines := strings.Split(source, "\n")
		if line > len(sourceLines) {
			return cursored
		}
		errLine := sourceLines[line-1]
		cursored = append(cursored, errLine)
		if column-1 > len(errLine) || column < 1 {
			return cursored
		} else if column == 1 {
			cursored = append(cursored, "^")
		} else {
			cursored = append(cursored, strings.Repeat("-", column-1) + "^")
		}
	}
	return cursored
}

// Walk the ast of expressions like (((x))) and return the inner *ParenExpr.
// Returns input Expr if it is not a *ParenExpr
func skipSuperfluousParens(expr Expr) Expr {
	if p, ok := expr.(*ParenExpr); ok {
		// Remove useless parens from (((x))) expressions
		var tmp *ParenExpr
		for ; ok; tmp, ok = p.X.(*ParenExpr) {
			p = tmp
		}

		// Remove parens from all expressions where order of evaluation is irrelevant
		switch p.X.(type) {
		case *BinaryExpr:
			return p
		default:
			return p.X.(Expr)
		}
	}
	return expr
}

