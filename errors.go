package interactive

import (
	"errors"
	"fmt"
	"reflect"

	"go/token"
)

var (
	ErrArrayKey = errors.New("array index must be non-negative integer constant")
	ErrMissingValue = errors.New("void used as value")
)

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

type ErrMultiInSingleContext struct {
	vals []reflect.Value
}

type ErrArrayIndexOutOfBounds struct {
	t reflect.Type
	i uint64
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

func (err ErrMultiInSingleContext) Error() string {
	strvals := ""
	for i, val := range err.vals {
		if i == 0 {
			strvals += fmt.Sprintf("%v", val)
		} else {
			strvals += fmt.Sprintf(", %v", val)
		}
	}
	return fmt.Sprintf("multiple-value (%s) in single value context", strvals)
}

func (err ErrArrayIndexOutOfBounds) Error() string {
	return fmt.Sprintf("array index %d out of bounds [0:%d]", err.i, err.t.Len())
}


