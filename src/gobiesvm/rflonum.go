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
	obj.methods["=="] = &RMethod{gofunc: RFlonum_equal}
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

	num := v[0].(*RObject).val.float

	// Convert if given object is a RFixnum
	if v[0].(*RObject).class.name == "RFixnum" {
		num = float64(v[0].(*RObject).val.fixnum)
	}

	dummy_args := []Object{obj.val.float + num}
	obj = RFlonum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFlonum_sub(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)

	num := v[0].(*RObject).val.float

	// Convert if given object is a RFixnum
	if v[0].(*RObject).class.name == "RFixnum" {
		num = float64(v[0].(*RObject).val.fixnum)
	}

	dummy_args := []Object{obj.val.float - num}
	obj = RFlonum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFlonum_mul(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)

	num := v[0].(*RObject).val.float

	// Convert if given object is a RFixnum
	if v[0].(*RObject).class.name == "RFixnum" {
		num = float64(v[0].(*RObject).val.fixnum)
	}

	dummy_args := []Object{obj.val.float * num}
	obj = RFlonum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFlonum_div(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)

	num := v[0].(*RObject).val.float

	// Convert if given object is a RFixnum
	if v[0].(*RObject).class.name == "RFixnum" {
		num = float64(v[0].(*RObject).val.fixnum)
	}

	dummy_args := []Object{obj.val.float / num}
	obj = RFlonum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFlonum_equal(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	operand_obj := addRObjectToSet(v[0].(*RObject), env)

	// Convert to numeric object to compare, or else return false
	if operand_obj.class.name == "RFixnum" {
		return RBoolean_new(vm, env, nil, []Object{obj.val.float == float64(operand_obj.val.fixnum)}).(*RObject)
	} else if operand_obj.class.name == "RFlonum" {
		return RBoolean_new(vm, env, nil, []Object{obj.val.float == operand_obj.val.float}).(*RObject)
	} else {
		return RBoolean_new(vm, env, nil, []Object{false})
	}
}

func RFlonum_to_s(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	return RString_new(vm, env, nil, []Object{strconv.FormatFloat(obj.val.float, 'f', 12, 64)})
}
