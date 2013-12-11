package main

import (
	"bufio"
	"fmt"
	"go/parser"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/0xfaded/go-interactive"
)

func readline(prompt string, in *bufio.Reader) (string, error) {
	fmt.Printf(prompt)
	line, err := in.ReadString('\n')
	if err == nil {
		line = strings.TrimRight(line, "\r\n")
	}
	return line, err
}

func run(env *interactive.Env, results *([]interface{})) {

	var err error
	exprs := 0
	in := bufio.NewReader(os.Stdin)
	line, err := readline("go> ", in)
	for line != "quit" {
		if err != nil {
			if err == io.EOF { break }
			panic(err)
		}
		//line = "func() {" + line + "}"
		if expr, err := parser.ParseExpr(line); err != nil {
			fmt.Printf("parse error: %s\n", err)
		} else if vals, _, err := interactive.EvalExpr(&interactive.Ctx{line}, expr, env); err != nil {
			fmt.Printf("eval error: %s\n", err)
		} else if vals == nil {
			fmt.Printf("nil\n")
		} else if len(*vals) == 0 {
			fmt.Printf("void\n")
		} else if len(*vals) == 1 {
			value := (*vals)[0]
			kind := value.Kind().String()
			fmt.Printf("Kind = %v\n", kind)
			if kind == "string" {
				fmt.Printf("Results[%d] = %s\n", exprs,
					strconv.QuoteToASCII(value.String()))
			} else {
				fmt.Printf("Results[%d] = %v\n", exprs, (value.Interface()))
			}
			exprs  += 1
			*results = append(*results, (*vals)[0].Interface())
		} else {
			sep := "("
			for _, v := range *vals {
				fmt.Printf("%s%v", sep, v.Interface())
			}
			fmt.Printf(")\n")
		}

		line, err = readline("go> ", in)
	}
	//WriteHistory("data/deleteme.history")
	// rl.Rl_reset_terminal(term)
}

func main() {
	var vars   map[string] reflect.Value = make(map[string] reflect.Value)
	var consts map[string] reflect.Value = make(map[string] reflect.Value)
	var funcs  map[string] reflect.Value = make(map[string] reflect.Value)
	var types  map[string] reflect.Type  = make(map[string] reflect.Type)
	{ var x log.Logger; types["Logger"] = reflect.TypeOf(x); }
	// Constants arent properly implemented, hence the conversion to int64
	consts["Ldate"] = reflect.ValueOf(int64(log.Ldate))
	consts["Ltime"] = reflect.ValueOf(int64(log.Ltime))
	consts["Lmicroseconds"] = reflect.ValueOf(int64(log.Lmicroseconds))
	consts["Llongfile"] = reflect.ValueOf(int64(log.Llongfile))
	consts["Lshortfile"] = reflect.ValueOf(int64(log.Lshortfile))
	consts["LstdFlags"] = reflect.ValueOf(int64(log.LstdFlags))
	funcs["SetOutput"] = reflect.ValueOf(log.SetOutput)
	funcs["Println"] = reflect.ValueOf(log.Println)
	funcs["Panicln"] = reflect.ValueOf(log.Panicln)
	funcs["New"] = reflect.ValueOf(log.New)
	funcs["Fatal"] = reflect.ValueOf(log.Fatal)
	funcs["Flags"] = reflect.ValueOf(log.Flags)
	funcs["SetFlags"] = reflect.ValueOf(log.SetFlags)
	funcs["Print"] = reflect.ValueOf(log.Print)
	funcs["Printf"] = reflect.ValueOf(log.Printf)
	funcs["Panic"] = reflect.ValueOf(log.Panic)
	funcs["Panicf"] = reflect.ValueOf(log.Panicf)
	funcs["Prefix"] = reflect.ValueOf(log.Prefix)
	funcs["SetPrefix"] = reflect.ValueOf(log.SetPrefix)
	funcs["Fatalf"] = reflect.ValueOf(log.Fatalf)
	funcs["Fatalln"] = reflect.ValueOf(log.Fatalln)

	type Alice struct {
		Bob int
		Secret string
	}

	var v2 []interface{} = make([] interface{}, 0, 10)
	vars["Results"] = reflect.ValueOf(&v2)
	env := interactive.Env {
		Name:   ".",
		Vars:   vars,
		Consts: make(map[string] reflect.Value),
		Funcs:  make(map[string] reflect.Value),
		Types:  map[string] reflect.Type{ "Alice": reflect.TypeOf(Alice{}) },
		Pkgs:   map[string] interactive.Pkg { "log": &interactive.Env {
				Name:   "log",
				Vars:   vars,
				Consts: consts,
				Funcs:  funcs,
				Types:  types,
				Pkgs:   make(map[string] interactive.Pkg),
			}, "os": &interactive.Env {
				Name:   "os",
				Vars:   map[string] reflect.Value { "Stdout": reflect.ValueOf(&os.Stdout) },
				Consts: make(map[string] reflect.Value),
				Funcs:  make(map[string] reflect.Value),
				Types:  make(map[string] reflect.Type),
				Pkgs:   make(map[string] interactive.Pkg),
			},
		},
	}

	run(&env, &v2)
}
