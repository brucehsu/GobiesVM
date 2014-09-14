package main

func initRBoolean() *RObject {
	obj := &RObject{}
	obj.name = "RBoolean"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RBoolean method initialization
	obj.methods["new"] = &RMethod{gofunc: RBoolean_new}
	obj.methods["to_s"] = &RMethod{gofunc: RBoolean_to_s}

	return obj
}

// Boolean.new(bool=false)
// v = [bool]
func RBoolean_new(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	val := false
	if len(v) == 1 {
		val = v[0].(bool)
	}

	obj := &RObject{}
	obj.class = vm.consts["RBoolean"]
	obj.val.boolean = val
	if env == nil {
		obj.rev = vm.rev
	} else {
		obj.rev = env.transactionPC.rev
	}

	return obj
}

func RBoolean_to_s(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	if obj.val.boolean {
		return RString_new(vm, env, nil, []Object{"True"})
	} else {
		return RString_new(vm, env, nil, []Object{"False"})
	}
}
