package main

import "github.com/kr/try"

type RObject struct {
	ivars     map[string]Object
	class     *RObject
	methods   map[string]*RMethod
	name      string
	val       RValue
	rev       int64
	writeLock try.Mutex
	blockVar  bool
}

type RMethod struct {
	gofunc func(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object
	def    []Instruction
}
