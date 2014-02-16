package eval

import (
	"fmt"
)

func InterpStmt(stmt Stmt, env Env) error {
	switch s := stmt.(type) {
	case *AssignStmt:
		if len(s.Rhs) == 1 {
			rs, err := evalTypedExpr(s.Rhs[0], s.types, env)
			if err != nil {
				if _, ok := err.(PanicInterfaceConversion); !ok || len(s.types) != 2 {
					return err
				}
			}
			for i, lhs := range s.Lhs {
				if name, ok := s.newNames[i]; !ok {
					if l, err := evalTypedExpr(lhs, lhs.KnownType(), env); err != nil {
						return err
					} else {
						l[i].Set(rs[i])
					}
				} else if name != "_" {
					v := hackedNew(s.types[i])
					v.Elem().Set(rs[i])
					env.AddVar(name, v)
				}
			}
		} else {
			for i, lhs := range s.Lhs {
				r, err := evalTypedExpr(s.Rhs[i], s.types[i:i+1], env)
				if err != nil {
					return err
				}
				if name, ok := s.newNames[i]; !ok {
					if l, err := evalTypedExpr(lhs, lhs.KnownType(), env); err != nil {
						return err
					} else {
						l[0].Set(r[0])
					}
				} else if name != "_" {
					v := hackedNew(s.types[i])
					v.Elem().Set(r[0])
					env.AddVar(name, v)
				}
			}
		}
		return nil
	default:
		panic(dytc(fmt.Sprintf("Unsupported statement %T", s)))
	}
}

