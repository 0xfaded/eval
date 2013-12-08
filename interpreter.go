package interactive

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"go/parser"
)

func readline(prompt string, in *bufio.Reader) (string, error) {
	fmt.Printf(prompt)
	line, err := in.ReadString('\n')
	if err == nil {
		line = strings.TrimRight(line, "r\n")
	}
	return line, err
}

func Run(env *Env, results *([]interface{})) {

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
		} else if vals, _, err := EvalExpr(expr, env); err != nil {
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
				fmt.Printf("Results[%d] = \"%v\"\n", exprs, (value.Interface()))
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
