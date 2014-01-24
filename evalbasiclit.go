package eval

import (
	"reflect"
)

func evalBasicLit(ctx *Ctx, lit *BasicLit) (reflect.Value, error) {
	if lit.IsConst() {
		return lit.Const(), nil
	}
	panic(dytc("non-const basic lit"))
}
