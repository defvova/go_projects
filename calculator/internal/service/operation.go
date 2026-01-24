package service

type Operation struct {
	Num1 int
	Num2 int
}

func (o *Operation) Add() int {
	return o.Num1 + o.Num2
}

func (o *Operation) Subtract() int {
	return o.Num1 - o.Num2
}

func (o *Operation) Multiply() int {
	return o.Num1 * o.Num2
}

func (o *Operation) Divide() int {
	return o.Num1 / o.Num2
}

func (o *Operation) Sum(n []int) int {
	result := 0
	for i := range n {
		result = result + n[i]
	}

	return result
}
