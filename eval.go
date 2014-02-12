package eval

import (
	"reflect"

	"go/parser"
	"go/scanner"
)

var emptyEnv Env = MakeSimpleEnv()

func Eval(expr string) (result []reflect.Value, panik error, compileErrors []error) {
	return EvalEnv(expr, emptyEnv)
}

func EvalEnv(expr string, env Env) (result []reflect.Value, panik error, compileErrors []error) {
	if e, err := parser.ParseExpr(expr); err != nil {
		errs := err.(scanner.ErrorList)
		for i := range errs {
			compileErrors = append(compileErrors, errs[i])
		}
	} else if cexpr, errs := CheckExpr(e, env); errs != nil {
		compileErrors = errs
	} else {
		result, panik = EvalExpr(cexpr, env)
	}
	return
}
