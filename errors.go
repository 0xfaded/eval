package eval

import (
	"fmt"
	"reflect"

	"go/ast"
	"go/token"
)

type ErrBadBasicLit struct {
	ErrorContext
}

type ErrUndefined struct {
	ErrorContext
}

type ErrInvalidOperand struct {
	x reflect.Value
	op token.Token
}

type ErrInvalidIndirect struct {
	ErrorContext
}

type ErrUndefinedFieldOrMethod struct {
	ErrorContext
}

type ErrMismatchedTypes struct {
	x reflect.Value
	op token.Token
	y reflect.Value
}

type ErrInvalidOperands struct {
	x reflect.Value
	op token.Token
	y reflect.Value
}

type ErrBadFunArgument struct {
	fun reflect.Value
	index int
	value reflect.Value
}

type ErrBadComplexArguments struct {
	a, b reflect.Value
}

type ErrBadBuiltinArgument struct {
	fun string
	value reflect.Value
}

type ErrCallNonFuncType struct {
	ErrorContext
}

type ErrWrongNumberOfArgs struct {
	ErrorContext
	numArgs int
}

type ErrWrongArgType struct {
	ErrorContext
	call *CallExpr
	argPos int
}

type ErrInvalidEllipsisInCall struct {
	ErrorContext
}

type ErrMissingValue struct {
	ErrorContext
}

type ErrMultiInSingleContext struct {
	ErrorContext
}

type ErrNonIntegerIndex struct {
	ErrorContext
}

type ErrIndexOutOfBounds struct {
	ErrorContext
	x Expr
	i int
}

type ErrInvalidIndexOperation struct {
	ErrorContext
}

type ErrInvalidSliceIndex struct {
	ErrorContext
}

type ErrInvalidSliceOperation struct {
	ErrorContext
}

type ErrUnaddressableSliceOperand struct {
	ErrorContext
}

type ErrInvalidIndex struct {
	ErrorContext
	indexValue reflect.Value
	containerType reflect.Type
}

type ErrDivideByZero struct {
	ErrorContext
}

type ErrInvalidBinaryOperation struct {
	ErrorContext
}

type ErrInvalidUnaryOperation struct {
	ErrorContext
}

type ErrInvalidAddressOf struct {
	ErrorContext
}

type ErrBadConversion struct {
	ErrorContext
	from reflect.Type
	to reflect.Type
	v reflect.Value
}

type ErrBadConstConversion struct {
	ErrorContext
	from reflect.Type
	to reflect.Type
	v reflect.Value
}

type ErrTruncatedConstant struct {
	ErrorContext
	to ConstType
	constant *ConstNumber
}

type ErrOverflowedConstant struct {
	ErrorContext
	from ConstType
	to reflect.Type
	constant *ConstNumber
}

type ErrUntypedNil struct {
	ErrorContext
}

type ErrBadArrayKey struct {
	ErrorContext
}

type ErrArrayKeyOutOfBounds struct {
	ErrorContext
	arrayT reflect.Type
	index int
}

type ErrDuplicateArrayKey struct {
	ErrorContext
	index int
}

type ErrBadArrayValue struct {
	ErrorContext
	eltT reflect.Type
}

type ErrUnknownStructField struct {
	ErrorContext
	structT reflect.Type
	field string
}

type ErrInvalidStructField struct {
	ErrorContext
}

type ErrDuplicateStructField struct {
	ErrorContext
	field string
}

type ErrMixedStructValues struct {
	ErrorContext
}

type ErrWrongNumberOfStructValues struct {
	ErrorContext
}

type ErrBadStructValue struct {
	ErrorContext
	eltT reflect.Type
}

type ErrInvalidTypeAssert struct {
	ErrorContext
}

type ErrImpossibleTypeAssert struct {
	ErrorContext
}

type ErrorContext struct {
	Input string
	ast.Node
}

func (err ErrBadBasicLit) Error() string {
	return fmt.Sprintf("Bad literal %s", err.Source())
}

func (err ErrUndefined) Error() string {
	return fmt.Sprintf("undefined: %v", err.Node)
}

func (err ErrInvalidOperand) Error() string {
	return fmt.Sprintf("invalid unary operation %v %v", err.op, err.x)
}

func (err ErrInvalidOperands) Error() string {
	return fmt.Sprintf("invalid binary operation %v %v %v", err.x, err.op, err.y)
}

func (err ErrMismatchedTypes) Error() string {
	return fmt.Sprintf("invalid operation %v %v %v (mismatched types %s and %s)", err.x, err.op, err.y, err.x.Kind(), err.y.Kind())
}

func (err ErrBadFunArgument) Error() string {
	return fmt.Sprintf("invalid type (%v) for argument %d of %v", err.value.Type(), err.index, err.fun)
}

func (err ErrBadComplexArguments) Error() string {
	return fmt.Sprintf("invalid operation: complex(%v, %v)", err.a, err.b)
}

func (err ErrBadBuiltinArgument) Error() string {
	return fmt.Sprintf("invalid operation: %s(%v)", err.fun, err.value)
}

func (err ErrInvalidIndexOperation) Error() string {
	t := err.Node.(*IndexExpr).X.(Expr).KnownType()[0]
	return fmt.Sprintf("invalid operation: %s (index of type %v)", err.Source(), t)
}

func (err ErrInvalidSliceIndex) Error() string {
	slice := err.Node.(*SliceExpr)
	return fmt.Sprintf("invalid slice index: %v > %v", slice.Low, slice.High)
}

