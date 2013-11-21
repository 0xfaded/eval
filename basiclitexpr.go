package interactive

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"go/ast"
	"go/token"
)

func evalBasicLit(lit *ast.BasicLit) (reflect.Value, bool, error) {
	switch lit.Kind {
	case token.STRING:
		return reflect.ValueOf(lit.Value), false, nil
	case token.INT:
		i, err := strconv.ParseInt(lit.Value, 0, 0)
		return reflect.ValueOf(i), false, err
	case token.FLOAT:
		f, err := strconv.ParseFloat(lit.Value, 64)
		return reflect.ValueOf(f), false, err
	case token.IMAG:
		f, err := strconv.ParseFloat(lit.Value[:len(lit.Value)-1], 64)
		return reflect.ValueOf(complex(0, f)), false, err
	default:
		return reflect.Value{}, false, errors.New(fmt.Sprintf("BasicLit: Bad token type (%+v)", lit))
	}
}

