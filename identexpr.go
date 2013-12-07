package interactive

import (
	"errors"
	"fmt"
	"reflect"

	"go/ast"
)

	//"github.com/sbinet/go-readline/pkg/readline"
func evalIdentExpr(ident *ast.Ident, env *Env) (reflect.Value, bool, error) {
	name := ident.Name
	if v, ok := env.Vars[name]; ok {
		return v.Elem(), true, nil
	} else if v, ok := env.Consts[name]; ok {
		return v, false, nil
	} else if v, ok := env.Funcs[name]; ok {
		return v, true, nil
	} else if v, ok := builtinFuncs[name]; ok {
		return v, false, nil
	} else if p, ok := env.Pkgs[name]; ok {
		return reflect.ValueOf(p), true, nil
	} else {
		return reflect.Value{}, false, errors.New(fmt.Sprintf("%s undefined", name))
	}
}
