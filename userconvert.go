package interactive
import (
	"reflect"
)

/* A way to allow a user-specificed conversion routine before accessing a value.

Here's the problem. The ssa interpreter has its version of a Value which isn't a
reflect.Value. For example structures are represented as slices and other means
are needed to associated a member of a struct to the index in the slice.

It would be slow and cumbersome to try to convert ssa.interpreter
Values to reflect Values even if it could be done.

Yet, we'd like to use as much of the code here at the higher level
rather than duplicate it elsewhere. Things like index accessing and
binary operations remain roughly the same.

So instead we try lazy conversion of an ssa.intepreter Value to a
reflect Value. So in some cases, Values in composites are reflect Values of the
corresponding interpreter Values. To unbox this just before evaluation, we use
a user-conversion callback routine defined below.  */

type UserConvertFunc func(reflect.Value, bool) (reflect.Value, bool, error)
var userConversion UserConvertFunc = nil

func SetUserConversion(callback UserConvertFunc) {
	userConversion = callback
}

func GetUserConversion() UserConvertFunc {
	return userConversion
}
