package main

type RObject struct {
	ivars   map[string]*Object
	class   *Object
	methods map[string]*RMethod
}

type RMethod struct {
	gofunc func(vm *GobiesVM, receiver Object, v ...interface{})
	def    []*Instruction
}

func (obj *RObject) getMethods() map[string]*RMethod {
	return obj.methods
}
