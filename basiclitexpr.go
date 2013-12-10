package interactive

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"go/ast"
	"go/token"
)

func evalBasicLit(ctx *Ctx, lit *ast.BasicLit) (reflect.Value, bool, error) {
	switch lit.Kind {
	case token.CHAR:
		if r, _, tail, err := strconv.UnquoteChar(lit.Value[1:len(lit.Value)-1], '\''); err != nil {
			return reflect.Value{}, false, ErrBadBasicLit{at(ctx, lit)}
		} else if tail != "" {
			// parser.ParseExpr() should raise a syntax error before we get here.
			panic("go-interactive: bad char lit " + lit.Value)
		} else {
			return reflect.ValueOf(r), false, nil
		}
	case token.STRING:
		str, err := strconv.Unquote(string(lit.Value))
		return reflect.ValueOf(str), true, err
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
