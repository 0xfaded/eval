package eval

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"go/ast"
	"go/token"
)

// Equivalent of reflect.New, but unwraps internal Types into their original reflect.Type
func hackedNew(t reflect.Type) reflect.Value {
	return reflect.New(unhackType(t))
}

// Get the underlying reflect.Type a hacked type
func unhackType(t reflect.Type) reflect.Type {
	switch tt := t.(type) {
	case Rune:
		return tt.Type
	case Byte:
		return tt.Type
	default:
		return t
	}
}

// Determine if type from is assignable to type to. From and To must not be ConstTypes
func typeAssignableTo(from, to reflect.Type) bool {
	return from.AssignableTo(unhackType(to))
}

// Check expr and then determine if it is assignable to type T. Essentially
// exprAssignableTo(CheckExpr(expr), t), but errors are accumulated and a
// bool value is returned indicating if an additional error is needed.
// The error is needed if expr successfully typechecks but is not
// assignable to t. In this case, the caller should append an appropriate,
// use case specific error.
func checkExprAssignableTo(ctx *Ctx, expr ast.Expr, t reflect.Type, env *Env) (Expr, bool, []error) {
	var errs []error
	aexpr, moreErrs := CheckExpr(ctx, expr, env)
	if moreErrs != nil {
		errs = append(errs, moreErrs...)
	} else if _, err := expectSingleType(ctx, aexpr.KnownType(), aexpr); err != nil {
		errs = append(errs, err)
	}
	if errs != nil {
		return aexpr, false, errs
	}
	ok, convErrs := exprAssignableTo(ctx, aexpr, t)
	if convErrs != nil {
		errs = append(errs, convErrs...)
	}
	return aexpr, !ok, errs
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
		return reflect.Value(cv).IsValid(), errs
	}

	return typeAssignableTo(fromType, to), nil
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

// Is op a boolean operator that produces a const bool type.
// Notably absent are LAND(&&) and LOR(||), which result
// in a value of the same type as their operands.
func isBooleanOp(op token.Token) bool {
	switch op {
	case token.EQL, token.NEQ, token.LEQ, token.GEQ, token.LSS, token.GTR:
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

	// This is slighly misleading. slices, funcs and maps are only
	// comparable if their paired operand is nil
	case reflect.Ptr, reflect.Interface, reflect.Struct,
		reflect.Slice, reflect.Map, reflect.Chan:
		return op == token.EQL || op == token.NEQ
	}
	return false
}

