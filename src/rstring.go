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
	obj.methods["+"] = &RMethod{gofunc: RString_concat}
	obj.methods["to_s"] = &RMethod{gofunc: RString_to_s}
	obj.methods["inspect"] = &RMethod{gofunc: RString_inspect}
	obj.methods["size"] = &RMethod{gofunc: RString_length}
	obj.methods["len"] = &RMethod{gofunc: RString_length}
	obj.methods["split"] = &RMethod{gofunc: RString_split}

	return obj
}

// String.new(str='')
// v = [string]
func RString_new(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	str := ""
	if len(v) == 1 {
		str = v[0].(string)
	}

	obj := &RObject{}
	obj.class = vm.consts["RString"]
	obj.val.str = str
	if env == nil {
		obj.rev = vm.rev
	} else {
		obj.rev = env.transactionPC.rev
	}

	return obj
}

func RString_concat(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	obj = env.transactionPC.objectSet[obj]
	substr := v[0].(*RObject).val.str

	return RString_new(vm, env, nil, []Object{obj.val.str + substr})
}

func RString_to_s(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	return obj
}

func RString_inspect(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	str := obj.val.str
	return RString_new(vm, env, nil, []Object{"'" + str + "'"})
}

func RString_length(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	arg := make([]Object, 1, 1)
	arg[0] = int64(len(obj.val.str))
	return RFixnum_new(vm, env, receiver, arg)
}

func RString_split(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	sep := v[0].(*RObject).val.str

	strList := []string{}

	if sep == " " {
		strList = strings.Fields(obj.val.str)
	} else {
		strList = strings.Split(obj.val.str, sep)
	}

	arg := make([]Object, len(strList))
	dummy_arg := []Object{nil}
	for i, v := range strList {
		dummy_arg[0] = v
		arg[i] = RString_new(vm, env, nil, dummy_arg)
	}
	return RArray_new(vm, env, nil, arg)
}
