package eval

type PanicDivideByZero struct {
}

type PanicIndexOutOfBounds struct {
}

func (err PanicDivideByZero) Error() string {
        return "runtime error: integer divide by zero"
}

func (err PanicIndexOutOfBounds) Error() string {
        return "runtime error: index out of range"
}
