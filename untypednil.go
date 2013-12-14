package interactive

type UntypedNil struct {}
func (UntypedNil) String() string {
	return "nil"
}