func (err ErrInvalidSliceOperation) Error() string {
	x := err.Node.(*SliceExpr).X.(Expr)
	xT := x.KnownType()[0]
	return fmt.Sprintf("cannot slice %v (type %v)", x, xT)
}

func (err ErrUnaddressableSliceOperand) Error() string {
	return fmt.Sprintf("invalid operation %v (slice of unaddressable value)", err.Node)
}

func (err ErrInvalidIndirect) Error() string {
	expr := err.Node.(Expr)
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
	selector := err.Node.(*SelectorExpr)
	t := selector.X.(Expr).KnownType()[0]
	return fmt.Sprintf("%v undefined (type %v has no field or method %v)",
		selector, t, selector.Sel.Name)
}

func (err ErrMissingValue) Error() string {
	return fmt.Sprintf("%s used as value", err.ErrorContext.Source())
}

func (err ErrMultiInSingleContext) Error() string {
	return fmt.Sprintf("multiple-value %s in single-value context", err.ErrorContext.Source())
}

func (err ErrNonIntegerIndex) Error() string {
	i := err.Node.(Expr)
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
	i := err.Node.(Expr)
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
	expr := err.Node.(Expr)
	return fmt.Sprintf("cannot call non-function %v (type %v)",
		expr, expr.KnownType()[0])
}

func (err ErrWrongNumberOfArgs) Error() string {
	call := err.ErrorContext.Node.(*CallExpr)
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
		actual := err.Node.(Expr).KnownType()[err.argPos]
		return fmt.Sprintf("cannot use %v as type %v in argument to %v",
			actual, expected, err.call.Fun)
	} else {
		arg := err.Node.(Expr)
		actual := arg.KnownType()[0]
		return fmt.Sprintf("cannot use %v (type %v) as type %v in function argument",
			arg, actual, expected)
	}
}

func (err ErrInvalidEllipsisInCall) Error() string {
	fun := err.Node.(*CallExpr).Fun
	return fmt.Sprintf("invalid use of ... in call to %v", fun)
}

func (err ErrInvalidUnaryOperation) Error() string {
	unary := err.ErrorContext.Node.(*UnaryExpr)
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
	return fmt.Sprintf("cannot take the address of %v", err.Node)
}

func (err ErrInvalidBinaryOperation) Error() string {
	binary := err.ErrorContext.Node.(*BinaryExpr)
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
			return fmt.Sprintf("illegal constant expression: ideal %v ideal", op)
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
			xFmt = "nil"
		} else {
			xx, _ := promoteConstToTyped(&Ctx{}, xct, constValue(x.Const()), yt, x)
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
			yFmt = "nil"
		} else {
			yy, _ := promoteConstToTyped(&Ctx{}, yct, constValue(y.Const()), xt, y)
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
		operandT := xt
		var mismatch bool
		// Interfaces are wierd. Non-bool ops produce mismatched type
		// errors in preference to operation not defined. Bool ops
		// are the reverse.
		if xt.Kind() == reflect.Interface || yt.Kind() == reflect.Interface {
			mismatch = xt != yt
		} else {
			mismatch = !areTypesCompatible(xt, yt)
		}
		if !mismatch && !isOpDefinedOn(op, operandT) {
			xtFmt := sprintOperandType(xt)
                        return fmt.Sprintf("invalid operation: %v %v %v (operator %v not defined on %s)",
                                x, op, y, op, xtFmt)
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
	return fmt.Sprintf("cannot convert %v (type %v) to type %v", err.Node.(Expr), err.from, err.to)
}

func (err ErrBadConstConversion) Error() string {
	return fmt.Sprintf("cannot convert %v to type %v", err.Node.(Expr), err.to)
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
	expr := err.Node.(Expr)
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
	return fmt.Sprintf("invalid field name %v in struct initializer", err.Node)
}

func (err ErrDuplicateStructField) Error() string {
	return fmt.Sprintf("duplicate field name in struct literal: %v", err.field)
}

func (err ErrMixedStructValues) Error() string {
	return fmt.Sprintf("mixture of field:value and value initializers")
}

func (err ErrWrongNumberOfStructValues) Error() string {
	lit := err.Node.(*CompositeLit)
	actual := len(lit.Elts)
	expected := lit.KnownType()[0].NumField()
	if actual < expected {
		return fmt.Sprintf("too few values in struct initializer")
	} else {
		return fmt.Sprintf("too many values in struct initializer")
	}
}

func (err ErrBadStructValue) Error() string {
	expr := err.Node.(Expr)
	t := expr.KnownType()[0]
	if t == ConstNil {
		return fmt.Sprintf("cannot use nil as type %v in field value", err.eltT)
	}
	return fmt.Sprintf("cannot use %v (type %v) as type %v in field value",
		expr, t, err.eltT)
}

func (err ErrInvalidTypeAssert) Error() string {
	assert := err.Node.(*TypeAssertExpr)
	xT := assert.X.(Expr).KnownType()[0]
	return fmt.Sprintf("invalid type assertion: %v (non-interface type %v on left)",
		assert, xT)
}

func (err ErrImpossibleTypeAssert) Error() string {
	assert := err.Node.(*TypeAssertExpr)
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

func at(ctx *Ctx, expr ast.Node) ErrorContext {
	return ErrorContext{ctx.Input, expr}
}

func (errCtx ErrorContext) Source() string {
	return errCtx.Input[errCtx.Node.Pos()-1:errCtx.Node.End()-1]
}

func drop0i(i interface{}) interface{} {
	if n, ok := i.(*ConstNumber); ok {
		return n.StringShow0i(false)
	}
	return i
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
	default:
		return t.String()
	}
}
