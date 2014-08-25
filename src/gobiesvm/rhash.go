package main

import (
	"strconv"
	"strings"
)

func initRHash() *RObject {
	obj := &RObject{}
	obj.name = "RHash"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RHash method initialization
	obj.methods["new"] = &RMethod{gofunc: RHash_new}
	obj.methods["[]"] = &RMethod{gofunc: RHash_find_by_key}
	obj.methods["[]="] = &RMethod{gofunc: RHash_assign_to_key}
	obj.methods["to_s"] = &RMethod{gofunc: RHash_inspect}
	obj.methods["inspect"] = &RMethod{gofunc: RHash_inspect}

	return obj
}

func RHash_new(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := &RObject{}
	obj.class = vm.consts["RHash"]
	obj.ivars = make(map[string]Object)
	internal_map := make(map[RValue]*RObject)
	obj.ivars["map"] = internal_map
	if env == nil {
		obj.rev = vm.rev
	} else {
		obj.rev = env.transactionPC.rev
	}

	if v != nil && len(v) > 0 {
		for i, length := 0, len(v); i < length; i += 2 {
			key, val := v[i].(*RObject).val, v[i+1].(*RObject)
			internal_map[key] = val
		}
	}

	return obj
}

func RHash_find_by_key(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	obj = env.transactionPC.objectSet[obj]

	hash := obj.ivars["map"].(map[RValue]*RObject)
	val, ok := hash[v[0].(*RObject).val]

	if !ok {
		return nil
	}
	return val
}

func RHash_assign_to_key(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	new_obj := env.transactionPC.objectSet[obj]
	hash := new_obj.ivars["map"].(map[RValue]*RObject)

	if obj == new_obj {
		args := make([]Object, 0, 1)
		for key, val := range hash {
			args = append(args, &RObject{val: key})
			args = append(args, val)
		}
		new_obj = RHash_new(vm, env, nil, args).(*RObject)
		env.transactionPC.objectSet[obj] = new_obj
	}

	hash = new_obj.ivars["map"].(map[RValue]*RObject)
	if v != nil && len(v) == 2 {
		hash[v[0].(*RObject).val] = v[1].(*RObject)
	}

	return obj
}

func RHash_inspect(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	hash := obj.ivars["map"].(map[RValue]*RObject)

	if len(hash) == 0 {
		return RString_new(vm, env, nil, []Object{"[]"})
	}

	strList := make([]string, 0, 0)

	for key, val := range hash {
		valStr := make([]string, 2, 2)
		if len(key.str) != 0 {
			valStr[0] = key.str
		} else { // Currently we only have fixnum
			valStr[0] = strconv.FormatInt(key.fixnum, 10)
		}
		valStr[1] = val.methodLookup("inspect").gofunc(vm, env, val, nil).(*RObject).val.str
		strList = append(strList, strings.Join(valStr, "=>"))
	}

	finalStr := strings.Join(strList, ", ")

	return RString_new(vm, env, nil, []Object{"{" + finalStr + "}"})
}
