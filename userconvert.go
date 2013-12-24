package interactive
import (
	"reflect"
)

/* A way to allow a user-specificed conversion routine before accessing a value.

Here's the problem. The ssa interpreter has its version of a Value which isn't a
reflect.Value. For example structures are represented as slices and other means
are needed to associated a member of a struct to the index in the slice.

It would be slow and cumbersome to try to convert all ssa.interpreter
Values to reflect Values even if it could be done.

Yet, we'd like to use as much of the code here at the higher level and
not duplicate it elsewhere. Things like index accessing and binary
operations remain roughly the same even if the representation of more
atomic values varies.

So instead, we do lazy conversion of an ssa.intepreter Value to a
reflect Value. I some cases, in gub Values in composites are a shallow
reflect conversion of the composite while the elements inside the
composite remove interpreter Values. To unbox this just before
evaluation, we use a user-conversion callback routine defined below.
This is called at various points in evaluation such as getting the
operands of a binary operation on a scalar value.

*/

type UserConvertFunc func(reflect.Value, bool) (reflect.Value, bool, error)
var userConversion UserConvertFunc = nil

func SetUserConversion(callback UserConvertFunc) {
	userConversion = callback
}

func GetUserConversion() UserConvertFunc {
	return userConversion
}
