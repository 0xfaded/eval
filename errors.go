package eval

import (
	"fmt"
	"reflect"

	"go/token"
)

type ErrBadBasicLit struct {
	*BasicLit
}

type ErrUndefined struct {
	Expr
}

type ErrInvalidIndirect struct {
	Expr
}

type ErrUndefinedFieldOrMethod struct {
	Expr
}

type ErrCallNonFuncType struct {
	Expr
}

type ErrWrongNumberOfArgs struct {
	*CallExpr
	numArgs int
}

type ErrWrongArgType struct {
	Expr
	call *CallExpr
	argPos int
}

type ErrInvalidEllipsisInCall struct {
	*CallExpr
}

type ErrMissingValue struct {
	Expr
}

type ErrMultiInSingleContext struct {
	Expr
}

type ErrBadMapIndex struct {
	Expr
	keyT reflect.Type
}

type ErrNonIntegerIndex struct {
	Expr
}

type ErrIndexOutOfBounds struct {
	Expr
	x Expr
	i int
}

type ErrInvalidIndexOperation struct {
	*IndexExpr
}

type ErrInvalidSliceIndex struct {
	*SliceExpr
}

type ErrInvalidSliceOperation struct {
	*SliceExpr
}

type ErrUnaddressableSliceOperand struct {
	*SliceExpr
}

type ErrInvalidIndex struct {
	Expr
	indexValue reflect.Value
	containerType reflect.Type
}

type ErrDivideByZero struct {
	*BinaryExpr
}

type ErrInvalidBinaryOperation struct {
	*BinaryExpr
}

type ErrInvalidUnaryOperation struct {
	*UnaryExpr
}

type ErrInvalidAddressOf struct {
	Expr
}

type ErrInvalidRecvFrom struct {
	Expr
}

type ErrBadConversion struct {
	Expr
	from reflect.Type
	to reflect.Type
}

type ErrBadConstConversion struct {
	Expr
	from reflect.Type
	to reflect.Type
}

type ErrTruncatedConstant struct {
	Expr
	to ConstType
	constant *ConstNumber
}

type ErrOverflowedConstant struct {
	Expr
	from ConstType
	to reflect.Type
	constant *ConstNumber
}

type ErrUntypedNil struct {
	Expr
}

type ErrTypeUsedAsExpression struct {
	Expr
}

type ErrUncomparableMapKey struct {
	Expr
	keyT reflect.Type
}

type ErrMissingMapKey struct {
	Expr
}

type ErrBadMapKey struct {
	Expr
	keyT reflect.Type
}

type ErrDuplicateMapKey struct {
	Expr
}

type ErrBadMapValue struct {
	Expr
	eltT reflect.Type
}

type ErrBadArrayKey struct {
	Expr
}

type ErrArrayKeyOutOfBounds struct {
	Expr
	arrayT reflect.Type
	index int
}

type ErrDuplicateArrayKey struct {
	Expr
	index int
}

type ErrBadArrayValue struct {
	Expr
	eltT reflect.Type
}

type ErrUnknownStructField struct {
	Expr
	structT reflect.Type
	field string
}

type ErrInvalidStructField struct {
	Expr
}

type ErrDuplicateStructField struct {
	*Ident
	field string
}

type ErrMixedStructValues struct {
	Expr
}

type ErrWrongNumberOfStructValues struct {
	*CompositeLit
}

type ErrMissingCompositeLitType struct {
	*CompositeLit
}

type ErrBadStructValue struct {
	Expr
	eltT reflect.Type
}

type ErrInvalidTypeAssert struct {
	*TypeAssertExpr
}

type ErrImpossibleTypeAssert struct {
	*TypeAssertExpr
}

type ErrBuiltinWrongNumberOfArgs struct {
	*CallExpr
}

type ErrBuiltinWrongArgType struct {
	Expr
	call *CallExpr
}

type ErrBuiltinMismatchedArgs struct {
	*CallExpr
	x, y reflect.Type
}

type ErrBuiltinNonTypeArg struct {
	Expr
}

type ErrBuiltinInvalidEllipsis struct {
	*CallExpr
}

