package eval

type PanicDivideByZero struct {}
type PanicInvalidDereference struct {}
type PanicIndexOutOfBounds struct {}
type PanicSliceOutOfBounds struct {}

func (err PanicDivideByZero) Error() string {
        return "runtime error: integer divide by zero"
}

func (err PanicInvalidDereference) Error() string {
	return "runtime error: invalid memory address or nil pointer dereference"
}

func (err PanicIndexOutOfBounds) Error() string {
        return "runtime error: index out of range"
}

func (err PanicSliceOutOfBounds) Error() string {
        return "runtime error: slice bounds out of range"
}
