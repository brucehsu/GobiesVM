package main

import "strings"

func initRArray() *RObject {
	obj := &RObject{}
	obj.name = "RArray"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RArray method initialization
	obj.methods["new"] = &RMethod{gofunc: RArray_new}
	obj.methods["<<"] = &RMethod{gofunc: RArray_append}
	obj.methods["[]"] = &RMethod{gofunc: RArray_at}
	obj.methods["[]="] = &RMethod{gofunc: RArray_assign_to_index}
	obj.methods["to_s"] = &RMethod{gofunc: RArray_to_s}
	obj.methods["inspect"] = &RMethod{gofunc: RArray_inspect}
	obj.methods["size"] = &RMethod{gofunc: RArray_length}
	obj.methods["length"] = &RMethod{gofunc: RArray_length}
	obj.methods["each"] = &RMethod{gofunc: RArray_each}

	return obj
}

func RArray_new(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := &RObject{}
	obj.class = vm.consts["RArray"]
	obj.ivars = make(map[string]Object)
	internal_array := []*RObject{}

	if v != nil && len(v) > 0 {
		for _, item := range v {
			internal_array = append(internal_array, item.(*RObject))
		}
	}

	obj.ivars["array"] = internal_array

	return obj
}

// array << obj
// [*RObject]
func RArray_append(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	internal_array := obj.ivars["array"].([]*RObject)

	obj.ivars["array"] = append(internal_array, v[0].(*RObject))

	return obj
}

func RArray_at(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	internal_array := obj.ivars["array"].([]*RObject)

	idx := v[0].(*RObject).val.fixnum

	return internal_array[idx]
}

func RArray_assign_to_index(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	internal_array := obj.ivars["array"].([]*RObject)

	idx := v[0].(*RObject).val.fixnum
	val := v[1].(*RObject)

	internal_array[idx] = val
	obj.ivars["array"] = internal_array

	return val
}

func RArray_to_s(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	internal_array := obj.ivars["array"].([]*RObject)
	strList := make([]string, 0, len(internal_array))
	for _, item := range internal_array {
		strList = append(strList, item.methodLookup("to_s").gofunc(vm, item, v).(string))
	}

	return strings.Join(strList, "\n")
}

func RArray_inspect(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	internal_array := obj.ivars["array"].([]*RObject)
	strList := make([]string, 0, len(internal_array))
	for _, item := range internal_array {
		strList = append(strList, item.methodLookup("to_s").gofunc(vm, item, v).(string))
	}

	if len(strList) == 0 {
		return "[]"
	}

	dummyList := []string{"[", strList[0]}
	strList[0] = strings.Join(dummyList, "")
	dummyList = []string{strList[len(strList)-1], "]"}
	strList[len(strList)-1] = strings.Join(dummyList, "")

	return strings.Join(strList, ", ")
}

func RArray_length(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	internal_array := obj.ivars["array"].([]*RObject)

	length := make([]Object, 1, 1)
	length[0] = int64(len(internal_array))

	return RFixnum_new(vm, receiver, length)
}

// RArray.each(&block)
// Given: [RBlock]
// Block parameters: [item]
func RArray_each(vm *GobiesVM, receiver Object, v []Object) Object {
	obj := receiver.(*RObject)
	if v != nil && len(v) == 1 { // Given a single RBlock
		block := v[0].(*RObject)
		internal_array := obj.ivars["array"].([]*RObject)

		params := make([]*RObject, 1, 1)

		for i := 0; i < len(internal_array); i++ {
			// Prepare block arguments
			params[0] = internal_array[i]

			// Let VM handle all other stuff such as clean call frame
			vm.executeBlock(block, params)
		}
	}
	return obj
}
