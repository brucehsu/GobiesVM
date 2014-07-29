package main

import "strconv"

func initRFlonum() *RObject {
	obj := &RObject{}
	obj.name = "RFlonum"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RFlonum method initialization
	obj.methods["new"] = &RMethod{gofunc: RFlonum_new}
	obj.methods["="] = &RMethod{gofunc: RFlonum_assign}
	obj.methods["+"] = &RMethod{gofunc: RFlonum_add}
	obj.methods["-"] = &RMethod{gofunc: RFlonum_sub}
	obj.methods["*"] = &RMethod{gofunc: RFlonum_mul}
	obj.methods["/"] = &RMethod{gofunc: RFlonum_div}
	obj.methods["to_s"] = &RMethod{gofunc: RFlonum_to_s}
	obj.methods["inspect"] = &RMethod{gofunc: RFlonum_to_s}

	return obj
}

// Fixnum.new(int=0)
// v = [float64]
func RFlonum_new(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	val := float64(0)
	if len(v) == 1 {
		val = v[0].(float64)
	}

	obj := &RObject{}
	obj.class = vm.consts["RFlonum"]
	obj.val.float = val
	if env == nil {
		obj.rev = vm.rev
	} else {
		obj.rev = env.transactionPC.rev
	}

	return obj
}

func RFlonum_assign(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	new_obj := env.transactionPC.objectSet[obj]

	if obj == new_obj {
		new_obj = RFlonum_new(vm, env, nil, v).(*RObject)
		env.transactionPC.objectSet[obj] = new_obj
	}

	new_obj.val.float = v[0].(*RObject).val.float

	return obj
}

func RFlonum_add(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	dummy_args := []Object{obj.val.float + v[0].(*RObject).val.float}
	obj = RFlonum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFlonum_sub(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	dummy_args := []Object{obj.val.float - v[0].(*RObject).val.float}
	obj = RFlonum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFlonum_mul(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	dummy_args := []Object{obj.val.float * v[0].(*RObject).val.float}
	obj = RFlonum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFlonum_div(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	dummy_args := []Object{obj.val.float / v[0].(*RObject).val.float}
	obj = RFlonum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFlonum_to_s(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	return RString_new(vm, env, nil, []Object{strconv.FormatFloat(obj.val.float, 'f', 12, 64)})
}
