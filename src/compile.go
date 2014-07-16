package main

type Bytecode uint8

const (
	BC_PUTSELF Bytecode = iota
	BC_PUTNIL
	BC_PUTOBJ
	BC_PUTTRUE
	BC_PUTFALSE
	BC_SETLOCAL
	BC_GETLOCAL
	BC_SETGLOBAL
	BC_GETGLOBAL
	BC_SETSYMBOL
	BC_GETSYMBOL
	BC_SETCONST
	BC_GETCONST
	BC_SETIVAR
	BC_GETIVAR
	BC_SETCVAR
	BC_GETCVAR
	BC_SEND
	BC_JUMP
)

type Instruction struct {
	inst_seq  int
	inst_type Bytecode
	obj       Object
	// block     Instruction[]
	argc int
}

func (vm *GobiesVM) AddInstruction(inst Bytecode, obj Object) {
	new_inst := &Instruction{inst_type: inst, inst_seq: len(vm.instList)}
	if obj != nil {
		new_inst.obj = obj
	}
	vm.instList = append(vm.instList, new_inst)
}

func (VM *GobiesVM) compile(node *AST) {

	switch node.Type {
	case NODE_ROOT:
		head := node.args[0]
		for head != nil {
			VM.compile(head)
			head = head.next
		}
	case NODE_LIST:
	case NODE_BLOCK:
		body := node.args[0]
		params := node.args[1]
		inst_last_idx := len(VM.instList)
		if params != nil { // Block has params
			// TODO:
		} else {

		}
		VM.compile(body)

		// block_inst_len := len(VM.instList) - inst_last_idx
		// block := VM.instList[inst_last_idx:-1]
		VM.instList = VM.instList[0:inst_last_idx]
		// VM.instList[inst_last_idx].block = block
	case NODE_VALUE:
		// NODE_VALUE can only be either NUMBER or SYMBOL
		astval := node.args[0]
		if len(astval.value.str) != 0 { // SYMBOL
			VM.AddInstruction(BC_SETSYMBOL, &RString{val: RValue{str: astval.value.str}})
		} else { // NUMBER
			VM.AddInstruction(BC_PUTOBJ, &RFixnum{val: RValue{fixnum: node.value.numeric}})
		}
	case NODE_ASTVAL:
		if len(node.value.str) != 0 {
			VM.AddInstruction(BC_PUTOBJ, &RString{val: RValue{str: node.value.str}})
		} else {
			VM.AddInstruction(BC_PUTOBJ, &RFixnum{val: RValue{fixnum: node.value.numeric}})
		}
	case NODE_STRING:
		astval := node.args[0]
		VM.AddInstruction(BC_PUTOBJ, &RString{val: RValue{str: astval.value.str}})
	case NODE_ASSIGN:
		// Set local variable
		astval := node.args[0]
		VM.compile(node.args[1])
		VM.AddInstruction(BC_SETLOCAL, &RString{val: RValue{str: astval.value.str}})
	case NODE_ARG:
		head := node.args[0]
		for {
			VM.compile(head)
			if head.next == nil {
				break
			}
			head = head.next
		}
	case NODE_SEND: // Appears in Call, AsgnCall, SpecCall and Method.
		rcv := node.args[0]
		msg := node.args[1]
		argc := 0
		block := node.args[2]

		if rcv == nil && block == nil {
			if msg.args[1] == nil { // Local variable access
				// TODO: Local method call without any argument
				VM.AddInstruction(BC_GETLOCAL, &RString{val: RValue{str: msg.args[0].value.str}})
			} else { // Calling local methods with one or more arguments
				if msg.args[1] != nil {
					VM.AddInstruction(BC_PUTSELF, nil)
					VM.compile(msg.args[1])
					argc = msg.args[1].args[0].length
				}
				VM.AddInstruction(BC_SEND, &RString{val: RValue{str: msg.args[0].value.str}})
				VM.instList[len(VM.instList)-1].argc = argc
			}
		} else {
			if rcv == nil { // SELF
				VM.AddInstruction(BC_PUTSELF, nil)
			} else {
				VM.compile(rcv)
			}

			if msg.args[1] != nil {
				VM.compile(msg.args[1])
			}

			if block != nil {
				VM.compile(block)
			}

			VM.AddInstruction(BC_SEND, &RString{val: RValue{str: msg.args[0].value.str}})
			VM.instList[len(VM.instList)-1].argc = argc
		}
	case NODE_MSG:
		VM.compile(node.args[0])
		if node.args[1] != nil { // Arguments
			VM.compile(node.args[1])
		}
	case NODE_IF:
	case NODE_UNLESS:
	case NODE_AND:
	case NODE_OR:
	case NODE_WHILE:
	case NODE_UNTIL:
	case NODE_TRUE:
		VM.AddInstruction(BC_PUTTRUE, nil)
	case NODE_FALSE:
		VM.AddInstruction(BC_PUTFALSE, nil)
	case NODE_NIL:
		VM.AddInstruction(BC_PUTNIL, nil)
	case NODE_SELF:
		VM.AddInstruction(BC_PUTSELF, nil)
	case NODE_LEAVE:
	case NODE_RETURN:
	case NODE_BREAK:
	case NODE_YIELD:
	case NODE_DEF:
	case NODE_METHOD: // Appears in method definition
	case NODE_PARAM: // Appears in method parameters and blocks
	case NODE_CLASS:
	case NODE_MODULE:
	case NODE_CONST:
		astval := node.args[0]
		VM.AddInstruction(BC_GETCONST, &RString{val: RValue{str: astval.value.str}})
	case NODE_SETCONST:
		name := node.args[0]
		val := node.args[1]
		VM.compile(val)
		VM.AddInstruction(BC_SETCONST, &RString{val: RValue{str: name.value.str}})
	case NODE_ARRAY:
	case NODE_HASH:
	case NODE_RANGE:
	case NODE_GETIVAR:
		name := node.args[0]
		VM.AddInstruction(BC_GETIVAR, &RString{val: RValue{str: name.value.str}})
	case NODE_SETIVAR:
	case NODE_GETCVAR:
		name := node.args[0]
		VM.AddInstruction(BC_GETCVAR, &RString{val: RValue{str: name.value.str}})
	case NODE_SETCVAR:
	case NODE_GETGLOBAL:
	case NODE_SETGLOBAL:
	case NODE_ADD:
	case NODE_SUB:
	case NODE_LT:
	case NODE_NEG:
	case NODE_NOT:
	}
}
