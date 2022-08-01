package interpreter

type FunctionType uint8

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
	FunctionTypeMethod
	FunctionTypeIinitalizer
)
