package interactive

import (
	"errors"
	"fmt"
	"reflect"
)

func evalIdentExpr(ctx *Ctx, ident *Ident, env *Env) (*reflect.Value, bool, error) {
	name := ident.Name
	if name == "nil" {
		// FIXME: Should this be done first or last?
		return nil, false, nil
	} else if v, ok := env.Vars[name]; ok {
		elem := v.Elem()
		return &elem, true, nil
	} else if v, ok := env.Consts[name]; ok {
		return &v, false, nil
	} else if v, ok := env.Funcs[name]; ok {
		return &v, true, nil
	} else if v, ok := builtinFuncs[name]; ok {
		return &v, false, nil
	} else if p, ok := env.Pkgs[name]; ok {
		val := reflect.ValueOf(p)
		return &val, true, nil
	} else {
		return nil, false, errors.New(fmt.Sprintf("%s undefined", name))
	}
}
