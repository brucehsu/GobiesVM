package main

import (
	"bufio"
	"os"
)

func initRIO() *RObject {
	obj := &RObject{}
	obj.name = "RIO"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RString method initialization
	obj.methods["readlines"] = &RMethod{gofunc: RIO_readlines}

	return obj
}

// IO.readlines(filename)
// [RString]
func RIO_readlines(vm *GobiesVM, receiver Object, v []Object) Object {
	filename := v[0].(*RObject).val.str

	obj := RArray_new(vm, receiver, nil)

	input, _ := os.Open(filename)
	defer input.Close()

	reader := bufio.NewReader(input)

	line, err := reader.ReadString('\n')
	for ; err == nil; line, err = reader.ReadString('\n') {
		dummy_obj := make([]Object, 1, 1)
		dummy_obj[0] = line
		rstr := RString_new(vm, receiver, dummy_obj)
		dummy_obj[0] = rstr
		RArray_append(vm, obj, dummy_obj)
	}

	return obj
}
