// Demos replacing the default identifier lookup value mechanism with
// our own custom version.

package interactive

import (
	"testing"
	"reflect"
)

// Here's our custom ident lookup.
func MyEvalIdentExpr(ctx *Ctx, ident *Ident, env *Env) (
	*reflect.Value, bool, error) {
	name := ident.Name
	if name == "nil" {
		return nil, false, nil
	} else if name[0] == 'v' {
		val := reflect.ValueOf(5)
		return &val, true, nil
	} else if name[0] == 'c' {
		val := reflect.ValueOf("constant")
		return &val, true, nil
	} else if name[0] == 'c' {
		val := reflect.ValueOf(true)
		return &val, true, nil
	} else {
		val := reflect.ValueOf('x')
		return &val, true, nil
	}
}


func TestReplaceIdentLookup(t *testing.T) {
	defer SetEvalIdentExprCallback(EvalIdentExpr)
	env := makeEnv()
	SetEvalIdentExprCallback(MyEvalIdentExpr)
	expectResult(t, "fdafdsa", env, 'x')
	expectResult(t, "c + \" value\"", env, "constant value")

}
