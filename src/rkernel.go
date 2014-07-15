package main

import "fmt"

type RKernel struct {
	RObject
}

func initRKernel() *RKernel {
	rkernel := &RKernel{}
	rkernel.ivars = make(map[string]*Object)
	rkernel.methods = make(map[string]*RMethod)
	rkernel.initRKernelMethods()
	return rkernel
}

func (obj *RKernel) initRKernelMethods() {
	obj.methods["puts"] = &RMethod{gofunc: RKernel_puts}
}

func (obj *RKernel) getMethods() map[string]*RMethod {
	return obj.methods
}

func (obj *RKernel) getString() string {
	return "RKernel"
}

func RKernel_puts(vm *GobiesVM, receiver Object, v []Object) {
	for _, obj := range v {
		fmt.Println(obj.getString())
	}
}
