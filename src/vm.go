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
	VM.consts["RFlonum"] = initRFlonum()
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

func (currentCallFrame *CallFrame) variableTableLookup(var_name string) map[string]Object {
	if _, ok := currentCallFrame.var_table[var_name]; ok {
		return currentCallFrame.var_table
	}
	if currentCallFrame.parent != nil {
		return currentCallFrame.parent.variableTableLookup(var_name)
	}
	return nil
}

func (t *Transaction) initTransaction(env *ThreadEnv, instList []Instruction) *Transaction {
	t.instList = []Instruction{}
	for _, inst := range instList {
		t.instList = append(t.instList, inst)
	}
	t.transactionStack = copyFrames(env.threadStack, false)
	t.objectSet = make(map[*RObject]*RObject)
	return t
}

func copyFrames(src []*CallFrame, initThread bool) []*CallFrame {
	var newStack []*CallFrame
	if initThread {
		newFrame := initCallFrame()
		newFrame.parent = nil
		newFrame.stack = make([]Object, 0, 0)
		// fmt.Println(test, "=>", src)
		for _, frame := range src {
			newFrame.me = frame.me
			// Copy every object "pointer" in instance variable table
			for key, obj := range frame.var_table {
				if obj.(*RObject).blockVar {
					newFrame.var_table[key] = obj.(*RObject).copyObject()
				} else {
					newFrame.var_table[key] = obj
				}
			}
		}
		newStack = []*CallFrame{newFrame}
	} else {
		newStack = make([]*CallFrame, len(src), len(src))
		for i := 0; i < len(src)-1; i++ {
			newStack[i] = src[i]
		}

		frame := src[len(src)-1]
		newFrame := initCallFrame()
		newFrame.me = frame.me
		newFrame.parent = frame.parent
		newFrame.stack = make([]Object, len(frame.stack), len(frame.stack))

		// Copy every object "pointer" in old frame stack
		for j, obj := range frame.stack {
			newFrame.stack[j] = obj
		}

		// Copy every object "pointer" in instance variable table
		for key, obj := range frame.var_table {
			newFrame.var_table[key] = obj
		}

		newStack[len(src)-1] = newFrame

	}

	return newStack
}

func (obj *RObject) copyObjectFrom(src *RObject) {
	// Copy instance variables
	ivars := src.ivars
	for k, v := range ivars {
		obj.ivars[k] = v
	}
}

func (obj *RObject) copyObject() *RObject {
	new_obj := &RObject{}
	new_obj.class = obj.class
	new_obj.rev = obj.rev
	new_obj.val = obj.val
	new_obj.blockVar = obj.blockVar
	new_obj.ivars = make(map[string]Object)
	new_obj.methods = make(map[string]*RMethod)

	ivars := obj.ivars
	for k, v := range ivars {
		new_obj.ivars[k] = v
	}
	return new_obj
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
	VM.executeThread(VM.instList, nil)
	wg.Wait()
}

func (VM *GobiesVM) executeThread(instList []Instruction, parentScope *ThreadEnv) {
	VM.globalLock.Lock()
	// Create clean call frame without pushing back to VM stack
	currentCallFrame := initCallFrame()
	env := &ThreadEnv{instList: instList}
	if parentScope == nil {
		currentCallFrame.parent = VM.callFrameStack[len(VM.callFrameStack)-1]
	} else {
		env.threadStack = copyFrames(parentScope.threadStack, true)
		currentCallFrame.parent = env.threadStack[len(env.threadStack)-1]
		// currentCallFrame.parent = parentScope.threadStack[len(parentScope.threadStack)-1]
	}
	VM.globalLock.Unlock()

	currentCallFrame.me = currentCallFrame.parent.me
	env.threadStack = append(env.threadStack, currentCallFrame)

	// t.inevitable = true
	VM.executeBytecodes(nil, env)
	if parentScope != nil {
		wg.Done()
	}
}

func (VM *GobiesVM) transactionBegin(env *ThreadEnv, inst []Instruction) *Transaction {
	t := &Transaction{}
	t.initTransaction(env, inst)

	VM.globalLock.Lock()
	t.rev = atomic.LoadInt64(&VM.rev)
	env.transactionPC = t
	t.env = env
	VM.globalLock.Unlock()

	return t
}

