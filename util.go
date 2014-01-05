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
		if column > len(errLine) || column < 1 {
			return cursored
		} else if column == 1 {
			cursored = append(cursored, "^")
		} else {
			cursored = append(cursored, strings.Repeat("-", column-2) + "^")
		}
	}
	return cursored
}