type ErrMakeBadType struct {
	Expr
	of reflect.Type
}

type ErrMakeNonIntegerArg struct {
	Expr
	i int
}

type ErrMakeLenGtrThanCap struct {
	*CallExpr
	length, capacity int
}

type ErrAppendFirstArgNotSlice struct {
	Expr
}

type ErrAppendFirstArgNotVariadic struct {
	Expr
}

type ErrCopyArgsMustBeSlices struct {
	*CallExpr
	xT, yT reflect.Type
}

type ErrCopyArgsHaveDifferentEltTypes struct {
	*CallExpr
	xT, yT reflect.Type
}

type ErrDeleteFirstArgNotMap struct {
	Expr
}

func (err ErrBadBasicLit) Error() string {
	return fmt.Sprintf("Bad literal %v", err.BasicLit)
}

func (err ErrUndefined) Error() string {
	return fmt.Sprintf("undefined: %v", err.Expr)
}

func (err ErrInvalidIndexOperation) Error() string {
	t := err.IndexExpr.X.(Expr).KnownType()[0]
	return fmt.Sprintf("invalid operation: %v (index of type %v)", err.IndexExpr, t)
}

func (err ErrInvalidSliceIndex) Error() string {
	slice := err.SliceExpr
	return fmt.Sprintf("invalid slice index: %v > %v", slice.Low, slice.High)
}

func (err ErrInvalidSliceOperation) Error() string {
	x := err.SliceExpr.X.(Expr)
	xT := x.KnownType()[0]
	return fmt.Sprintf("cannot slice %v (type %v)", x, xT)
}

func (err ErrUnaddressableSliceOperand) Error() string {
	return fmt.Sprintf("invalid operation %v (slice of unaddressable value)", err.SliceExpr)
}

func (err ErrInvalidIndirect) Error() string {
	expr := err.Expr
	t := expr.KnownType()[0]
	if ct, ok := t.(ConstType); ok {
		if ct == ConstNil {
			return "invalid indirect of nil"
		}
		return fmt.Sprintf("invalid indirect of %v (type %s)",
			expr, ct.ErrorType())
	}
	return fmt.Sprintf("invalid indirect of %v (type %s)", expr, t)
}

func (err ErrUndefinedFieldOrMethod) Error() string {
	selector := err.Expr.(*SelectorExpr)
	t := selector.X.(Expr).KnownType()[0]
	return fmt.Sprintf("%v undefined (type %v has no field or method %v)",
		selector, t, selector.Sel.Name)
}

func (err ErrMissingValue) Error() string {
	return fmt.Sprintf("%v used as value", err.Expr)
}

func (err ErrMultiInSingleContext) Error() string {
	return fmt.Sprintf("multiple-value %v in single-value context", err.Expr)
}

func (err ErrBadMapIndex) Error() string {
	i := err.Expr
	iT := i.KnownType()[0]
	if _, ok := iT.(ConstType); ok {
		return fmt.Sprintf("cannot use %v as type %v in map index", i, err.keyT)
	} else {
		return fmt.Sprintf("cannot use %v (type %v) as type %v in map index", uc(i), iT, err.keyT)
	}
}

func (err ErrNonIntegerIndex) Error() string {
	i := err.Expr
	iT := i.KnownType()[0]
	var xname string
	if iT.Kind() == reflect.String {
		xname = "string"
	} else {
		xname = "array"
	}
	return fmt.Sprintf("non-integer %s index %v", xname, i)
}

func (err ErrIndexOutOfBounds) Error() string {
	i := err.Expr
	x := err.x
	var xname string
	var eltname string
	var length int
	if x.KnownType()[0].Kind() == reflect.String {
		length = x.Const().Len()
		xname = "string"
		eltname = "byte"
	} else {
		length = x.KnownType()[0].Len()
		xname = "array"
		eltname = "element"
	}
	if err.i < 0 {
		return fmt.Sprintf("invalid %s index %v (index must be non negative)",
			xname, i)
	} else {
		return fmt.Sprintf("invalid %s index %v (out of bounds for %d-%s %s)",
			xname, i, length, eltname, xname)
	}
}

