package eval

import (
	"errors"
	"fmt"
	"reflect"

	"go/ast"
	"go/token"
)

var (
	ErrArrayKey = errors.New("array index must be non-negative integer constant")
)

type ErrBadBasicLit struct {
	ErrorContext
}

type ErrInvalidOperand struct {
	x reflect.Value
	op token.Token
}

type ErrInvalidIndirect struct {
	t reflect.Type
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

// TODO remove when checker complete
type ErrWrongNumberOfArgsOld struct {
	fun reflect.Value
	numArgs int
}

type ErrWrongNumberOfArgs struct {
	ErrorContext
}

type ErrMissingValue struct {
	ErrorContext
}

type ErrMultiInSingleContext struct {
	ErrorContext
}

type ErrArrayIndexOutOfBounds struct {
	t reflect.Type
	i uint64
}

type ErrInvalidIndexOperation struct {
	ErrorContext
	t reflect.Type
}

type ErrInvalidSliceType struct {
	ErrorContext
	t reflect.Type
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

type ErrorContext struct {
	Input string
	ast.Node
}

func (err ErrBadBasicLit) Error() string {
	return fmt.Sprintf("Bad literal %s", err.Source())
}

func (err ErrInvalidOperand) Error() string {
	return fmt.Sprintf("invalid unary operation %v %v %v", err.x, err.op)
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

func (err ErrWrongNumberOfArgsOld) Error() string {
	expected := err.fun.Type().NumIn()
	if err.numArgs < expected {
		return fmt.Sprintf("not enough args (%d) to call %v (%d)", err.numArgs, err.fun, expected)
	} else {
		return fmt.Sprintf("too many args (%d) to call %v (%d)", err.numArgs, err.fun, expected)
	}
}

func (err ErrInvalidIndexOperation) Error() string {
	return fmt.Sprintf("invalid operation: %s (index of type %v)", err.Source(), err.t)
}

func (err ErrInvalidSliceType) Error() string {
	return fmt.Sprintf("cannot slice (type %v)", err.t)
}

func (err ErrInvalidIndex) Error() string {
	var ct string

	switch err.containerType.Kind() {
	//case reflect.Map:
	case reflect.Array:
		ct = "array"
	case reflect.Slice:
		ct = "slice"
	case reflect.String:
		ct = "string"
	default:
		panic("go-interactive error: ErrInvalidIndex requires indexable err.containerType")
	}

	switch err.indexValue.Type().Kind() {
	case reflect.Int:
		var reason string
		i := int(err.indexValue.Int())
		if err.containerType.Kind() == reflect.Array && i >= err.containerType.Len() {
			reason = fmt.Sprintf("out of bounds for %d-element array", err.containerType.Len())

		// TODO string errors for constant strings constant values are implemented
		// out of bounds for 3-byte string
		// } else if err.containerType.Kind == reflect.String && i > err.containerType.Len() {
		} else {
			reason = "index must be non-negative"
		}
		return fmt.Sprintf("invalid %s index %s (%s)", ct, err.Source(), reason)
	default:
		return fmt.Sprintf("non-integer %s index %s", ct, err.Source())
	}
}

func (err ErrInvalidIndirect) Error() string {
	return fmt.Sprintf("invalid indirect (type %v)", err.t)
}

func (err ErrMissingValue) Error() string {
	return fmt.Sprintf("%s used as value", err.ErrorContext.Source())
}

func (err ErrMultiInSingleContext) Error() string {
	return fmt.Sprintf("multiple-value %s in single-value context", err.ErrorContext.Source())
}

func (err ErrArrayIndexOutOfBounds) Error() string {
	return fmt.Sprintf("array index %d out of bounds [0:%d]", err.i, err.t.Len())
}

func (err ErrWrongNumberOfArgs) Error() string {
	call := err.ErrorContext.Node.(*CallExpr)
	if call.isTypeConversion {
		to := call.KnownType()[0]
		if call.Args == nil {
			return fmt.Sprintf("missing argument to conversion to %v", to)
		} else {
			return fmt.Sprintf("too many arguments to conversion to %v", to)
		}
	}
	return "TODO ErrWrongNumberOfArgs"
}

func (err ErrInvalidUnaryOperation) Error() string {
	unary := err.ErrorContext.Node.(*UnaryExpr)
	x := unary.X.(Expr)
	t := x.KnownType()[0]
	switch t {
	case ConstNil:
		return fmt.Sprintf("invalid operation: %v nil", unary.Op)
	case ConstBool:
		return fmt.Sprintf("invalid operation: %v ideal bool", unary.Op)
	case ConstString:
		return fmt.Sprintf("invalid operation: %v ideal string", unary.Op)
	default:
		if unary.Op == token.XOR {
			return fmt.Sprintf("illegal constant expression %v ideal", unary.Op)
		} else {
			return fmt.Sprintf("invalid operation: %v ideal", unary.Op)
		}
	}
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

	xn, xnok := x.Const().Interface().(*ConstNumber)
	yn, ynok := y.Const().Interface().(*ConstNumber)

	if xcok && ycok {

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
                xx, errs := promoteConstToTyped(&Ctx{}, xct, constValue(x.Const()), yt, x)
                mismatch := false
                if errs != nil {
                        if _, ok := errs[0].(ErrBadConstConversion); ok {
                                mismatch = true
                        }
                }
		if !mismatch && !isOpDefinedOn(op, yt) {
                        return fmt.Sprintf("invalid operation: %v %v %v (operator %v not defined on %v)",
                                sprintConstValue(xt, reflect.Value(xx), false), op, y, op, yt)
                }
	} else if ycok {
                yy, errs := promoteConstToTyped(&Ctx{}, yct, constValue(y.Const()), xt, y)
                mismatch := false
                if errs != nil {
                        if _, ok := errs[0].(ErrBadConstConversion); ok {
                                mismatch = true
                        }
                }
		if !mismatch && !isOpDefinedOn(op, xt) {
                        return fmt.Sprintf("invalid operation: %v %v %v (operator %v not defined on %v)",
                                x, op, sprintConstValue(yt, reflect.Value(yy), false), op, xt)
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
