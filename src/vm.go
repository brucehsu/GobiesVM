package main

type CallFrame struct {
	parent    *CallFrame
	var_table map[string]Object
	stack     []Object
	me        Object
}

type GobiesVM struct {
	instList       []*Instruction
	callFrameStack []*CallFrame
	consts         map[string]*RObject
	symbols        map[string]int
}

func initVM() *GobiesVM {
	VM := &GobiesVM{}
	top := initCallFrame()
	VM.callFrameStack = append(VM.callFrameStack, top)
	VM.consts = make(map[string]*RObject)
	VM.initConsts()
	VM.symbols = make(map[string]int)
	top.me = initRKernel()
	return VM
}

func initCallFrame() *CallFrame {
	frame := &CallFrame{}
	frame.var_table = make(map[string]Object)
	return frame
}

func (VM *GobiesVM) initConsts() {
	VM.consts["RString"] = initRString()
	VM.consts["RFixnum"] = initRFixnum()
	VM.consts["RArray"] = initRArray()
	VM.consts["IO"] = initRIO()
}

func (obj *RObject) methodLookup(method_name string) *RMethod {
	if val, ok := obj.methods[method_name]; ok {
		return val
	}
	if obj.class != nil {
		return obj.class.methodLookup(method_name)
	}
	return nil
}

func (VM *GobiesVM) executeBytecode() {
	for _, v := range VM.instList {
		currentCallFrame := VM.callFrameStack[len(VM.callFrameStack)-1]
		switch v.inst_type {
		case BC_PUTSELF:
			currentCallFrame.stack = append(currentCallFrame.stack, currentCallFrame.me)
		case BC_PUTNIL:
			currentCallFrame.stack = append(currentCallFrame.stack, nil)
		case BC_PUTOBJ:
			currentCallFrame.stack = append(currentCallFrame.stack, v.obj)
		case BC_PUTTRUE:
		case BC_PUTFALSE:
		case BC_SETLOCAL:
			top := currentCallFrame.stack[len(currentCallFrame.stack)-1]
			currentCallFrame.var_table[v.obj.(string)] = top
			currentCallFrame.stack = currentCallFrame.stack[0 : len(currentCallFrame.stack)-1]
		case BC_GETLOCAL:
			currentCallFrame.stack = append(currentCallFrame.stack, currentCallFrame.var_table[v.obj.(string)])
		case BC_SETGLOBAL:
		case BC_GETGLOBAL:
		case BC_SETSYMBOL:
		case BC_GETSYMBOL:
		case BC_SETCONST:
		case BC_GETCONST:
			currentCallFrame.stack = append(currentCallFrame.stack, VM.consts[v.obj.(string)])
		case BC_SETIVAR:
		case BC_GETIVAR:
		case BC_SETCVAR:
		case BC_GETCVAR:
		case BC_SEND:
			argLists := currentCallFrame.stack[len(currentCallFrame.stack)-(v.argc+1):] // argc + 1 ensures inclusion of receiver
			currentCallFrame.stack = currentCallFrame.stack[:len(currentCallFrame.stack)-(v.argc+1)]
			recv := argLists[0].(*RObject)
			argLists = argLists[1:]
			return_val := recv.methodLookup(v.obj.(string)).gofunc(VM, recv, argLists)
			if return_val != nil {
				currentCallFrame.stack = append(currentCallFrame.stack, return_val)
			}
		case BC_JUMP:
		}
	}
}
