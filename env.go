package interactive

import (
	"reflect"
)

type Pkg *Env

type Env struct {
	Name string

	// Values
	Vars map[string] reflect.Value
	Consts map[string] reflect.Value
	Funcs map[string] reflect.Value

	// Types
	Types map[string] reflect.Type

	// Packages
	Pkgs map[string] Pkg
}

