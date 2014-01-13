package eval

import (
	"reflect"
	"go/ast"
)

func checkIdent(ctx *Ctx, ident *ast.Ident, env *Env) (_ *Ident, errs []error) {
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
        // TODO[crc] Remove this when builtin identifiers are handled
        // in the general case. This complex method was added only to
        // reable the binaryexpr tests
        case "complex":
		aexpr.knownType = []reflect.Type{reflect.TypeOf(complex128(0))}
        default:
                if v, ok := env.Vars[aexpr.Name]; ok {
                        aexpr.knownType = knownType{v.Elem().Type()}
                } else if v, ok := env.Consts[aexpr.Name]; ok {
                        if n, ok := v.Interface().(*ConstNumber); ok {
                                aexpr.knownType = knownType{n.Type}
                        } else {
                                aexpr.knownType = knownType{v.Type()}
                        }
                        aexpr.constValue = constValue(v)
                } else if v, ok := env.Funcs[aexpr.Name]; ok {
                        aexpr.knownType = knownType{v.Type()}
                } else {
                        errs = append(errs, ErrUndefined{at(ctx, ident)})
                }
        }
	return aexpr, errs
}
