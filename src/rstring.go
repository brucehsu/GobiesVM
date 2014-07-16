package main

func initRString() *RObject {
	obj := &RObject{}
	obj.name = "RString"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RString method initialization
	obj.methods["new"] = &RMethod{gofunc: RString_new}
	obj.methods["to_s"] = &RMethod{gofunc: RString_to_s}

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
	obj := receiver.(RObject)
	return obj.val.str
}