func (VM *GobiesVM) transactionEnd(env *ThreadEnv) bool {
	t := env.transactionPC
	locked := []*RObject{}

	// Validate the read-set
	for orig_obj, _ := range t.objectSet {
		if orig_obj.rev > t.rev {
			// Retry
			goto TRANSACTION_RETRY
		}
	}

	// Lock the write-set
	for orig_obj, new_obj := range t.objectSet {
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

	VM.globalLock.Lock()
	// Increment global revision
	for !atomic.CompareAndSwapInt64(&VM.rev, VM.rev, VM.rev+1) {
	}

	// Re-validate the read-set
	for orig_obj, _ := range t.objectSet {
		if orig_obj.rev > t.rev {
			// Release write locks
			for _, locked_obj := range locked {
				locked_obj.writeLock.Unlock()
			}
			VM.globalLock.Unlock()
			// Retry
			goto TRANSACTION_RETRY
		}
	}

	// Commit then release the locks
	for _, locked_obj := range locked {
		src := t.objectSet[locked_obj]
		locked_obj.copyObjectFrom(src)
		locked_obj.rev = atomic.LoadInt64(&VM.rev)
	}
	VM.globalLock.Unlock()

	for _, locked_obj := range locked {
		locked_obj.writeLock.Unlock()
	}

	env.threadStack = copyFrames(t.transactionStack, false)
	env.transactionPC = nil
	return true

TRANSACTION_RETRY:
	t.initTransaction(env, t.instList)
	VM.globalLock.Lock()
	t.rev = atomic.LoadInt64(&VM.rev)
	VM.globalLock.Unlock()
	VM.executeBytecodes(t.instList, env)
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
		// Validate the read-set
		for orig_obj, _ := range t.objectSet {
			if orig_obj.rev > t.rev {
				// Retry
				goto TRANSACTION_RETRY
			}
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
			table := currentCallFrame.variableTableLookup(v.obj.(string))
			if table != nil {
				table[v.obj.(string)] = top
			} else {
				currentCallFrame.var_table[v.obj.(string)] = top
			}
			currentCallFrame.stack = currentCallFrame.stack[0 : len(currentCallFrame.stack)-1] // Pop object from stack
		case BC_GETLOCAL:
			if VM.transactionEnd(env) == false {
				return
			}
			currentCallFrame = env.threadStack[len(env.threadStack)-1]
			obj := currentCallFrame.variableLookup(v.obj.(string)).(*RObject)
			currentCallFrame.stack = append(currentCallFrame.stack, obj)

		case BC_SETGLOBAL:
		case BC_GETGLOBAL:
		case BC_SETSYMBOL:
		case BC_GETSYMBOL:
		case BC_SETCONST:
		case BC_GETCONST:
			if v.obj.(string) == "Join" {
				wg.Wait()
			} else {
				currentCallFrame.stack = append(currentCallFrame.stack, VM.consts[v.obj.(string)])
			}
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
			// Update address since some functions might init new transaction
			currentCallFrame = env.transactionPC.transactionStack[len(env.transactionPC.transactionStack)-1]
			if return_val != nil {
				currentCallFrame.stack = append(currentCallFrame.stack, return_val)
			}
			if VM.transactionEnd(env) == false {
				return
			}
		case BC_JUMP:

		}
	}
	// End transaction if any
	if env.transactionPC != nil {
		VM.transactionEnd(env)
	}
	return

TRANSACTION_RETRY:
	t.initTransaction(env, t.instList)
	VM.globalLock.Lock()
	t.rev = atomic.LoadInt64(&VM.rev)
	VM.globalLock.Unlock()
	VM.executeBytecodes(t.instList, env)
}

func (VM *GobiesVM) executeBlock(env *ThreadEnv, block *RObject, args []*RObject) {
	if env.transactionPC != nil { // Before
		VM.transactionEnd(env)
	}

	// Create clean call frame
	currentCallFrame := initCallFrame()
	currentCallFrame.parent = env.threadStack[len(env.threadStack)-1]
	currentCallFrame.me = currentCallFrame.parent.me
	env.threadStack = append(env.threadStack, currentCallFrame)

	// Fill in arguments to current call frame
	if block.ivars["params"] != nil {
		params := block.ivars["params"].(*RObject).ivars["array"].([]*RObject)
		for i, v := range params {
			var_name := v.val.str
			currentCallFrame.var_table[var_name] = args[i].copyObject()
			currentCallFrame.var_table[var_name].(*RObject).blockVar = true
		}
	}

	VM.transactionBegin(env, block.methods["def"].def)

	// Execute block definition
	VM.executeBytecodes(block.methods["def"].def, env)
	env.threadStack = env.threadStack[0 : len(env.threadStack)-1]

	VM.transactionBegin(env, []Instruction{})
}
