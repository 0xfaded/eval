package interactive

import (
	"reflect"
	"go/ast"
)

func checkIdent(ctx *Ctx, ident *ast.Ident, env *Env) (*Ident, []error) {
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
	}

	return aexpr, nil
}
