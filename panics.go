package eval

type PanicDivideByZero struct {
}

func (err PanicDivideByZero) Error() string {
        return "runtime error: integer divide by zero"
}
