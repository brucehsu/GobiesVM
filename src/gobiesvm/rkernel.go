package main

import (
	"fmt"
	"math/rand"
	"math/big"
)

func initRKernel() *RObject {
	rkernel := &RObject{}
	rkernel.ivars = make(map[string]Object)
	rkernel.methods = make(map[string]*RMethod)
	rkernel.initRKernelMethods()

	obj := &RObject{}
	obj.class = rkernel
	obj.ivars = make(map[string]Object)
	obj.methods = make(map[string]*RMethod)
	return obj
}

func (obj *RObject) initRKernelMethods() {
	obj.methods["puts"] = &RMethod{gofunc: RKernel_puts}
	obj.methods["p"] = &RMethod{gofunc: RKernel_p}
	obj.methods["rand"] = &RMethod{gofunc: RKernel_rand}
}

// Irreversible IO functions
func RKernel_puts(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	vm.transactionEnd(env)

	for _, obj := range v {
		robj := obj.(*RObject)
		fmt.Println(robj.methodLookup("to_s").gofunc(vm, nil, robj, nil).(*RObject).val.str)
	}

	// Begin empty transaction
	vm.transactionBegin(env, []Instruction{})
	return nil
}

func RKernel_p(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	vm.transactionEnd(env)

	for _, obj := range v {
		robj := obj.(*RObject)
		fmt.Println(robj.methodLookup("inspect").gofunc(vm, nil, robj, nil).(*RObject).val.str)
	}

	// Begin empty transaction
	vm.transactionBegin(env, []Instruction{})

	return nil
}

func RKernel_rand(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	return RFixnum_new(vm, env, nil, []Object{big.NewInt(rand.Int63n(v[0].(*RObject).val.fixnum.Int64()))})
}
