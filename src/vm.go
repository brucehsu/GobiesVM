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
	startInst        Instruction
	instList         []Instruction
	objectSet        map[*RObject]*RObject
	transactionStack []*CallFrame
	env              *ThreadEnv
	inevitable       bool
	rev              int64
}

type ThreadEnv struct {
	instList      []Instruction
	threadStack   []*CallFrame
	transactionPC *Transaction
	id            int
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
	VM.consts["Thread"] = initRThread()
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

func (t *Transaction) initTransaction(env *ThreadEnv, instList []Instruction) *Transaction {
	t.instList = []Instruction{}
	for _, inst := range instList {
		t.instList = append(t.instList, inst)
	}
	t.transactionStack = copyFrames(env.threadStack)
	t.objectSet = make(map[*RObject]*RObject)
	return t
}

func copyFrames(src []*CallFrame) []*CallFrame {
	newStack := []*CallFrame{}
	for _, frame := range src {
		newFrame := initCallFrame()
		newFrame.me = frame.me
		newFrame.parent = frame.parent

		// Copy every object "pointer" in old frame stack
		for _, obj := range frame.stack {
			newFrame.stack = append(newFrame.stack, obj)
		}

		// Copy every object "pointer" in instance variable table
		for key, obj := range frame.var_table {
			newFrame.var_table[key] = obj
		}

		newStack = append(newStack, newFrame)
	}
	return newStack
}

func (obj *RObject) copyObject(src *RObject) {
	// Copy instance variables
	ivars := src.ivars
	for k, v := range ivars {
		obj.ivars[k] = v
	}
}

func addRObjectToSet(obj *RObject, env *ThreadEnv) *RObject {
	if env == nil { // created during complication
		return obj
	}
	if _, ok := env.transactionPC.objectSet[obj]; !ok {
		env.transactionPC.objectSet[obj] = obj
	}
	return obj
}

func (VM *GobiesVM) execute() {
	// 	root := initTransaction(VM.instList)
	// 	root.inevitable = true

	VM.rev = 0

	// Execute root transaction which is inevitable
	wg.Add(1)
	go VM.executeThread(VM.instList, nil)
	wg.Wait()
}

func (VM *GobiesVM) executeThread(instList []Instruction, parentScope *ThreadEnv) {
	// Create clean call frame without pushing back to VM stack
	currentCallFrame := initCallFrame()

	env := &ThreadEnv{instList: instList}
	if parentScope == nil {
		currentCallFrame.parent = VM.callFrameStack[len(VM.callFrameStack)-1]
	} else {
		env.threadStack = copyFrames(parentScope.threadStack)
		currentCallFrame.parent = env.threadStack[len(env.threadStack)-1]
		// currentCallFrame.parent = parentScope.threadStack[len(parentScope.threadStack)-1]
	}
	currentCallFrame.me = currentCallFrame.parent.me
	env.threadStack = append(env.threadStack, currentCallFrame)

	// t.inevitable = true
	VM.executeBytecodes(nil, env)
	wg.Done()
}

func (VM *GobiesVM) transactionBegin(env *ThreadEnv, inst []Instruction) *Transaction {
	t := &Transaction{}
	t.initTransaction(env, inst)
	t.rev = atomic.LoadInt64(&VM.rev)
	env.transactionPC = t
	t.env = env

	// Initialize environment
	if t.inevitable {
		VM.globalLock.Lock()
		atomic.StoreInt32(&VM.has_inevitable, 1)
	} else {
		VM.globalLock.RLock()
	}

	return t
}

func (VM *GobiesVM) transactionEnd(env *ThreadEnv) bool {
	t := env.transactionPC

	// Lock the write-set
	locked := []*RObject{}
	for orig_obj, new_obj := range t.objectSet {
		// fmt.Println(orig_obj, new_obj)
		if orig_obj != new_obj {
			if orig_obj.writeLock.TryLock() {
				locked = append(locked, orig_obj)
			} else { // Attempt to acquire lock failed
				// Release all locks
				for _, locked_obj := range locked {
					locked_obj.writeLock.Unlock()
				}
				// Retry
				goto TRANSACTION_RETRY

			}
		}
	}

	// Increment global revision
	for !atomic.CompareAndSwapInt64(&VM.rev, VM.rev, VM.rev+1) {
	}

	// Validate the read-set
	for orig_obj, new_obj := range t.objectSet {
		if orig_obj == new_obj {
			if orig_obj.rev > t.rev {
				// Release write locks
				for _, locked_obj := range locked {
					locked_obj.writeLock.Unlock()
				}
				// Retry
				goto TRANSACTION_RETRY
			}
		}
	}

	// Commit then release the locks
	for _, locked_obj := range locked {
		src := t.objectSet[locked_obj]
		locked_obj.copyObject(src)
	}
	for _, locked_obj := range locked {
		locked_obj.writeLock.Unlock()
	}

	env.threadStack = copyFrames(t.transactionStack)

	if t.inevitable {
		VM.globalLock.Unlock()
		atomic.StoreInt32(&VM.has_inevitable, 0)
	} else {
		VM.globalLock.RUnlock()
	}

	env.transactionPC = nil
	return true

TRANSACTION_RETRY:
	t.initTransaction(env, t.instList)
	t.rev = atomic.LoadInt64(&VM.rev)
	return false
}

func (VM *GobiesVM) executeBytecodes(instList []Instruction, env *ThreadEnv) {
	t := env.transactionPC

	if instList == nil {
		instList = env.instList
	}

	// SPECULATIVE_EXEC:
	// Speculative execution
	for i, v := range instList {
		t = env.transactionPC
		if t == nil {
			t = VM.transactionBegin(env, instList[i:])
		}
		currentCallFrame := t.transactionStack[len(t.transactionStack)-1]

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
			currentCallFrame.stack = currentCallFrame.stack[0 : len(currentCallFrame.stack)-1] // Pop object from stack
		case BC_GETLOCAL:
			obj := currentCallFrame.variableLookup(v.obj.(string)).(*RObject)
			if _, ok := t.objectSet[obj]; ok {
				obj = t.objectSet[obj]
			}
			currentCallFrame.stack = append(currentCallFrame.stack, obj)
			if VM.transactionEnd(env) == false {

			}
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
			return_val := recv.methodLookup(v.obj.(string)).gofunc(VM, env, recv, argLists)
			// fmt.Println(env.transactionPC)
			// Update address since some functions might init new transaction
			currentCallFrame = env.transactionPC.transactionStack[len(env.transactionPC.transactionStack)-1]
			if return_val != nil {
				currentCallFrame.stack = append(currentCallFrame.stack, return_val)
			}
			if VM.transactionEnd(env) == false {

			}
		case BC_JUMP:

		}
	}
	// End transaction if any
	if env.transactionPC != nil {
		VM.transactionEnd(env)
	}
}

func (VM *GobiesVM) executeBlock(env *ThreadEnv, block *RObject, args []*RObject) {
	if env.transactionPC != nil { // Before
		VM.transactionEnd(env)
	}
	t := VM.transactionBegin(env, block.methods["def"].def)

	// Create clean call frame
	currentCallFrame := initCallFrame()
	currentCallFrame.parent = t.transactionStack[len(t.transactionStack)-1]
	currentCallFrame.me = currentCallFrame.parent.me
	t.transactionStack = append(t.transactionStack, currentCallFrame)

	// Fill in arguments to current call frame
	if block.ivars["params"] != nil {
		params := block.ivars["params"].(*RObject).ivars["array"].([]*RObject)
		for i, v := range params {
			var_name := v.val.str
			currentCallFrame.var_table[var_name] = args[i]
		}
	}

	// Execute block definition
	VM.executeBytecodes(block.methods["def"].def, env)

	VM.transactionBegin(env, []Instruction{})
}