func isUnaryOpDefinedOn(op token.Token, t reflect.Type) bool {
	if _, ok := t.(ConstNilType); ok {
		return false
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch op {
		case token.ADD, token.SUB, token.XOR:
			return true
		}
	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		switch op {
		case token.ADD, token.SUB:
			return true
		}
	case reflect.Bool:
		switch op {
		case token.NOT:
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

// Returns the float type that is half the width of the input complex type
func comprisingFloatType(complexType reflect.Type) reflect.Type {
	if complexType == c128 {
		return f64
	} else {
		return f32
	}
}

// Evals an expression with a known result type. If the node is an
// untyped constant, it is converted to type t. This function assumes
// the input is successfully type checked, and therefore is undefined
// incorrectly typed inputs.
func evalTypedExpr(ctx *Ctx, expr Expr, t knownType, env *Env) (xs []reflect.Value, err error) {
        if expr.IsConst() {
                x := expr.Const()
                if ct, ok := expr.KnownType()[0].(ConstType); ok {
                        cx, _ := promoteConstToTyped(ctx, ct, constValue(x), t[0], expr)
                        xs = []reflect.Value{reflect.Value(cx)}
                } else {
                        xs = []reflect.Value{x}
                }
        } else {
                var xxs *[]reflect.Value
                xxs, _, err = EvalExpr(ctx, expr, env)
                xs = *xxs
        }
        return xs, err
}

// Type check an integral node. Returns the type checked node, the
// integer value if constant, ok if the node was indeed integral,
// and checkErrs which occur during the type check. It is possible
// that checkErrs will be non-nil yet ok is still true. In this case
// the errors are non-fatal, such as integer truncation.
func checkInteger(ctx *Ctx, expr ast.Expr, env *Env) (aexpr Expr, i int, ok bool, checkErrs []error) {
	aexpr, checkErrs = CheckExpr(ctx, expr, env)
	if checkErrs != nil && !aexpr.IsConst() {
		return aexpr, 0, false, checkErrs
	}
	t, err := expectSingleType(ctx, aexpr.KnownType(), aexpr)
	if err != nil {
		return aexpr, 0, false, append(checkErrs, err)
	}

	var ii int64
	if ct, ok := t.(ConstType); ok {
		c, moreErrs := promoteConstToTyped(ctx, ct, constValue(aexpr.Const()), intType, aexpr)
		if moreErrs != nil {
			checkErrs = append(checkErrs, moreErrs...)
		}
		v := reflect.Value(c)
		if v.IsValid() {
			ii = v.Int()
		} else {
			return aexpr, 0, false, checkErrs
		}
	} else {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if aexpr.IsConst() {
				ii = aexpr.Const().Int()
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if aexpr.IsConst() {
				ii = int64(aexpr.Const().Uint())
			}
		default:
			return aexpr, 0, false, checkErrs
		}
	}
	return aexpr, int(ii), true, checkErrs
}

// Eval a node and cast it to an int. expr must be a *ConstNumber or integral type
func evalInteger(ctx *Ctx, expr Expr, env *Env) (int, error) {
        if expr.IsConst() {
                x := expr.Const()
                if ct, ok := expr.KnownType()[0].(ConstType); ok {
                        cx, _ := promoteConstToTyped(ctx, ct, constValue(x), intType, expr)
			return int(reflect.Value(cx).Int()), nil
                } else {
			panic(dytc("const bool or string evaluated as int"))
                }
        } else {
                xs, _, err := EvalExpr(ctx, expr, env);
		if err != nil {
			return 0, err
		}
		x := (*xs)[0]
		switch x.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int(x.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int(x.Uint()), nil
		default:
			panic(dytc("non-integral type evaluated as int"))
		}
        }
}

func checkArrayIndex(ctx *Ctx, expr ast.Expr, env *Env) (aexpr Expr, i int, ok bool, checkErrs []error) {
	aexpr, checkErrs = CheckExpr(ctx, expr, env)
	if !aexpr.IsConst() {
		return aexpr, 0, false, checkErrs
	}
	t := aexpr.KnownType()[0]
	var ii int64
	if ct, ok := t.(ConstType); ok {
		c, moreErrs := promoteConstToTyped(ctx, ct, constValue(aexpr.Const()), intType, aexpr)
		if moreErrs != nil {
			checkErrs = append(checkErrs, moreErrs...)
		}
		v := reflect.Value(c)
		if v.IsValid() {
			ii = v.Int()
		} else {
			return aexpr, 0, false, checkErrs
		}
	} else {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ii = aexpr.Const().Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ii = int64(aexpr.Const().Uint())
		default:
			return aexpr, 0, false, checkErrs
		}
	}
	// The limit of 2^31-1 is derived from the gc implementation,
	// which seems to use this definition whilst type checking.
	// The actual definition is the "largest value representable by an int"
	return aexpr, int(ii), 0 <= ii && ii <= 0x7fffffff, checkErrs
}

func isValidGoName(name string) bool {
	if name == "" {
		return false
	}

	ok := true
	// TODO[crc] go 1.2 includes unicode.DecodeRuneInString for decoding the first rune
	for _, r := range name {
		ok = ok && (r == '_' || unicode.IsLetter(r))
		break
	}
	for _, r := range name {
		ok = ok && (r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r))
	}
	return ok
}

// spec: addressable, that is, either a
// variable,
// pointer indirection,
// or slice indexing operation;
// or a field selector of an addressable struct operand;
// or an array indexing operation of an addressable array.
// As an exception to the addressability requirement, x may also be a (possibly parenthesized) composite literal
func isAddressable(expr Expr) bool {
	expr = skipSuperfluousParens(expr)
	switch n := expr.(type) {
	case *Ident:
		return n.source == envVar
	case *StarExpr:
		return true
	case *IndexExpr:
		x := n.X.(Expr)
		t := x.KnownType()[0]
		switch t.Kind() {
		case reflect.Slice:
			return true
		case reflect.Array:
			return isAddressable(x)
		case reflect.Ptr:
			return true
		}
	case *SelectorExpr:
		x := n.X.(Expr)
		t := x.KnownType()[0]
		switch t.Kind() {
		case reflect.Struct:
			return isAddressable(x)
		case reflect.Ptr:
			return true
		}
	}
	return false
}

func isAddressableOrCompositeLit(expr Expr) bool {
	expr = skipSuperfluousParens(expr)
	if _, ok := expr.(*CompositeLit); ok {
		return true
	} else {
		return isAddressable(expr)
	}
}

func isStructComparable(structT reflect.Type) bool {
	numField := structT.NumField()
	for i := 0; i < numField; i += 1 {
		field := structT.Field(i)
		if field.Name == "_" {
			continue
		}
		k := field.Type.Kind()
		if k == reflect.Slice || k == reflect.Map || k == reflect.Chan {
			return false
		}
	}
	return true
}
