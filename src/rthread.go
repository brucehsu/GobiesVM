package main

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
	vm.transactionEnd(env)

	wg.Add(1)
	go vm.executeThread(v[0].(*RObject).methods["def"].def)

	vm.transactionBegin(env, []Instruction{})

	return nil
}
