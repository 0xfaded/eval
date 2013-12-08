# Comments starting with #: below are remake GNU Makefile comments. See
# https://github.com/rocky/remake/wiki/Rake-tasks-for-gnu-make

.PHONY: all eval test check

#: Same as eval
all: eval

#: The front-end to the evaluator
eval: lib
	go build -o eval demo/eval.go

#: The evaluator library
lib:
	go build

#: Same as "check"
test: check

#: Run all tests (quick and interpreter)
check:
	go test -i && go test
