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
	obj.methods["inspect"] = &RMethod{gofunc: RHash_inspect}

	return obj
}

func RHash_new(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := &RObject{}
	obj.class = vm.consts["RHash"]
	obj.ivars = make(map[string]Object)
	obj.ivars["map"] = make(map[RValue]*RObject)

	if v != nil && len(v) > 0 {
		for i := 0; i < len(v); i += 2 {
			key, val := v[i].(*RObject), v[i+1].(*RObject)
			dummy_pair := make([]Object, 2, 2)
			dummy_pair[0] = key
			dummy_pair[1] = val
			RHash_assign_to_key(vm, obj, dummy_pair)
		}
	}

	return obj
}

func RHash_find_by_key(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	hash := obj.ivars["map"].(map[RValue]*RObject)
	val, ok := hash[v[0].(*RObject).val]

	if !ok {
		return nil
	}
	return val
}

func RHash_assign_to_key(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	hash := obj.ivars["map"].(map[RValue]*RObject)
	if v != nil && len(v) == 2 {
		hash[v[0].(*RObject).val] = v[1].(*RObject)
	}

	return obj
}

func RHash_inspect(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	hash := obj.ivars["map"].(map[RValue]*RObject)

	if len(hash) == 0 {
		return "[]"
	}

	strList := make([]string, 0, 0)

	for key, val := range hash {
		valStr := make([]string, 2, 2)
		if len(key.str) != 0 {
			valStr[0] = key.str
		} else { // Currently we only have fixnum
			valStr[0] = strconv.FormatInt(key.fixnum, 10)
		}
		valStr[1] = val.methodLookup("inspect").gofunc(vm, val, nil).(string)
		strList = append(strList, strings.Join(valStr, "=>"))
	}

	finalStr := strings.Join(strList, ", ")

	strList = []string{"{", finalStr, "}"}

	return strings.Join(strList, "")
}
