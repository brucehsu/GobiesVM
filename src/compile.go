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
	BC_INITTRANS
)

type Instruction struct {
	inst_seq  int
	inst_type Bytecode
	obj       Object
	argc      int
}

func (vm *GobiesVM) AddInstruction(inst Bytecode, obj Object) {
	new_inst := Instruction{inst_type: inst, inst_seq: len(vm.instList)}
	if obj != nil {
		new_inst.obj = obj
	}
	vm.instList = append(vm.instList, new_inst)
}

func (VM *GobiesVM) compile(node *AST) {
	for node != nil {
		switch node.Type {
		case NODE_ROOT:
			head := node.args[0]
			VM.compile(head)
		case NODE_LIST:
		case NODE_BLOCK:
			body := node.args[0]
			params := node.args[1]

			inst_last_idx := len(VM.instList) - 1
			VM.compile(body)
			block := newRBlock(VM.instList[inst_last_idx+1:])
			VM.instList = VM.instList[:inst_last_idx+1]
			VM.AddInstruction(BC_PUTOBJ, block)

			if params != nil { // Block has params
				argc := params.length
				param_array := RArray_new(VM, nil, nil)
				for i := 0; i < argc; i++ {
					dummy_arg := []Object{params.args[0].value.str}
					param := RString_new(VM, nil, dummy_arg)
					dummy_arg[0] = param
					RArray_append(VM, param_array, dummy_arg)
					params = params.next
				}
				block.ivars["params"] = param_array
			} else {
				block.ivars["params"] = nil
			}
		case NODE_VALUE:
			// NODE_VALUE can only be either NUMBER or SYMBOL
			astval := node.args[0]
			if len(astval.value.str) != 0 { // SYMBOL
				VM.AddInstruction(BC_SETSYMBOL, astval.value.str)
			} else { // NUMBER
				val := []Object{astval.value.numeric}
				VM.AddInstruction(BC_PUTOBJ, RFixnum_new(VM, nil, val))
			}
		case NODE_ASTVAL:
			if len(node.value.str) != 0 {
				val := []Object{node.value.str}
				VM.AddInstruction(BC_PUTOBJ, RString_new(VM, nil, val))
			} else {
				val := []Object{node.value.numeric}
				VM.AddInstruction(BC_PUTOBJ, RFixnum_new(VM, nil, val))
			}
		case NODE_STRING:
			astval := node.args[0]
			val := []Object{astval.value.str}
			VM.AddInstruction(BC_PUTOBJ, RString_new(VM, nil, val))
		case NODE_ASSIGN:
			// Set local variable
			astval := node.args[0]
			VM.compile(node.args[1])
			VM.AddInstruction(BC_SETLOCAL, astval.value.str)
		case NODE_ARG:
			head := node.args[0]
			VM.compile(head)
		case NODE_SEND: // Appears in Call, AsgnCall, SpecCall and Method.
			rcv := node.args[0]
			msg := node.args[1]
			argc := 0
			block := node.args[2]

			if rcv == nil && block == nil {
				if msg.args[1] == nil { // Local variable access
					// TODO: Local method call without any argument
					VM.AddInstruction(BC_GETLOCAL, msg.args[0].value.str)
				} else { // Calling local methods with one or more arguments
					if msg.args[1] != nil {
						VM.AddInstruction(BC_PUTSELF, nil)
						VM.compile(msg.args[1])
						argc = msg.args[1].args[0].length
					}
					VM.AddInstruction(BC_SEND, msg.args[0].value.str)
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
					argc = msg.args[1].args[0].length
				}

				if block != nil {
					VM.compile(block)
					argc += 1
				}

				VM.AddInstruction(BC_SEND, msg.args[0].value.str)
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
			VM.AddInstruction(BC_GETCONST, astval.value.str)
		case NODE_SETCONST:
			name := node.args[0]
			val := node.args[1]
			VM.compile(val)
			VM.AddInstruction(BC_SETCONST, name.value.str)
		case NODE_ARRAY:
			args := node.args[0]
			argc := 0
			var head *AST
			if args != nil {
				argc = args.length
				head = args.head
			}
			VM.AddInstruction(BC_GETCONST, "RArray")
			VM.compile(head)
			VM.AddInstruction(BC_SEND, "new")
			VM.instList[len(VM.instList)-1].argc = argc
		case NODE_HASH:
			args := node.args[0]
			argc := 0
			var head *AST
			if args != nil {
				argc = args.length
				head = args.head
			}
			VM.AddInstruction(BC_GETCONST, "RHash")
			VM.compile(head)
			VM.AddInstruction(BC_SEND, "new")
			VM.instList[len(VM.instList)-1].argc = argc
		case NODE_RANGE:
		case NODE_GETIVAR:
			name := node.args[0]
			VM.AddInstruction(BC_GETIVAR, name.value.str)
		case NODE_SETIVAR:
		case NODE_GETCVAR:
			name := node.args[0]
			VM.AddInstruction(BC_GETCVAR, name.value.str)
		case NODE_SETCVAR:
		case NODE_GETGLOBAL:
		case NODE_SETGLOBAL:
		case NODE_ADD:
			rcv := node.args[0]
			args := node.args[1]
			VM.compile(rcv)
			VM.compile(args)
			VM.AddInstruction(BC_SEND, "+")
			VM.instList[len(VM.instList)-1].argc = 1
		case NODE_SUB:
			rcv := node.args[0]
			args := node.args[1]
			VM.compile(rcv)
			VM.compile(args)
			VM.AddInstruction(BC_SEND, "-")
			VM.instList[len(VM.instList)-1].argc = 1
		case NODE_LT:
		case NODE_NEG:
		case NODE_NOT:
		}
		node = node.next
	}
}

func findTransactions(instList []Instruction) []Instruction {
FIND_TRANSACTION_BEGIN:
	for i, node := range instList {
		if node.inst_type == BC_GETCONST && node.obj.(string) == "Thread" {
			// Assume Thread always comes with .new(&block)
			block := instList[i+1].obj.(*RObject)
			node.inst_type = BC_INITTRANS
			node.obj = block
			instList[i] = node
			instList = append(instList[:i+1], instList[i+3:]...)
			goto FIND_TRANSACTION_BEGIN
		} else {
			if node.inst_type == BC_PUTOBJ {
				switch node.obj.(type) {
				case *RObject:
					obj := node.obj.(*RObject)
					if obj.name == "RBlock" {
						obj.methods["def"].def = findTransactions(obj.methods["def"].def)
					}
				}
			}
		}
	}
	return instList
}
