package main

import "strings"

func initRString() *RObject {
	obj := &RObject{}
	obj.name = "RString"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RString method initialization
	obj.methods["new"] = &RMethod{gofunc: RString_new}
	obj.methods["to_s"] = &RMethod{gofunc: RString_to_s}
	obj.methods["inspect"] = &RMethod{gofunc: RString_inspect}
	obj.methods["size"] = &RMethod{gofunc: RString_length}
	obj.methods["len"] = &RMethod{gofunc: RString_length}
	obj.methods["split"] = &RMethod{gofunc: RString_split}

	return obj
}

// String.new(str='')
// v = [string]
func RString_new(vm *GobiesVM, receiver Object, v []Object) Object {
	str := ""
	if len(v) == 1 {
		str = v[0].(string)
	}

	obj := &RObject{}
	obj.class = vm.consts["RString"]
	obj.val.str = str

	return obj
}

func RString_to_s(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	return obj.val.str
}

func RString_inspect(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	str := obj.val.str
	array := []string{"'", str, "'"}
	return strings.Join(array, "")
}

func RString_length(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	arg := make([]Object, 1, 1)
	arg[0] = int64(len(obj.val.str))
	return RFixnum_new(vm, receiver, arg)
}

func RString_split(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	sep := v[0].(*RObject).val.str

	// Manually escape linebreak if any
	sep = strings.Replace(sep, "\\n", "\n", -1)

	strList := strings.Split(obj.val.str, sep)
	arg := make([]Object, len(strList), len(strList))
	for i, v := range strList {
		dummy_arg := make([]Object, 1, 1)
		dummy_arg[0] = v
		arg[i] = RString_new(vm, nil, dummy_arg)
	}
	return RArray_new(vm, nil, arg)
}
