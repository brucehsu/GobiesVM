package main

type CallFrame struct {
	parent    *CallFrame
	var_table map[string]*Object
	stack     []*Object
	me        Object
}

type GobiesVM struct {
	instList  []*Instruction
	callStack []*CallFrame
	consts    map[string]*Object
	symbols   map[string]int
}

func initVM() *GobiesVM {
	VM := &GobiesVM{}
	top := &CallFrame{}
	VM.callStack = append(VM.callStack, top)
	top.me = initRKernel()
	methods := top.me.getMethods()
	methods["puts"].gofunc(VM, top.me, "Hello, world!")
	return VM
}
