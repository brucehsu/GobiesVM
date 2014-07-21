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
func RIO_readlines(vm *GobiesVM, receiver Object, v []Object) Object {
	filename := v[0].(*RObject).val.str

	content, _ := ioutil.ReadFile(filename)
	str := string(content[:])

	items := []Object{}

	lines := strings.SplitAfter(str, "\n")
	for _, line := range lines {
		dummy_obj := make([]Object, 1, 1)
		dummy_obj[0] = line
		rstr := RString_new(vm, receiver, dummy_obj)
		items = append(items, rstr)
	}

	obj := RArray_new(vm, receiver, items)

	return obj
}
