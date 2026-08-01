package main

import "slices"

func initRThread() *RObject {
	obj := &RObject{}
	obj.name = "RThread"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RThread method initialization
	obj.methods["new"] = &RMethod{gofunc: RThread_new}

	return obj
}

func RThread_new(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	if vm.transactionEnd(env) == false {
		return nil
	}

	vm.stateLock.RLock()
	definition := slices.Clone(v[0].(*RObject).methods["def"].def)
	parentStack := slices.Clone(env.threadStack)
	vm.stateLock.RUnlock()

	vm.wg.Go(func() {
		vm.executeThread(definition, &ThreadEnv{threadStack: parentStack})
	})
	vm.transactionBegin(env, []Instruction{})
	return nil
}
