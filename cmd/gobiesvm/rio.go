package main

import (
	"fmt"
	"os"
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
	if vm.transactionEnd(env) == false {
		return nil
	}

	vm.stateLock.RLock()
	filename := v[0].(*RObject).val.str
	vm.stateLock.RUnlock()

	content, err := os.ReadFile(filename)
	if err != nil {
		vm.recordError(fmt.Errorf("IO.readlines %q: %w", filename, err))
		vm.transactionBegin(env, []Instruction{})
		return nil
	}
	str := string(content)

	vm.stateLock.RLock()
	items := []Object{}
	dummy_obj := []Object{nil}

	lines := strings.SplitAfter(str, "\n")
	for _, line := range lines {
		dummy_obj[0] = line
		rstr := RString_new(vm, nil, receiver, dummy_obj)
		items = append(items, rstr)
	}

	obj := RArray_new(vm, nil, receiver, items)
	vm.stateLock.RUnlock()

	// Begin empty transaction
	vm.transactionBegin(env, []Instruction{})

	return obj
}
