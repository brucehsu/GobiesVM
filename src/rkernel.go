package main

import "fmt"

type RKernel struct {
	RObject
}

func initRKernel() *RKernel {
	rkernel := &RKernel{}
	rkernel.ivars = make(map[string]*Object, 256)
	rkernel.methods = make(map[string]*RMethod, 256)
	rkernel.initRKernelMethods()
	return rkernel
}

func (obj *RKernel) initRKernelMethods() {
	obj.methods["puts"] = &RMethod{gofunc: RKernel_puts}
}

func (obj *RKernel) getMethods() map[string]*RMethod {
	return obj.methods
}

func RKernel_puts(vm *GobiesVM, receiver Object, v ...interface{}) {
	str := v[0]
	fmt.Println(str)
}
