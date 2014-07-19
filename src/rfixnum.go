package main

import "strconv"

func initRFixnum() *RObject {
	obj := &RObject{}
	obj.name = "RFixnum"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RFixnum method initialization
	obj.methods["new"] = &RMethod{gofunc: RFixnum_new}
	obj.methods["to_s"] = &RMethod{gofunc: RFixnum_to_s}
	obj.methods["inspect"] = &RMethod{gofunc: RFixnum_to_s}

	return obj
}

// Fixnum.new(int=0)
// v = [int64]
func RFixnum_new(vm *GobiesVM, receiver Object, v []Object) Object {
	val := int64(0)
	if len(v) == 1 {
		val = v[0].(int64)
	}

	obj := &RObject{}
	obj.class = vm.consts["RFixnum"]
	obj.val.fixnum = val

	return obj
}

func RFixnum_to_s(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	return strconv.FormatInt(obj.val.fixnum, 10)
}
