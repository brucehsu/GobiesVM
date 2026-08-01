package main

type RObject struct {
	ivars    map[string]Object
	class    *RObject
	methods  map[string]*RMethod
	name     string
	val      RValue
	rev      int64
	blockVar bool
}

type RMethod struct {
	gofunc func(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object
	def    []Instruction
}
