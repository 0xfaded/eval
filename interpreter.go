package interactive

import (
	"os"
	"fmt"

	"reflect"
	"sort"
	"strings"

	"go/parser"

	rl "code.google.com/p/go-gnureadline"
)

func Run(env *Env) {
	var line string
	term := os.ExpandEnv("TERM")

	names := completions(env)

	rl.SetAttemptedCompletionFunction(func(text string, start, end int) []string {
		if text == "" {
			return nil
		}

		n := len(names)
		top := sort.Search(n, func(i int) bool { return text <= names[i] })
		bot := sort.Search(n, func(i int) bool { return i >= top && !strings.HasPrefix(names[i], text) })

		if bot == top {
			return nil
		} else if bot - top == 1 {
			return names[top:bot]
		} else {
			ntop := names[top]
			nbot := names[bot-1]
			n := len(ntop)
			if len(nbot) < n {
				n = len(nbot)
			}
			var i int
			for i = 0; i < n && ntop[i] == nbot[i]; i += 1 {}

			return append([]string{ntop[0:i]}, names[top:bot]...)
		}
	})

	line, rlerr := rl.Readline("go> ")
	for rlerr == nil && line != "quit" {

		//line = "func() {" + line + "}"
		if expr, err := parser.ParseExpr(line); err != nil {
			fmt.Printf("%s\n", err)
		} else if vals, _, err := evalExpr(expr, env); err != nil {
			fmt.Printf("%s\n", err)
		} else if len(vals) == 0 {
			fmt.Printf("void")
		} else if len(vals) == 1 {
			fmt.Printf("%v\n", vals[0].Interface())
		} else {
			sep := "("
			for _, v := range vals {
				fmt.Printf("%s%v", sep, v.Interface())
			}
			fmt.Printf(")\n")
		}

		line, rlerr = rl.Readline("go> ")

	}
	//WriteHistory("data/deleteme.history")
	rl.Rl_reset_terminal(term)
}

func completions (env *Env) (names []string) {
	prefix := env.Name
	if env.Name == "." {
		for k := range builtinTypes {
			names = append(names, k)
		}
		for k := range builtinFuncs {
			names = append(names, k)
		}
		for _, v := range env.Pkgs {
			names = append(names, completions(v)...)
		}
		prefix = ""
	} else {
		prefix += "."
	}

	for k := range env.Vars {
		names = append(names, prefix + k)
	}

	for k := range env.Consts {
		names = append(names, prefix + k)
	}

	for k := range env.Funcs {
		names = append(names, k)
	}

	for k, v := range env.Types {
		names = append(names, prefix + k)
		if v.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i += 1 {
				names = append(names, v.Field(i).Name)
			}
		}
		if v.Kind() == reflect.Struct || v.Kind() == reflect.Interface {
			for i := 0; i < v.NumMethod(); i += 1 {
				names = append(names, v.Method(i).Name)
			}
		}
	}
	sort.Strings(names)
	return names
}
