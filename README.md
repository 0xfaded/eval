eval - A library for providing an eval function in Go
============================================================================

This project adds an *eval()* function to go.

Right now only, Go expressions are handled.

Using
-----

Right now the basic dance is a 4-step process.

* Create an evaluation environment
* Parse the expression using `go/parser`
* type check the expression returned by the parser, and
* evaluate the expression.

A skeletal structure is below:

```
    package ...

	import ("reflect"; "go/parser"; "github.com/0xfaded/eval")

	...
	env := &eval.Env {
		Name:   ".",
		Vars:   make(map[string] reflect.Value),
		Consts: make(map[string] reflect.Value),
		Funcs:  make(map[string] reflect.Value),
		Types:  make(map[string] reflect.Type),
		Pkgs :  make(map[string] eval.Pkg),
	}
	// Populate env with a useful evaluation environment

    line := `5 * 6 + int32(len("abc"[0:1])))` // something to eval
	if expr, err := parser.ParseExpr(line); err != nil {
		fmt.Printf("parse error: %s\n", err)
	} else if cexpr, errs := eval.CheckExpr(expr, env); len(errs) != 0 {
		for _, cerr := range errs {
			fmt.Printf("%v\n", cerr)
		}
	} else if vals, _, err := eval.EvalExpr(cexpr, env); err != nil {
		fmt.Printf("eval error: %s\n", err)
	} else {
	  // do something with pointer to reflect.Value array vals, e.g.:
	  fmt.Println((*vals)[0].Interface())
	}
```

The program [repl.go](https://github.com/0xfaded/eval/tree/master/demo/repl.go) is a full Go program showing this.

Right now, values are retuned as a pointer to an array of
*reflect.Value()* and *reflect.Values* are used as intermediate
results. However a callback to a conversion routine is provided for
applications that need to use a different representation of a
value. The *gub* debugger is an instance where this occurs.

See Also
--------

* [What's left to do?](https://github.com/0xfaded/eval/wiki/What's-left-to-do%3F)
* [go-fish](https://github.com/rocky/go-fish): an interactive read, eval, print loop which uses this to handle the *eval()* step. In that project, see the program *make_eval* for how to create a complete environment given an initial import.
* [gub debugger](https://github.com/rocky/ssa-interp): a debugger that uses this to handle the *eval* debugger command
