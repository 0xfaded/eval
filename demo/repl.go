// Example repl is a simple REPL (read-eval-print loop) for GO using
// http://github.com/0xfaded/eval to the heavy lifting to implement
// the eval() part.
//
// The intent here is to show how more to use the library, rather than
// be a full-featured REPL.
//
// A more complete REPL including command history, tab completion and
// readline editing is available as a separate package:
// http://github.com/rocky/go-fish
//
// (rocky) My intent here is also to have something that I can debug in
// the ssa-debugger tortoise/gub.sh. Right now that can't handle the
// unsafe package, pointers, and calls to C code. So that let's out
// go-gnureadline and lineedit.
package main

import (
	"bufio"
	"fmt"
	"go/parser"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/0xfaded/eval"
)

// Simple replacement for GNU readline
func readline(prompt string, in *bufio.Reader) (string, error) {
	fmt.Printf(prompt)
	line, err := in.ReadString('\n')
	if err == nil {
		line = strings.TrimRight(line, "\r\n")
	}
	return line, err
}

func intro_text() {
	fmt.Printf(`=== A simple Go eval REPL ===

Results of expression are stored in variable slice "results".
The environment is stored in global variable "env".

Enter expressions to be evaluated at the "go>" prompt.

To see all results, type: "results".

To quit, enter: "quit" or Ctrl-D (EOF).
`)

}

// REPL is the a read, eval, and print loop.
func REPL(env *eval.Env, results *([]interface{})) {

	var err error
	exprs := 0
	in := bufio.NewReader(os.Stdin)
	line, err := readline("go> ", in)
	for line != "quit" {
		if err != nil {
			if err == io.EOF { break }
			panic(err)
		}
		ctx := &eval.Ctx{line}
		if expr, err := parser.ParseExpr(line); err != nil {
			if pair := eval.FormatErrorPos(line, err.Error()); len(pair) == 2 {
				fmt.Println(pair[0])
				fmt.Println(pair[1])
			}
			fmt.Printf("parse error: %s\n", err)
		} else if cexpr, errs := eval.CheckExpr(ctx, expr, env); len(errs) != 0 {
			for _, cerr := range errs {
				fmt.Printf("%v\n", cerr)
			}
		} else if vals, _, err := eval.EvalExpr(ctx, cexpr, env); err != nil {
			fmt.Printf("eval error: %s\n", err)
		} else if vals == nil {
			fmt.Printf("Kind=nil\nnil\n")
		} else if len(*vals) == 0 {
			fmt.Printf("Kind=Slice\nvoid\n")
		} else if len(*vals) == 1 {
			value := (*vals)[0]
			if value.IsValid() {
				kind := value.Kind().String()
				typ  := value.Type().String()
				if typ != kind {
					fmt.Printf("Kind = %v\n", kind)
					fmt.Printf("Type = %v\n", typ)
				} else {
					fmt.Printf("Kind = Type = %v\n", kind)
				}
				fmt.Printf("results[%d] = %s\n", exprs, eval.Inspect(value))
				exprs += 1
				*results = append(*results, (*vals)[0].Interface())
			} else {
				fmt.Printf("%s\n", value)
			}
		} else {
			fmt.Printf("Kind = Multi-Value\n")
			size := len(*vals)
			for i, v := range *vals {
				fmt.Printf("%s", eval.Inspect(v))
				if i < size-1 { fmt.Printf(", ") }
			}
			fmt.Printf("\n")
			exprs += 1
			*results = append(*results, (*vals))
		}

		line, err = readline("go> ", in)
	}
}

var results []interface{} = make([] interface{}, 0, 10)

// Create an eval.Env environment to use in evaluation.
// This is a bit ugly here, because we are rolling everything by hand, but
// we want some sort of environment to show off in demo'ing.
// The artifical environment we create here consists of
//   fmt:
//      fns: fmt.Println, fmt.Printf
//   os:
//      types: MyInt
//      vars: Stdout, Args
//   main:
//      type Alice
//      var  results, alice, aliceptr
//
// See make_env in github.com/rocky/go-fish for an automated way to
// create more complete environment from a starting import.
func makeBogusEnv() eval.Env {
	var vars   map[string] reflect.Value = make(map[string] reflect.Value)
	var consts map[string] reflect.Value = make(map[string] reflect.Value)
	var types  map[string] reflect.Type  = make(map[string] reflect.Type)

	var global_funcs map[string] reflect.Value = make(map[string] reflect.Value)
	var global_vars map[string]  reflect.Value = make(map[string] reflect.Value)

	// A place to store result values of expressions entered
	// interactively
	global_vars["results"] = reflect.ValueOf(&results)

	// What we have from the fmt package.
	var fmt_funcs    map[string] reflect.Value = make(map[string] reflect.Value)
	fmt_funcs["Println"] = reflect.ValueOf(fmt.Println)
	fmt_funcs["Printf"] = reflect.ValueOf(fmt.Printf)

	// Some "alice" things for testing
	type Alice struct {
		Bob int
		Secret string
	}

	var alice = Alice{1, "shhh"}
	alicePtr := &alice
	global_vars["alice"]    = reflect.ValueOf(alice)
	global_vars["alicePtr"] = reflect.ValueOf(alicePtr)

	// And a simple type
	type MyInt int

	// A stripped down package environment.  See
	// http://github.com/rocky/go-fish and repl_imports.go for a more
	// complete environment.
	pkgs := map[string] eval.Pkg {
			"fmt": &eval.Env {
				Name:   "fmt",
				Path:   "fmt",
				Vars:   vars,
				Consts: consts,
				Funcs:  fmt_funcs,
				Types:  types,
				Pkgs:   make(map[string] eval.Pkg),
			}, "os": &eval.Env {
				Name:   "os",
				Path:   "os",
				Vars:   map[string] reflect.Value {
					"Stdout": reflect.ValueOf(&os.Stdout),
					"Args"  : reflect.ValueOf(&os.Args)},
				Consts: make(map[string] reflect.Value),
				Funcs:  make(map[string] reflect.Value),
				Types:  map[string] reflect.Type{
					"MyInt": reflect.TypeOf(*new(MyInt))},
				Pkgs:   make(map[string] eval.Pkg),
			},
		}

	env := eval.Env {
		Name:   ".",
		Path:   "",
		Vars:   global_vars,
		Consts: make(map[string] reflect.Value),
		Funcs:  global_funcs,
		Types:  map[string] reflect.Type{ "Alice": reflect.TypeOf(Alice{}) },
		Pkgs:   pkgs,
	}
	return env
}

func main() {
	env := makeBogusEnv()
	intro_text()
	REPL(&env, &results)
}
