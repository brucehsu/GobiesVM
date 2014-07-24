package main

import (
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup

type CallFrame struct {
	parent    *CallFrame
	var_table map[string]Object
	stack     []Object
	me        Object
}

type Transaction struct {
	instList            []Instruction
	objectSet           map[*RObject]*RObject
	localCallFrameStack []*CallFrame
	inevitable          bool
	rev                 int64
}

type GobiesVM struct {
	instList       []Instruction
	callFrameStack []*CallFrame
	consts         map[string]*RObject
	symbols        map[string]int
	rev            int64
	globalLock     sync.RWMutex
	has_inevitable int32
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
	VM.consts["RHash"] = initRHash()
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

func (currentCallFrame *CallFrame) variableLookup(var_name string) Object {
	if obj, ok := currentCallFrame.var_table[var_name]; ok {
		return obj
	}
	if currentCallFrame.parent != nil {
		return currentCallFrame.parent.variableLookup(var_name)
	}
	return nil
}

func initTransaction(instList []Instruction) *Transaction {
	t := &Transaction{}
	t.instList = []Instruction{}
	for _, inst := range instList {
		t.instList = append(t.instList, inst)
	}
	return t
}

func (VM *GobiesVM) execute() {
	// 	root := initTransaction(VM.instList)
	// 	root.inevitable = true

	VM.rev = 0

	// Execute root transaction which is inevitable
	wg.Add(1)
	go VM.executeThread(VM.instList)

	wg.Wait()
}

func (VM *GobiesVM) executeThread(instList []Instruction) {
	// Create clean call frame without pushing back to VM stack
	currentCallFrame := initCallFrame()
	currentCallFrame.parent = VM.callFrameStack[len(VM.callFrameStack)-1]
	currentCallFrame.me = currentCallFrame.parent.me

	// Determine transactions and execute them
	t := initTransaction(instList)
	t.localCallFrameStack = append(t.localCallFrameStack, currentCallFrame)
	// t.inevitable = true
	VM.executeTransaction(t)
	wg.Done()
}

func (VM *GobiesVM) executeTransaction(t *Transaction) {
	// Initialize environment
	if t.inevitable {
		VM.globalLock.Lock()
		atomic.StoreInt32(&VM.has_inevitable, 1)
	} else {
		VM.globalLock.RLock()
	}
	t.rev = atomic.LoadInt64(&VM.rev)

	// Run through a speculative execution
	// Create clean call frame without pushing back to VM stack
	VM.executeBytecodes(nil, t)

	// Lock the write-set

	// Increment global revision

	// Validate the read-set

	// Commit and release the locks
	if t.inevitable {
		VM.globalLock.Unlock()
		atomic.StoreInt32(&VM.has_inevitable, 0)
	} else {
		VM.globalLock.RUnlock()
	}
}

func (VM *GobiesVM) executeBytecodes(instList []Instruction, t *Transaction) {
	if instList == nil {
		instList = t.instList
	}
	for _, v := range instList {
		// currentCallFrame := VM.callFrameStack[len(VM.callFrameStack)-1]
		currentCallFrame := t.localCallFrameStack[len(t.localCallFrameStack)-1]
		// printInstructions([]Instruction{v}, false)
		switch v.inst_type {
		case BC_PUTSELF:
			currentCallFrame.stack = append(currentCallFrame.stack, currentCallFrame.me)
		case BC_PUTNIL:
			currentCallFrame.stack = append(currentCallFrame.stack, nil)
		case BC_PUTOBJ:
			currentCallFrame.stack = append(currentCallFrame.stack, v.obj)
		case BC_PUTFIXNUM:
			currentCallFrame.stack = append(currentCallFrame.stack, RFixnum_new(VM, nil, nil, v.obj.([]Object)))
		case BC_PUTSTRING:
			currentCallFrame.stack = append(currentCallFrame.stack, RString_new(VM, nil, nil, v.obj.([]Object)))
		case BC_PUTTRUE:
		case BC_PUTFALSE:
		case BC_SETLOCAL:
			top := currentCallFrame.stack[len(currentCallFrame.stack)-1]
			currentCallFrame.var_table[v.obj.(string)] = top
			top.(*RObject).name = v.obj.(string)                                               // Change object from anonymous to named
			currentCallFrame.stack = currentCallFrame.stack[0 : len(currentCallFrame.stack)-1] // Pop object from stack
		case BC_GETLOCAL:
			currentCallFrame.stack = append(currentCallFrame.stack, currentCallFrame.variableLookup(v.obj.(string)))
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
			// fmt.Println(currentCallFrame.stack)
			argLists := currentCallFrame.stack[len(currentCallFrame.stack)-(v.argc+1):] // argc + 1 ensures inclusion of receiver
			currentCallFrame.stack = currentCallFrame.stack[:len(currentCallFrame.stack)-(v.argc+1)]
			recv := argLists[0].(*RObject)
			argLists = argLists[1:]
			return_val := recv.methodLookup(v.obj.(string)).gofunc(VM, t, recv, argLists)
			if return_val != nil {
				currentCallFrame.stack = append(currentCallFrame.stack, return_val)
			}
		case BC_JUMP:
		case BC_INIT_THREAD:
			wg.Add(1)
			go VM.executeThread(v.obj.(*RObject).methods["def"].def)
		}
	}
}

func (VM *GobiesVM) executeBlock(t *Transaction, block *RObject, args []*RObject) {
	// Create clean call frame
	currentCallFrame := initCallFrame()
	currentCallFrame.parent = t.localCallFrameStack[len(t.localCallFrameStack)-1]
	currentCallFrame.me = currentCallFrame.parent.me
	t.localCallFrameStack = append(t.localCallFrameStack, currentCallFrame)

	// Fill in arguments to current call frame
	if block.ivars["params"] != nil {
		params := block.ivars["params"].(*RObject).ivars["array"].([]*RObject)
		for i, v := range params {
			var_name := v.val.str
			currentCallFrame.var_table[var_name] = args[i]
		}
	}

	// Execute block definition
	VM.executeBytecodes(block.methods["def"].def, t)

	// Pop temporary call frame
	t.localCallFrameStack = t.localCallFrameStack[:len(t.localCallFrameStack)-1]
}
