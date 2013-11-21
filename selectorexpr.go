package interactive

import (
	"errors"
	"fmt"
	"reflect"

	"go/ast"
)

func evalSelectorExpr(selector *ast.SelectorExpr, env *Env) (reflect.Value, bool, error) {
       var err error
       var x []reflect.Value
       if x, _, err = evalExpr(selector.X, env); err != nil {
               return reflect.Value{}, true, err
       }
       sel := selector.Sel.Name
       xname := x[0].Type().Name()

       if x[0].Kind() == reflect.Ptr {
               // Special case for handling packages
               if x[0].Type() == reflect.TypeOf(Pkg(nil)) {
                       return evalIdentExpr(selector.Sel, x[0].Interface().(Pkg))
               } else if !x[0].IsNil() && x[0].Elem().Kind() == reflect.Struct {
                       x[0] = x[0].Elem()
               }
       }

       switch x[0].Type().Kind() {
       case reflect.Struct:
               if v := x[0].FieldByName(sel); v.IsValid() {
                       return v, true, nil
               } else if x[0].CanAddr() {
                       if v := x[0].Addr().MethodByName(sel); v.IsValid() {
                               return v, true, nil
                       }
               }
               return reflect.Value{}, true, errors.New(fmt.Sprintf("%s has no field or method %s", xname, sel))
       case reflect.Interface:
               if v := x[0].MethodByName(sel); !v.IsValid() {
                       return v, true, errors.New(fmt.Sprintf("%s has no method %s", xname, sel))
               } else {
                       return v, true, nil
               }
       default:
               err = errors.New(fmt.Sprintf("%s.%s undefined (%s has no field or method %s)",
                       xname, sel, xname, sel))
               return reflect.Value{}, true, err
       }
}
