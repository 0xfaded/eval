package eval

import (
	"reflect"
	"go/ast"
)

func checkIdent(ident *ast.Ident, env *Env) (_ *Ident, errs []error) {
	aexpr := &Ident{Ident: ident}
	switch aexpr.Name {
	case "nil":
		aexpr.constValue = constValueOf(UntypedNil{})
		aexpr.knownType = []reflect.Type{ConstNil}

	case "true":
		aexpr.constValue = constValueOf(true)
		aexpr.knownType = []reflect.Type{ConstBool}

	case "false":
		aexpr.constValue = constValueOf(false)
		aexpr.knownType = []reflect.Type{ConstBool}
        default:
                if v, ok := env.Vars[aexpr.Name]; ok {
                        aexpr.knownType = knownType{v.Elem().Type()}
			aexpr.source = envVar
                } else if v, ok := env.Consts[aexpr.Name]; ok {
                        if n, ok := v.Interface().(*ConstNumber); ok {
                                aexpr.knownType = knownType{n.Type}
                        } else {
                                aexpr.knownType = knownType{v.Type()}
                        }
                        aexpr.constValue = constValue(v)
			aexpr.source = envConst
                } else if v, ok := env.Funcs[aexpr.Name]; ok {
                        aexpr.knownType = knownType{v.Type()}
			aexpr.source = envFunc
                } else {
                        errs = append(errs, ErrUndefined{aexpr})
                }
        }
	return aexpr, errs
}
