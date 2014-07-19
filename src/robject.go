package main

type RObject struct {
	ivars   map[string]Object
	class   *RObject
	methods map[string]*RMethod
	name    string
	val     RValue
}

type RMethod struct {
	gofunc func(vm *GobiesVM, receiver Object, v []Object) Object
	def    []Instruction
}