func (err ErrCallNonFuncType) Error() string {
	expr := err.Expr
	return fmt.Sprintf("cannot call non-function %v (type %v)",
		expr, expr.KnownType()[0])
}

func (err ErrWrongNumberOfArgs) Error() string {
	call := err.CallExpr
	if call.isTypeConversion {
		to := call.KnownType()[0]
		if err.numArgs == 0 {
			return fmt.Sprintf("missing argument to conversion to %v", to)
		} else {
			return fmt.Sprintf("too many arguments to conversion to %v", to)
		}
	} else {
		if err.numArgs < call.Fun.(Expr).KnownType()[0].NumIn() {
			return fmt.Sprintf("not enough arguments in call to %v", call.Fun)
		} else {
			return fmt.Sprintf("too many arguments in call to %v", call.Fun)
		}
	}
}

func (err ErrWrongArgType) Error() string {
	ft := err.call.Fun.(Expr).KnownType()[0]
	var expected reflect.Type
	if ft.IsVariadic() && !err.call.argNEllipsis && err.argPos >= ft.NumIn() - 1 {
		expected = ft.In(ft.NumIn() - 1).Elem()
	} else {
		expected = ft.In(err.argPos)
	}

	if err.call.arg0MultiValued {
		actual := err.Expr.KnownType()[err.argPos]
		return fmt.Sprintf("cannot use %v as type %v in argument to %v",
			actual, expected, err.call.Fun)
	} else {
		arg := err.Expr
		actual := arg.KnownType()[0]
		return fmt.Sprintf("cannot use %v (type %v) as type %v in function argument",
			arg, actual, expected)
	}
}

func (err ErrInvalidEllipsisInCall) Error() string {
	fun := err.CallExpr.Fun
	return fmt.Sprintf("invalid use of ... in call to %v", fun)
}

func (err ErrInvalidUnaryOperation) Error() string {
	unary := err.UnaryExpr
	x := unary.X.(Expr)
	t := x.KnownType()[0]
	if ct, ok := t.(ConstType); ok {
		if unary.Op == token.XOR && ct.IsNumeric() {
			return fmt.Sprintf("illegal constant expression ^ %v", ct.ErrorType())
		}
		return fmt.Sprintf("invalid operation: %v %v", unary.Op, ct.ErrorType())
	}
	return fmt.Sprintf("invalid operation: %v %v", unary.Op, t)
}

func (err ErrInvalidAddressOf) Error() string {
	return fmt.Sprintf("cannot take the address of %v", err.Expr)
}

func (err ErrInvalidRecvFrom) Error() string {
	operand := err.Expr
	t := operand.KnownType()[0]
	var cause string
	if t.Kind() != reflect.Chan {
		cause = fmt.Sprintf("receive from non-chan type %v", t)
	} else {
		cause = fmt.Sprintf("receive from send-only type %v", t)
	}
	return fmt.Sprintf("invalid operation: <-%v (%s)", err.Expr, cause)
}

