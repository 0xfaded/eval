package interactive

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

type ErrWrongNumberOfArgs struct {
	fun reflect.Value
	numArgs int
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
	target reflect.Type
}

type ErrTruncatedConstant struct {
	ErrorContext
	constant *BigComplex
}

type ErrIllegalConstantExpr struct {
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

func (err ErrBadFunArgument) Error() string {
	return fmt.Sprintf("invalid type (%v) for argument %d of %v", err.value.Type(), err.index, err.fun)
}

func (err ErrBadComplexArguments) Error() string {
	return fmt.Sprintf("invalid operation: complex(%v, %v)", err.a, err.b)
}

func (err ErrBadBuiltinArgument) Error() string {
	return fmt.Sprintf("invalid operation: %s(%v)", err.fun, err.value)
}

func (err ErrWrongNumberOfArgs) Error() string {
	expected := err.fun.Type().NumIn()
	if err.numArgs < expected {
		return fmt.Sprintf("not enouch args (%d) to call %v (%d)", err.numArgs, err.fun, expected)
	} else {
		return fmt.Sprintf("too many args (%d) to call %v (%d)", err.numArgs, err.fun, expected)
	}
}

func (err ErrInvalidIndexOperation) Error() string {
	return fmt.Sprintf("invalid operation: %s (index of type %v)", err.Source(), err.t)
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

func (err ErrMissingValue) Error() string {
	return fmt.Sprintf("%s used as value", err.ErrorContext.Source())
}

func (err ErrMultiInSingleContext) Error() string {
	return fmt.Sprintf("multiple-value %s in single-value context", err.ErrorContext.Source())
}

func (err ErrArrayIndexOutOfBounds) Error() string {
	return fmt.Sprintf("array index %d out of bounds [0:%d]", err.i, err.t.Len())
}

func (err ErrInvalidUnaryOperation) Error() string {
	return "ErrInvalidUnaryOperation TODO"
}

func (err ErrInvalidBinaryOperation) Error() string {
	return "ErrInvalidBinaryOperation TODO"
}

func (err ErrIllegalConstantExpr) Error() string {
	return "ErrInvalidBinaryOperation TODO"
}

func (err ErrDivideByZero) Error() string {
	return "division by zero"
}

func (err ErrBadConversion) Error() string {
	return "division by zero"
}

func (err ErrTruncatedConstant) Error() string {
	if err.constant.IsReal() {
		return fmt.Sprintf("constant %v truncated to integer", err.constant)
	} else {
		return fmt.Sprintf("constant %v truncated to real", err.constant)
	}
}

func at(ctx *Ctx, expr ast.Node) ErrorContext {
	return ErrorContext{ctx.Input, expr}
}

func (errCtx ErrorContext) Source() string {
	return errCtx.Input[errCtx.Node.Pos()-1:errCtx.Node.End()-1]
}
