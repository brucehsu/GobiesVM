package main

import (
	"math/big"
)

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
	val := big.NewInt(0)
	if len(v) == 1 {
		val.Add(val, v[0].(*big.Int))
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

	new_int := big.NewInt(0)
	new_obj.val.fixnum = new_int.Add(new_int, v[0].(*RObject).val.fixnum)

	return obj
}

func RFixnum_add(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)

	// Convert to RFlonum if the second operand is RFlonum
	if v[0].(*RObject).class.name == "RFlonum" {
		new_rat := new(big.Rat)
		new_rat.SetString(obj.val.fixnum.String())
		f, _ := new_rat.Float64()
		dummy_args := []Object{f + v[0].(*RObject).val.float}
		return RFlonum_new(vm, env, nil, dummy_args).(*RObject)
	}

	dummy_args := []Object{big.NewInt(0).Add(obj.val.fixnum, v[0].(*RObject).val.fixnum)}
	obj = RFixnum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFixnum_atomic_add(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	new_obj := env.transactionPC.objectSet[obj]

	if obj == new_obj {
		new_int := big.NewInt(0)
		new_obj = RFixnum_new(vm, env, nil, []Object{new_int.Add(new_int, obj.val.fixnum)}).(*RObject)
		env.transactionPC.objectSet[obj] = new_obj
	}
	new_obj.val.fixnum.Add(new_obj.val.fixnum, v[0].(*RObject).val.fixnum)

	return obj
}

func RFixnum_sub(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)

	// Convert to RFlonum if the second operand is RFlonum
	if v[0].(*RObject).class.name == "RFlonum" {
		new_rat := new(big.Rat)
		new_rat.SetString(obj.val.fixnum.String())
		f, _ := new_rat.Float64()
		dummy_args := []Object{f - v[0].(*RObject).val.float}
		return RFlonum_new(vm, env, nil, dummy_args).(*RObject)
	}

	dummy_args := []Object{big.NewInt(0).Sub(obj.val.fixnum, v[0].(*RObject).val.fixnum)}
	obj = RFixnum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFixnum_mul(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)

	// Convert to RFlonum if the second operand is RFlonum
	if v[0].(*RObject).class.name == "RFlonum" {
		new_rat := new(big.Rat)
		new_rat.SetString(obj.val.fixnum.String())
		f, _ := new_rat.Float64()
		dummy_args := []Object{f * v[0].(*RObject).val.float}
		return RFlonum_new(vm, env, nil, dummy_args).(*RObject)
	}

	dummy_args := []Object{big.NewInt(0).Mul(obj.val.fixnum, v[0].(*RObject).val.fixnum)}
	obj = RFixnum_new(vm, env, nil, dummy_args).(*RObject)

	return obj
}

func RFixnum_div(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)

	// Abandon dummy_args pre-declaration cause it may trigger divide by zero
	// Convert to RFlonum if the second operand is RFlonum
	if v[0].(*RObject).class.name == "RFlonum" {
		new_rat := new(big.Rat)
		new_rat.SetString(obj.val.fixnum.String())
		f, _ := new_rat.Float64()
		return RFlonum_new(vm, env, nil, []Object{f / v[0].(*RObject).val.float}).(*RObject)
	}

	obj = RFixnum_new(vm, env, nil, []Object{big.NewInt(0).Div(obj.val.fixnum, v[0].(*RObject).val.fixnum)}).(*RObject)

	return obj
}

func RFixnum_to_s(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	return RString_new(vm, env, nil, []Object{obj.val.fixnum.String()})
}

func RFixnum_to_f(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	new_rat := new(big.Rat)
	new_rat.SetString(obj.val.fixnum.String())
	f, _ := new_rat.Float64()
	return RFlonum_new(vm, env, nil, []Object{f})
}

// RFixnum.times(&block)
// Given: [RBlock]
// Block parameters: [i]
func RFixnum_times(vm *GobiesVM, env *ThreadEnv, receiver Object, v []Object) Object {
	obj := addRObjectToSet(receiver.(*RObject), env)
	obj = env.transactionPC.objectSet[obj]
	if v != nil && len(v) == 1 { // Given a single RBlock
		block := v[0].(*RObject)

		params := []*RObject{nil}

		for i := big.NewInt(0); i.Cmp(obj.val.fixnum) == -1; i.Add(i, big.NewInt(1)) {
			// Prepare block arguments
			new_int := big.NewInt(0)			
			dummy_args := []Object{new_int}
			dummy_args[0] = new_int.Add(new_int, i)
			params[0] = RFixnum_new(vm, nil, nil, dummy_args).(*RObject)

			// Let VM handle all other stuff such as clean call frame
			vm.executeBlock(env, block, params)
		}
	}
	return obj
}
