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
			for i, l := range s.Lhs {
				if ident, ok := l.(*ast.Ident); ok {
					names[i] = ident.Name
				} else {
					errs = append(errs, ErrNonNameInDeclaration{fakeCheckExpr(l, env)})
				}
			}
			if len(names) == 0 {
				errs = append(errs, ErrNoNewNamesInDeclaration{a})
			}
		}

		// Check lhs
		lhsChecked := true
		for i, l := range s.Lhs {
			if isBlankIdentifier(l) {
				a.Lhs[i] = fakeCheckExpr(l, env)
				continue
			}
			a.Lhs[i], moreErrs = CheckExpr(l, env)
			if moreErrs != nil && !a.Lhs[i].IsConst() {
				if _, ok := names[i]; !ok {
					lhsChecked = false
					errs = append(errs, moreErrs...)
				}
				continue
			} else {
				// We are only interested in new names. If this was a
				// name and the type check passed, this means the variable
				// was already in scope.
				delete(names, i)
			}
			// Must be addressable or map index expr
			ll := skipSuperfluousParens(a.Lhs[i])
			if index, ok := ll.(*IndexExpr); ok {
				k := index.X.(Expr).KnownType()[0].Kind()
				if k == reflect.Map || k == reflect.Slice {
					continue
				}
			} else {
				_, err := expectSingleType(a.Lhs[i].KnownType(), a.Lhs[i])
				errs = append(errs, err)
			}
			if isAddressable(ll) {
				errs = append(errs, ErrCannotAssignToUnaddressable{a.Lhs[i]})
			}
		}

		// Check rhs
		var types []reflect.Type
		isMulti := false
		if len(a.Rhs) == 1 {
			a.Rhs[0], moreErrs = CheckExpr(s.Rhs[0], env)
			errs = append(errs, moreErrs...)
			if errs != nil && !a.Rhs[0].IsConst() {
				goto done
			}
			types = a.Rhs[0].KnownType()
			if len(a.Lhs) == 1 {
				if t, err := expectSingleType(types, a.Rhs[0]); err != nil {
					errs = append(errs, err)
				} else if lhsChecked && names[0] != "" && !t.AssignableTo(a.Lhs[0].KnownType()[0]) {
					errs = append(errs, ErrCannotAssignToType{a.Lhs[0], a.Rhs[0], -1})
				}
				goto done
			} else if len(a.Lhs) == 2 && len(types) == 1 && multivalueOk(a.Rhs[0]) {
				types = append(types, boolType)
			}
			isMulti = true
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

		// Check for assignability
		if !lhsChecked {
			goto done
		}
		if len(a.Lhs) != len(types) {
			errs = append(errs, ErrAssignCountMismatch{a, len(a.Lhs), len(types)})
			goto done
		}
		for i, t := range types {
			if t == nil || names[i] != "" {
				// new variable or typecheck failed
				continue
			}
			if isMulti {
				if !typeAssignableTo(t, a.Lhs[i].KnownType()[0]) {
					errs = append(errs, ErrCannotAssignToType{a.Lhs[i], a.Rhs[i], i})
				}
			} else {
				assignable, convErrs := exprAssignableTo(a.Rhs[i], t)
				errs = append(errs, convErrs...)
				if !assignable {
					errs = append(errs, ErrCannotAssignToType{a.Lhs[i], a.Rhs[i], -1})
				}
			}
		}
done:
		return a, errs
	default:
		return nil, []error{errors.New("Only assign statements are currently supported")}
	}
}
