package eval

import (
	"errors"
	"reflect"

	"go/ast"
	"go/token"
)

// Place holder for something more substantial
func CheckStmt(stmt ast.Stmt, env Env) (Stmt, []error) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		var errs, moreErrs []error
		a := &AssignStmt{
			AssignStmt: s,
			Lhs: make([]Expr, len(s.Lhs)),
			Rhs: make([]Expr, len(s.Rhs)),
		}

		names := map[int]string{}
		// Identify names
		if s.Tok == token.DEFINE {
			newName := false
			for i, l := range s.Lhs {
				if ident, ok := l.(*ast.Ident); ok {
					if ident.Name != "_" && !inTopEnv(ident.Name, env) {
						newName = true
						names[i] = ident.Name
					}
				} else {
					errs = append(errs, ErrNonNameInDeclaration{fakeCheckExpr(l, env)})
				}
			}
			if !newName {
				errs = append(errs, ErrNoNewNamesInDeclaration{a})
			}
		}

		// Check lhs
		lhsChecked := true
		for i, l := range s.Lhs {
			if isBlankIdentifier(l) {
				names[i] = "_"
				a.Lhs[i] = fakeCheckExpr(l, env)
				continue
			} else if _, ok := names[i]; ok {
				a.Lhs[i] = fakeCheckExpr(l, env)
				continue
			}
			a.Lhs[i], moreErrs = CheckExpr(l, env)
			if moreErrs != nil && !a.Lhs[i].IsConst() {
				lhsChecked = false
				errs = append(errs, moreErrs...)
				continue
			}
			// Must be addressable or map index expr
			ll := skipSuperfluousParens(a.Lhs[i])
			if index, ok := ll.(*IndexExpr); ok {
				k := index.X.(Expr).KnownType()[0].Kind()
				if k == reflect.Map || k == reflect.Slice {
					continue
				}
			} else {
				if _, err := expectSingleType(a.Lhs[i].KnownType(), a.Lhs[i]); err != nil {
					errs = append(errs, err)
				}
			}
			if !isAddressable(ll) {
				errs = append(errs, ErrCannotAssignToUnaddressable{a.Lhs[i]})
			}
		}

		isMulti := false
		var types []reflect.Type
		if len(a.Rhs) == 1 {
			a.Rhs[0], moreErrs = CheckExpr(s.Rhs[0], env)
			errs = append(errs, moreErrs...)
			if errs != nil && !a.Rhs[0].IsConst() {
				goto done
			}
			types = make([]reflect.Type, len(a.Rhs[0].KnownType()))
			copy(types, a.Rhs[0].KnownType())
			if len(a.Lhs) == 2 && len(a.Rhs) == 1 && multivalueOk(a.Rhs[0]) {
				types = append(types, boolType)
			}
			isMulti = len(a.Lhs) > 1
		} else {
			types = make([]reflect.Type, len(a.Rhs))
			for i, r := range s.Rhs {
				a.Rhs[i], moreErrs = CheckExpr(r, env)
				errs = append(errs, moreErrs...)
				if moreErrs != nil && !a.Rhs[i].IsConst() {
					continue
				}
				if t, err := expectSingleType(a.Rhs[i].KnownType(), a.Rhs[i]); err != nil {
					errs = append(errs, err)
				} else {
					types[i] = t
				}
			}
		}

		// Check rhs
		if len(a.Lhs) != len(types) {
			errs = append(errs, ErrAssignCountMismatch{a, len(a.Lhs), len(types)})
			goto done
		}
		// Check for assignability
		if !lhsChecked {
			goto done
		}

		for i := range a.Rhs {
			if types[i] == nil {
				// new variable or typecheck failed
				continue
			}
			kt := a.Lhs[i].KnownType()
			assignable := true
			if kt == nil {
				// _ or new name
				if ct, ok := types[i].(ConstType); ok {
					if ct == ConstNil {
						errs = append(errs, ErrUntypedNil{a.Rhs[i]})
						continue
					} else {
						types[i] = ct.DefaultPromotion()
					}
				}
			} else if _, ok := kt[0].(ConstType); ok {
				// Corner case for assigning to const basic lits. e.g. 1 = 1
				_, assignable = names[i]
			} else {
				// expect the left type
				types[i] = kt[0]
				var convErrs []error
				assignable, convErrs = exprAssignableTo(a.Rhs[i], types[i])
				if assignable {
					errs = append(errs, convErrs...)
				}
			}
			if !assignable {
				if isMulti {
					errs = append(errs, ErrCannotAssignToType{a.Lhs[i], a.Rhs[0], i})
				} else {
					errs = append(errs, ErrCannotAssignToType{a.Lhs[i], a.Rhs[i], -1})
				}
			}
		}
done:
		a.newNames = names
		a.types = types
		return a, errs
	default:
		return nil, []error{errors.New("Only assign statements are currently supported")}
	}
}
