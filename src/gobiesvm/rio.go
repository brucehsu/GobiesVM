package main

import (
	"io/ioutil"
	"strings"
)

func initRIO() *RObject {
	obj := &RObject{}
	obj.name = "RIO"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RIO method initialization
	obj.methods["readlines"] = &RMethod{gofunc: RIO_readlines}

	return obj
}

// IO.readlines(filename)
// [RString]
func RIO_readlines(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	vm.transactionEnd(env)

	filename := v[0].(*RObject).val.str

	content, _ := ioutil.ReadFile(filename)
	str := string(content[:])

	items := []Object{}
	dummy_obj := []Object{nil}

	lines := strings.SplitAfter(str, "\n")
	for _, line := range lines {
		dummy_obj[0] = line
		rstr := RString_new(vm, nil, receiver, dummy_obj)
		items = append(items, rstr)
	}

	obj := RArray_new(vm, nil, receiver, items)

	// Begin empty transaction
	vm.transactionBegin(env, []Instruction{})

	return obj
}
