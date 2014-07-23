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
	obj.methods["p"] = &RMethod{gofunc: RKernel_p}
}

func RKernel_puts(vm *GobiesVM, t *Transaction, receiver Object, v []Object) Object {
	for _, obj := range v {
		robj := obj.(*RObject)
		fmt.Println(robj.methodLookup("to_s").gofunc(vm, t, robj, nil).(*RObject).val.str)
	}
	return nil
}

func RKernel_p(vm *GobiesVM, t *Transaction, receiver Object, v []Object) Object {
	for _, obj := range v {
		robj := obj.(*RObject)
		fmt.Println(robj.methodLookup("inspect").gofunc(vm, t, robj, nil).(*RObject).val.str)
	}
	return nil
}
