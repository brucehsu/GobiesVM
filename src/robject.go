package main

import "sync"

type RObject struct {
	ivars     map[string]Object
	class     *RObject
	methods   map[string]*RMethod
	name      string
	val       RValue
	rev       int64
	writeLock sync.RWMutex
}

type RMethod struct {
	gofunc func(vm *GobiesVM, t *Transaction, receiver Object, v []Object) Object
	def    []Instruction
}