func (err ErrInvalidBinaryOperation) Error() string {
	binary := err.BinaryExpr
	op := binary.Op
	x := binary.X.(Expr)
	y := binary.Y.(Expr)

	xt := x.KnownType()[0]
	yt := y.KnownType()[0]

	xct, xcok := xt.(ConstType)
	yct, ycok := yt.(ConstType)

	if xcok && ycok {
		xn, xnok := x.Const().Interface().(*ConstNumber)
		yn, ynok := y.Const().Interface().(*ConstNumber)

		if xnok && ynok {
			switch op {
			case token.REM:
				if xn.Type.IsReal() && yn.Type.IsReal() {
					return "illegal constant expression: floating-point % operation"
				}
			}
			return fmt.Sprintf("illegal constant expression: %s %v %s", xct.ErrorType(), op, yct.ErrorType())
		} else if xt == yt {
			// const nil value prints as <T>, as an operand we should print nil
			var operandType interface{}
			if xt == ConstNil {
				operandType = "nil"
			} else {
				operandType = xt
			}
			return fmt.Sprintf("invalid operation: %v %v %v (operator %v not defined on %v)",
				x, op, y, op, operandType)
		}
	} else if xcok {
                // The gc implementation re-types nodes in const expressions, so that both sides
                // have type yt. We don't do this, so we will have to make the conversion again.
                // Runes get printed out verbatim
		var xFmt string
		if xt == ConstNil {
			// strings always produce mismatched types when
			// used with nil
			if yt.Kind() != reflect.String {
				xFmt = "nil"
			}
		} else {
			xx, _ := promoteConstToTyped(xct, constValue(x.Const()), yt, x)
			if reflect.Value(xx).IsValid() {
				xFmt = sprintConstValue(xt, reflect.Value(xx), false)
			}
		}
		if xFmt != "" && !isOpDefinedOn(op, yt) {
			ytFmt := sprintOperandType(yt)
                        return fmt.Sprintf("invalid operation: %v %v %v (operator %v not defined on %s)",
                                xFmt, op, y, op, ytFmt)
                }
	} else if ycok {
		var yFmt string
		if yt == ConstNil {
			if xt.Kind() != reflect.String {
				yFmt = "nil"
			}
		} else {
			yy, _ := promoteConstToTyped(yct, constValue(y.Const()), xt, y)
			if reflect.Value(yy).IsValid() {
				yFmt = sprintConstValue(yt, reflect.Value(yy), false)
			}
		}
		if yFmt != "" && !isOpDefinedOn(op, xt) {
			xtFmt := sprintOperandType(xt)
                        return fmt.Sprintf("invalid operation: %v %v %v (operator %v not defined on %s)",
                                x, op, yFmt, op, xtFmt)
		}
	} else {
		// Interfaces produce mismatched type errors unless
		// their types are identical
		var mismatch bool
		if xt.Kind() == reflect.Interface || yt.Kind() == reflect.Interface {
			mismatch = xt != yt
		} else {
			mismatch = !areTypesCompatible(xt, yt)
		}
		if !mismatch && !isOpDefinedOn(op, xt) {
			xtFmt := sprintOperandType(xt)
                        return fmt.Sprintf("invalid operation: %v %v %v (operator %v not defined on %s)",
                                x, op, y, op, xtFmt)
		} else if !mismatch && xt.Kind() == reflect.Struct {
			if field, ok := nonComparableField(xt); ok {
				return fmt.Sprintf("invalid operation: %v %v %v (struct containing %v cannot be compared)",
					x, op, y, field.Type)
			}
		} else if !mismatch && comparableToNilOnly(xt) {
			return fmt.Sprintf("invalid operation: %v %v %v (%v can only be compared to nil)",
				x, op, y, sprintOperandType(xt))
		}
        }

        // This hack is again to do with the retyping, if half the expression is
        // typed, then the untyped half of the expression assumes its default type.
        var xi, yi interface{} = x, y
        if !ycok {
                xi = sprintUntypedConstAsTyped(x)
        }
        if !xcok {
                yi = sprintUntypedConstAsTyped(y)
        }
        // One last hack to display nil types as "nil", not the usual "<T>"
        var xti, yti interface{} = xt, yt
        if !ycok && xt == ConstNil {
                xti = "nil"
        }
        if !xcok && yt == ConstNil {
                yti = "nil"
        }
	return fmt.Sprintf("invalid operation: %v %v %v (mismatched types %v and %v)",
		xi, op, yi, xti, yti,
	)
}

func (err ErrDivideByZero) Error() string {
	return "division by zero"
}

func (err ErrBadConversion) Error() string {
	return fmt.Sprintf("cannot convert %v (type %v) to type %v", err.Expr, err.from, err.to)
}

func (err ErrBadConstConversion) Error() string {
	return fmt.Sprintf("cannot convert %v to type %v", err.Expr, err.to)
}

func (err ErrTruncatedConstant) Error() string {
	if err.to.IsIntegral() {
		return fmt.Sprintf("constant %v truncated to integer", err.constant)
	} else {
		return fmt.Sprintf("constant %v truncated to real", err.constant)
	}
}

