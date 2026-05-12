package multityperror_test

import (
	"errors"
	"fmt"

	"github.com/azghr/forge/multityperror"
)

func Example() {
	var me multityperror.MultiError
	me.Append(fmt.Errorf("first"))
	me.Append(nil)
	me.Append(fmt.Errorf("second"))
	if !me.IsEmpty() {
		fmt.Println(me.Error())
	}
	// Output: first; second
}

func Example_empty() {
	var me multityperror.MultiError
	if me.IsEmpty() {
		fmt.Println("no errors")
	}
	// Output: no errors
}

func ExampleMultiError_Errors() {
	var me multityperror.MultiError
	me.Append(fmt.Errorf("one"))
	me.Append(fmt.Errorf("two"))

	for _, err := range me.Errors() {
		fmt.Println(err)
	}
	// Output:
	// one
	// two
}

var ErrPermission = errors.New("permission denied")

func Example_errorsIs() {
	var me multityperror.MultiError
	me.Append(fmt.Errorf("network error"))
	me.Append(ErrPermission)

	if errors.Is(&me, ErrPermission) {
		fmt.Println("found permission error")
	} else {
		fmt.Println("permission error not found")
	}
	// Output: found permission error
}
