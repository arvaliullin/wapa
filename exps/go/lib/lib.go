package main

import (
	"syscall/js"
)

func f(x float64) float64 {
	return x * x
}

func Integrate(xmin float64, xmax float64, intervals_count int) float64 {
	dx := (xmax - xmin) / float64(intervals_count)
	total := 0.0
	x := xmin
	for i := 0; i < intervals_count; i++ {
		total = total + dx*(f(x)+f(x+dx))/2.0
		x = x + dx
	}
	return total
}

func x2Integrate(this js.Value, args []js.Value) interface{} {
	xmin := args[0].Float()
	xmax := args[1].Float()
	intervals_count := args[2].Int()
	return Integrate(xmin, xmax, intervals_count)
}

func x2IntegrateMock(this js.Value, args []js.Value) interface{} {
	return 0.0
}

func FibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func fibonacciRecursive(this js.Value, args []js.Value) interface{} {
	n := args[0].Int()
	return FibonacciRecursive(n)
}

func fibonacciRecursiveMock(this js.Value, args []js.Value) interface{} {
	return 0
}

func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func fibonacciIterative(this js.Value, args []js.Value) interface{} {
	n := args[0].Int()
	return FibonacciIterative(n)
}

func fibonacciIterativeMock(this js.Value, args []js.Value) interface{} {
	return 0
}

func Multiply(size int) int {
	a, b := 33, 10
	result := 0
	for i := 0; i < size; i++ {
		result = a * b
	}
	return result
}

func multiply(this js.Value, args []js.Value) interface{} {
	size := args[0].Int()
	return Multiply(size)
}

func multiplyMock(this js.Value, args []js.Value) interface{} {
	return 0
}

func MultiplyVector(size int) []int {
	a, b := 33, 10
	aVector := make([]int, size)
	bVector := make([]int, size)
	resultVector := make([]int, size)
	for i := 0; i < size; i++ {
		aVector[i] = a
		bVector[i] = b
		resultVector[i] = aVector[i] * bVector[i]
	}
	return resultVector
}

func multiplyVector(this js.Value, args []js.Value) interface{} {
	size := args[0].Int()
	_ = MultiplyVector(size)
	return nil
}

func multiplyVectorMock(this js.Value, args []js.Value) interface{} {
	return nil
}

func Factorize(n int) []int {
	factors := make([]int, 0)
	d := 2
	for n > 1 {
		for n%d == 0 {
			factors = append(factors, d)
			n /= d
		}
		d++
	}
	return factors
}

func factorize(this js.Value, args []js.Value) interface{} {
	n := args[0].Int()
	_ = Factorize(n)
	return nil
}

func factorizeMock(this js.Value, args []js.Value) interface{} {
	return nil
}

func main() {
	js.Global().Set("x2Integrate", js.FuncOf(x2Integrate))
	js.Global().Set("x2IntegrateMock", js.FuncOf(x2IntegrateMock))

	js.Global().Set("fibonacciRecursive", js.FuncOf(fibonacciRecursive))
	js.Global().Set("fibonacciRecursiveMock", js.FuncOf(fibonacciRecursiveMock))

	js.Global().Set("fibonacciIterative", js.FuncOf(fibonacciIterative))
	js.Global().Set("fibonacciIterativeMock", js.FuncOf(fibonacciIterativeMock))

	js.Global().Set("multiply", js.FuncOf(multiply))
	js.Global().Set("multiplyMock", js.FuncOf(multiplyMock))

	js.Global().Set("multiplyVector", js.FuncOf(multiplyVector))
	js.Global().Set("multiplyVectorMock", js.FuncOf(multiplyVectorMock))

	js.Global().Set("factorize", js.FuncOf(factorize))
	js.Global().Set("factorizeMock", js.FuncOf(factorizeMock))

	select {}
}
