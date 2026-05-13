package mathutil_test

import (
	"fmt"

	"github.com/azghr/forge/mathutil"
)

func ExampleClamp() {
	fmt.Println(mathutil.Clamp(10, 0, 5))
	fmt.Println(mathutil.Clamp(2, 0, 5))
	fmt.Println(mathutil.Clamp(-1, 0, 5))
	// Output:
	// 5
	// 2
	// 0
}

func ExampleSign() {
	fmt.Println(mathutil.Sign(3.0))
	fmt.Println(mathutil.Sign(-0.1))
	fmt.Println(mathutil.Sign(0.0))
	// Output:
	// 1
	// -1
	// 0
}

func ExampleLerp() {
	fmt.Println(mathutil.Lerp(0, 10, 0.5))
	fmt.Println(mathutil.Lerp(0, 10, 0))
	fmt.Println(mathutil.Lerp(0, 10, 1))
	// Output:
	// 5
	// 0
	// 10
}

func ExampleGCD() {
	fmt.Println(mathutil.GCD(8, 12))
	fmt.Println(mathutil.GCD(0, 5))
	fmt.Println(mathutil.GCD(7, 13))
	// Output:
	// 4
	// 5
	// 1
}

func ExampleApproxEqual() {
	fmt.Println(mathutil.ApproxEqual(1.0, 1.0+1e-10))
	fmt.Println(mathutil.ApproxEqual(1.0, 2.0))
	// Output:
	// true
	// false
}

func ExampleApproxEqualEpsilon() {
	fmt.Println(mathutil.ApproxEqualEpsilon(1.0, 1.5, 0.6))
	fmt.Println(mathutil.ApproxEqualEpsilon(1.0, 2.0, 0.6))
	// Output:
	// true
	// false
}