func (err ErrOverflowedConstant) Error() string {
	switch err.to.(type) {
	case ConstStringType:
		return fmt.Sprintf("overflow in int -> string")
	default:
		var constant string

		// Runes print their actual value in overflow errors
		if err.constant.Type == ConstRune {
			constant = err.constant.Value.Re.Num().String()
		} else {
			constant = err.constant.String()
		}

		return fmt.Sprintf("constant %v overflows %v", constant, err.to)
	}
}

func (ErrUntypedNil) Error() string {
	return "use of untyped nil"
}

func (err ErrTypeUsedAsExpression) Error() string {
	t := err.Expr
	return fmt.Sprintf("type %v is not an expression", t)
}

func (err ErrUncomparableMapKey) Error() string {
	return fmt.Sprintf("invalid map key type %v", err.keyT)
}

func (err ErrMissingMapKey) Error() string {
	return "missing key in map literal"
}

func (err ErrBadMapKey) Error() string {
	expr := err.Expr
	t := expr.KnownType()[0]
	if t == ConstNil {
		return fmt.Sprintf("cannot use nil as type %v in map key", err.keyT)
	}
	return fmt.Sprintf("cannot use %v (type %v) as type %v in map key",
		expr, t, err.keyT)
}

func (err ErrDuplicateMapKey) Error() string {
	key := err.Expr
	return fmt.Sprintf("duplicate key %v in map literal", key)
}

func (err ErrBadMapValue) Error() string {
	expr := err.Expr
	t := expr.KnownType()[0]
	if t == ConstNil {
		return fmt.Sprintf("cannot use nil as type %v in map value", err.eltT)
	}
	return fmt.Sprintf("cannot use %v (type %v) as type %v in map value",
		expr, t, err.eltT)
}

func (ErrBadArrayKey) Error() string {
	return "array index must be non-negative integer constant"
}

func (err ErrArrayKeyOutOfBounds) Error() string {
	length := err.arrayT.Len()
	return fmt.Sprintf("array index %d out of bounds [0:%d]", err.index+1, length)
}

func (err ErrDuplicateArrayKey) Error() string {
	return fmt.Sprintf("duplicate index in array literal: %v", err.index)
}

func (err ErrBadArrayValue) Error() string {
	expr := err.Expr
	t := expr.KnownType()[0]
	if t == ConstNil {
		return fmt.Sprintf("cannot use nil as type %v in array element", err.eltT)
	}
	return fmt.Sprintf("cannot use %v (type %v) as type %v in array element",
		expr, t, err.eltT)
}

func (err ErrUnknownStructField) Error() string {
	return fmt.Sprintf("unknown %v field '%v' in struct literal",
		err.structT, err.field)
}

func (err ErrInvalidStructField) Error() string {
	return fmt.Sprintf("invalid field name %v in struct initializer", err.Expr)
}

func (err ErrDuplicateStructField) Error() string {
	return fmt.Sprintf("duplicate field name in struct literal: %v", err.field)
}

func (err ErrMixedStructValues) Error() string {
	return fmt.Sprintf("mixture of field:value and value initializers")
}

func (err ErrWrongNumberOfStructValues) Error() string {
	lit := err.CompositeLit
	actual := len(lit.Elts)
	expected := lit.KnownType()[0].NumField()
	if actual < expected {
		return fmt.Sprintf("too few values in struct initializer")
	} else {
		return fmt.Sprintf("too many values in struct initializer")
	}
}

func (err ErrBadStructValue) Error() string {
	expr := err.Expr
	t := expr.KnownType()[0]
	if t == ConstNil {
		return fmt.Sprintf("cannot use nil as type %v in field value", err.eltT)
	}
	return fmt.Sprintf("cannot use %v (type %v) as type %v in field value",
		expr, t, err.eltT)
}

func (err ErrInvalidTypeAssert) Error() string {
	assert := err.TypeAssertExpr
	xT := assert.X.(Expr).KnownType()[0]
	return fmt.Sprintf("invalid type assertion: %v (non-interface type %v on left)",
		assert, xT)
}

