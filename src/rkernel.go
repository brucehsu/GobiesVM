package main

import "fmt"

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
}

func RKernel_puts(vm *GobiesVM, receiver Object, v []Object) Object {
	for _, obj := range v {
		robj := obj.(*RObject)
		fmt.Println(robj.methodLookup("to_s").gofunc(vm, *robj, nil))
	}
	return nil
}
