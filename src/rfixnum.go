package main

import "strconv"

func initRFixnum() *RObject {
	obj := &RObject{}
	obj.name = "RFixnum"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	// RFixnum method initialization
	obj.methods["new"] = &RMethod{gofunc: RFixnum_new}
	obj.methods["="] = &RMethod{gofunc: RFixnum_assign}
	obj.methods["+"] = &RMethod{gofunc: RFixnum_add}
	obj.methods["-"] = &RMethod{gofunc: RFixnum_sub}
	obj.methods["*"] = &RMethod{gofunc: RFixnum_mul}
	obj.methods["/"] = &RMethod{gofunc: RFixnum_div}
	obj.methods["atomic_add"] = &RMethod{gofunc: RFixnum_atomic_add}
	obj.methods["to_s"] = &RMethod{gofunc: RFixnum_to_s}
	obj.methods["to_f"] = &RMethod{gofunc: RFixnum_to_f}
	obj.methods["inspect"] = &RMethod{gofunc: RFixnum_to_s}
	obj.methods["times"] = &RMethod{gofunc: RFixnum_times}

	return obj
}

// Fixnum.new(int=0)
// v = [int64]
func RFixnum_new(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	val := int64(0)
	if len(v) == 1 {
		val = v[0].(int64)
	}

	obj := &RObject{}
	obj.class = vm.consts["RFixnum"]
	obj.val.fixnum = val
	if env == nil {
		obj.rev = vm.rev
	} else {
		obj.rev = env.transactionPC.rev
	}

	return obj
}

func RFixnum_assign(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	new_obj := env.transactionPC.objectSet[obj]

	if obj == new_obj {
		new_obj = RFixnum_new(vm, env, nil, v).(*RObject)
		env.transactionPC.objectSet[obj] = new_obj
	}

	new_obj.val.fixnum = v[0].(*RObject).val.fixnum

	return obj
}

func RFixnum_add(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	dummy_args := []Object{obj.val.fixnum + v[0].(*RObject).val.fixnum}

	// Convert to RFlonum if the second operand is RFlonum
	if v[0].(*RObject).class.name == "RFlonum" {
		dummy_args[0] = float64(obj.val.fixnum) + (v[0].(*RObject).val.float)
		return RFlonum_new(vm, env, nil, dummy_args).(*RObject)
	}

	obj = RFixnum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFixnum_atomic_add(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	new_obj := env.transactionPC.objectSet[obj]

	if obj == new_obj {
		new_obj = RFixnum_new(vm, env, nil, []Object{obj.val.fixnum}).(*RObject)
		env.transactionPC.objectSet[obj] = new_obj
	}
	new_obj.val.fixnum += v[0].(*RObject).val.fixnum

	return obj
}

func RFixnum_sub(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	dummy_args := []Object{obj.val.fixnum - v[0].(*RObject).val.fixnum}

	// Convert to RFlonum if the second operand is RFlonum
	if v[0].(*RObject).class.name == "RFlonum" {
		dummy_args[0] = float64(obj.val.fixnum) - (v[0].(*RObject).val.float)
		return RFlonum_new(vm, env, nil, dummy_args).(*RObject)
	}

	obj = RFixnum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFixnum_mul(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	dummy_args := []Object{obj.val.fixnum * v[0].(*RObject).val.fixnum}

	// Convert to RFlonum if the second operand is RFlonum
	if v[0].(*RObject).class.name == "RFlonum" {
		dummy_args[0] = float64(obj.val.fixnum) * (v[0].(*RObject).val.float)
		return RFlonum_new(vm, env, nil, dummy_args).(*RObject)
	}

	obj = RFixnum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFixnum_div(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)

	// Abandon dummy_args pre-declaration cause it may trigger divide by zero
	// Convert to RFlonum if the second operand is RFlonum
	if v[0].(*RObject).class.name == "RFlonum" {
		return RFlonum_new(vm, env, nil, []Object{float64(obj.val.fixnum) / (v[0].(*RObject).val.float)}).(*RObject)
	}

	obj = RFixnum_new(vm, env, nil, []Object{obj.val.fixnum / v[0].(*RObject).val.fixnum}).(*RObject)

	return obj
}

func RFixnum_to_s(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	return RString_new(vm, env, nil, []Object{strconv.FormatInt(obj.val.fixnum, 10)})
}

func RFixnum_to_f(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	return RFlonum_new(vm, env, nil, []Object{float64(obj.val.fixnum)})
}

// RFixnum.times(&block)
// Given: [RBlock]
// Block parameters: [i]
func RFixnum_times(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	obj = env.transactionPC.objectSet[obj]
	if v != nil && len(v) == 1 { // Given a single RBlock
		block := v[0].(*RObject)

		dummy_args := []Object{0}
		params := []*RObject{nil}

		for i := int64(0); i < obj.val.fixnum; i++ {
			// Prepare block arguments
			dummy_args[0] = i
			params[0] = RFixnum_new(vm, nil, nil, dummy_args).(*RObject)

			// Let VM handle all other stuff such as clean call frame
			vm.executeBlock(env, block, params)
		}
	}
	return obj
}
