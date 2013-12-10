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
	vals []reflect.Value
}

type ErrArrayIndexOutOfBounds struct {
	t reflect.Type
	i uint64
}

type ErrorContext struct {
	Input string
	ast.Node
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

func (err ErrMissingValue) Error() string {
	return fmt.Sprintf("%s used as value", err.ErrorContext.Source())
}

func (err ErrMultiInSingleContext) Error() string {
	return fmt.Sprintf("multiple-value %s in single-value context", err.ErrorContext.Source())
}

func (err ErrArrayIndexOutOfBounds) Error() string {
	return fmt.Sprintf("array index %d out of bounds [0:%d]", err.i, err.t.Len())
}

func at(ctx *Ctx, expr ast.Node) ErrorContext {
	return ErrorContext{ctx.Input, expr}
}

func (errCtx ErrorContext) Source() string {
	return errCtx.Input[errCtx.Node.Pos()-1:errCtx.Node.End()-1]
}