func (err ErrImpossibleTypeAssert) Error() string {
	assert := err.TypeAssertExpr
	iT := assert.KnownType()[0]
	xT := assert.X.(Expr).KnownType()[0]

	var missingMethod string
	numMethod := iT.NumMethod()
	for i := 0; i < numMethod; i += 1 {
		missingMethod = iT.Method(i).Name
		if _, ok := xT.MethodByName(missingMethod); !ok {
			break
		}
	}

	return fmt.Sprintf("impossible type assertion:\n" +
		"\t%v does not implement %v (missing %s method)",
		xT, iT, missingMethod)
}

func (err ErrMissingCompositeLitType) Error() string {
	return "missing type in composite literal"
}

func (err ErrBuiltinWrongNumberOfArgs) Error() string {
	call := err.CallExpr
	ident := call.Fun.(*Ident)
	tooMany := false
	plural := ""
	var cause string
	switch ident.Name {
	case "complex":
		if len(call.Args) == 0 {
			cause = " - complex(<N>, <N>)"
		} else {
			tooMany = len(call.Args) > 2
			cause = fmt.Sprintf(" - complex(%v, <N>)", uc(call.Args[0].(Expr)))
		}
	case "new":
		if len(call.Args) != 0 {
			tooMany = true
			cause = fmt.Sprintf("(%v)", uc(call.Args[0].(Expr)))
		}
	case "make":
		if len(call.Args) == 1 {
			return fmt.Sprintf("too few arguments to make: %v", uc(call))
		} else if len(call.Args) != 0 {
			tooMany = true
			cause = fmt.Sprintf(": %v", uc(call))
		}
	case "copy":
		if len(call.Args) < 2 {
			plural = "s"
		} else if len(call.Args) != 0 {
			tooMany = true
		}
	case "delete":
		if len(call.Args) == 0 {
			plural = "s"
		} else if len(call.Args) == 1 {
			return "missing second (key) argument to delete"
		} else {
			tooMany = true
		}
	case "append":
		// Note the s on arguments, which
		return "missing arguments to append"
	default:
		cause = fmt.Sprintf(": %v", uc(call))
		tooMany = len(call.Args) != 0
	}
	if tooMany {
		return fmt.Sprintf("too many arguments to %s%s", ident.Name, cause)
	} else {
		return fmt.Sprintf("missing argument%s to %s%s", plural, ident.Name, cause)
	}
}

func (err ErrBuiltinWrongArgType) Error() string {
	ident := err.call.Fun.(*Ident)
	arg := err.Expr
	var t string
	kt := arg.KnownType()[0]
	if ct, ok := kt.(ConstType); ok {
		t = ct.ErrorType()
	} else {
		t = kt.String()
	}
	switch ident.Name {
	case "complex":
		call := uc(err.call).(*CallExpr)
		// ... doesn't get printed. uc() returns a clone of the root node, so we can safely change argNEllipsis
		call.argNEllipsis = false
		return fmt.Sprintf("invalid operation: %v (arguments have type %s, expected floating-point)",
			call, t)
	case "append":
		expected := err.call.Args[0].(Expr).KnownType()[0].Elem()
		if kt == ConstNil {
			return fmt.Sprintf("cannot use nil as type %v in append", expected)
		}
		return fmt.Sprintf("cannot use %v (type %s) as type %v in append", uc(arg), kt, expected)
	case "delete":
		expected := err.call.Args[0].(Expr).KnownType()[0].Key()
		if kt == ConstNil {
			return fmt.Sprintf("cannot use nil as type %v in delete", expected)
		}
		return fmt.Sprintf("cannot use %v (type %s) as type %v in delete", uc(arg), kt, expected)
	default:
		return fmt.Sprintf("invalid argument %v (type %s) for %s", uc(arg), kt, ident.Name)
	}
}

