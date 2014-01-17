package main

import (
	"bufio"
	"go/build"
	"fmt"

	"io"
	"io/ioutil"

	"strings"

	"os"
	"os/exec"
)

func compileExpr(expr string) (compileErrors []string, err error) {
	return compileExprWithDefs(expr, "")
}

func compileExprWithDefs(expr, defs string) (compileErrors []string, err error) {
	f, err := ioutil.TempFile("/tmp", "testgen")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	_, err = fmt.Fprintf(f,
`package main
func f(...interface{}) {}
func main() {
` + defs + `
	(func(...interface{}) {})(%s)
}
`, expr);

	if err != nil {
		return nil, err
	}

	// -e prints all errors
	cmd := exec.Command(build.ToolDir + "/8g", "-e", "-o", "/dev/null", f.Name())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	buf := bufio.NewReader(stdout)

	line, rerr := buf.ReadString('\n')
	for rerr == nil {
		if strings.Index(line, ": ") != -1 {
			// Remove filename prefix
			s := strings.SplitN(line, ": ", 2)[1]
			// Remove trailing \n
			s = s[:len(s)-1]
			compileErrors = append(compileErrors, s)
		}
		line, rerr = buf.ReadString('\n')
	}
	if rerr != io.EOF {
		return nil, rerr
	} else {
		return compileErrors, nil
	}
}

