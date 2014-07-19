package main

func newRBlock(inst []Instruction) *RObject {
	obj := &RObject{}
	obj.name = "RBlock"
	obj.ivars = make(map[string]Object)
	obj.class = nil
	obj.methods = make(map[string]*RMethod)

	var new_inst []Instruction
	for _, v := range inst {
		new_inst = append(new_inst, v)
	}

	obj.methods["def"] = &RMethod{def: new_inst}

	return obj
}
