package eval

import (
	"reflect"
)

type Pkg *Env

type envSource int
const (
	envUnknown envSource = iota
	envVar
	envConst
	envFunc
)

// A Environment used for evaluation
type Env struct {
	Name string  // e.g "fmt"
	Path string  // e.g. "github.com/0xfaded/eval"

	// Values
	Vars map[string] reflect.Value
	Consts map[string] reflect.Value
	Funcs map[string] reflect.Value

	// Types
	Types map[string] reflect.Type

	// Packages
	Pkgs map[string] Pkg
}