func (err ErrBuiltinMismatchedArgs) Error() string {
	call := err.CallExpr
	var x, y string
	cx, cxok := err.x.(ConstType)
	cy, cyok := err.y.(ConstType)
	if cxok && cyok {
		x = cx.ErrorType()
		y = cy.ErrorType()
	} else if cx == ConstNil {
		x = "nil"
		y = err.y.String()
	} else if cy == ConstNil {
		x = err.x.String()
		y = "nil"
	} else {
		x = err.x.String()
		y = err.y.String()
	}
	call = uc(call).(*CallExpr)
	// ... doesn't get printed. uc() returns a clone of the root node, so we can safely change argNEllipsis
	call.argNEllipsis = false
	return fmt.Sprintf("invalid operation: %v (mismatched types %s and %s)", uc(call), x, y)
}

func (err ErrBuiltinNonTypeArg) Error() string {
	return fmt.Sprintf("%v is not a type", uc(err.Expr))
}

func (err ErrBuiltinInvalidEllipsis) Error() string {
	ident := err.CallExpr.Fun.(*Ident)
	return fmt.Sprintf("invalid use of ... with builtin %s", ident.Name)
}

func (err ErrMakeBadType) Error() string {
	return fmt.Sprintf("cannot make type %v", err.of)
}

func (err ErrMakeNonIntegerArg) Error() string {
	var culprit string
	if err.i == 1 {
		culprit = "len"
	} else {
		culprit = "cap"
	}
	return fmt.Sprintf("make: non-integer %s argument %v", culprit, uc(err.Expr))
}

func (err ErrMakeLenGtrThanCap) Error() string {
	return fmt.Sprintf("len larger than cap in %v", err.CallExpr)
}

func (err ErrAppendFirstArgNotSlice) Error() string {
	arg := err.Expr
	t := arg.KnownType()[0]
	if t == ConstNil {
		return "first argument to append must be typed slice; have untyped nil"
	} else {
		var s string
		if ct, ok := t.(ConstType); ok {
			s = ct.ErrorType()
		} else {
			s = t.String()
		}
		return fmt.Sprintf("first argument to append must be slice; have %s", s)
	}
}

func (err ErrAppendFirstArgNotVariadic) Error() string {
	return "cannot use ... on first argument to append"
}

func (err ErrCopyArgsMustBeSlices) Error() string {
	if err.yT != ConstNil && err.yT.Kind() == reflect.Slice {
		return fmt.Sprintf("first argument to copy should be slice; have %v", err.xT)
	} else if err.xT != ConstNil && err.xT.Kind() == reflect.Slice {
		return fmt.Sprintf("second argument to copy should be slice or string; have %v", err.yT)
	} else {
		return fmt.Sprintf("arguments to copy must be slices; have %v, %v", err.xT, err.yT)
	}
}

func (err ErrCopyArgsHaveDifferentEltTypes) Error() string {
	return fmt.Sprintf("arguments to copy have different element types: %v and %v", err.xT, err.yT)
}

func (err ErrDeleteFirstArgNotMap) Error() string {
	arg := err.Expr
	t := arg.KnownType()[0]
	var s string
	if ct, ok := t.(ConstType); ok {
		s = ct.ErrorType()
	} else {
		s = t.String()
	}
	return fmt.Sprintf("first argument to delete must be map; have %s", s)
}

// For display purposes only, display untyped const nodes as they would be
// displayed as a typed const node.
func sprintUntypedConstAsTyped(expr Expr) string {
        if !expr.IsConst() {
                return expr.String()
        }
        switch expr.KnownType()[0].(type) {
        case ConstRuneType:
                return sprintConstValue(RuneType, reflect.Value(expr.Const()), false)
        default:
                return expr.String()
        }
}

// Determines if two types can be automatically converted between.
func areTypesCompatible(xt, yt reflect.Type) bool {
	return xt.AssignableTo(unhackType(yt)) || yt.AssignableTo(unhackType(xt))
}

func sprintOperandType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Array:
		return "array"
	case reflect.Slice:
		return "slice"
	case reflect.Interface:
		return "interface"
	case reflect.Ptr:
		return "pointer"
	case reflect.Struct:
		return "struct"
	case reflect.Map:
		return "map"
	default:
		return t.String()
	}
}
