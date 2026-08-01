package main

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"math/big"
	"slices"
	"sync"
	"sync/atomic"
)

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
	stateLock      sync.RWMutex
	wg             sync.WaitGroup
	errLock        sync.Mutex
	errs           []error
	failed         atomic.Bool
	outLock        sync.Mutex
	out            io.Writer
	logTransaction bool
	stats          *transactionStats
}

func initVM(out io.Writer, logTransaction bool, stats *transactionStats) *GobiesVM {
	VM := &GobiesVM{out: out, logTransaction: logTransaction, stats: stats}
	top := initCallFrame()
	VM.callFrameStack = append(VM.callFrameStack, top)
	VM.consts = make(map[string]*RObject)
	VM.initConsts()
	VM.symbols = make(map[string]int)
	top.me = initRKernel()
	return VM
}

func (VM *GobiesVM) recordError(err error) {
	if err == nil {
		return
	}
	VM.failed.Store(true)
	VM.errLock.Lock()
	VM.errs = append(VM.errs, err)
	VM.errLock.Unlock()
}

func (VM *GobiesVM) hasErrors() bool {
	return VM.failed.Load()
}

func (VM *GobiesVM) abortTransaction(env *ThreadEnv) {
	if env.transactionPC == nil {
		return
	}
	VM.stateLock.RUnlock()
	env.transactionPC = nil
}

func (VM *GobiesVM) println(values ...any) {
	VM.outLock.Lock()
	defer VM.outLock.Unlock()
	fmt.Fprintln(VM.out, values...)
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
	obj.class = src.class
	obj.rev = src.rev
	obj.name = src.name
	obj.val = cloneRValue(src.val)
	obj.blockVar = src.blockVar
	obj.ivars = cloneIvars(src.ivars)
	obj.methods = maps.Clone(src.methods)
}

func (obj *RObject) copyObject() *RObject {
	return &RObject{
		class:    obj.class,
		methods:  maps.Clone(obj.methods),
		name:     obj.name,
		val:      cloneRValue(obj.val),
		rev:      obj.rev,
		ivars:    cloneIvars(obj.ivars),
		blockVar: obj.blockVar,
	}
}

func cloneRValue(value RValue) RValue {
	cloned := value
	if value.fixnum != nil {
		cloned.fixnum = new(big.Int).Set(value.fixnum)
	}
	return cloned
}

func cloneIvars(ivars map[string]Object) map[string]Object {
	cloned := make(map[string]Object, len(ivars))
	for key, value := range ivars {
		switch value := value.(type) {
		case []*RObject:
			cloned[key] = slices.Clone(value)
		case map[RValue]*RObject:
			cloned[key] = maps.Clone(value)
		default:
			cloned[key] = value
		}
	}
	return cloned
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

func (VM *GobiesVM) execute() error {
	VM.rev = 0

	VM.executeThread(VM.instList, nil)
	VM.wg.Wait()

	VM.errLock.Lock()
	defer VM.errLock.Unlock()
	return errors.Join(VM.errs...)
}

func (VM *GobiesVM) executeThread(instList []Instruction, parentScope *ThreadEnv) {
	VM.stateLock.RLock()
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
	VM.stateLock.RUnlock()

	currentCallFrame.me = currentCallFrame.parent.me
	env.threadStack = append(env.threadStack, currentCallFrame)

	VM.executeBytecodes(nil, env)
}

func (VM *GobiesVM) transactionBegin(env *ThreadEnv, inst []Instruction) *Transaction {
	if env.transactionPC != nil {
		return env.transactionPC
	}

	t := &Transaction{}
	t.initTransaction(env, inst)

	VM.stateLock.RLock()
	t.rev = VM.rev
	env.transactionPC = t
	t.env = env

	return t
}

func (VM *GobiesVM) transactionEnd(env *ThreadEnv) bool {
	t := env.transactionPC
	if t == nil {
		return true
	}

	VM.stateLock.RUnlock()
	VM.stateLock.Lock()
	rev := int64(0)

	// Validate the read-set
	for orig_obj := range t.objectSet {
		if orig_obj.rev > t.rev {
			VM.stateLock.Unlock()
			if VM.logTransaction {
				VM.stats.status[LOG_COMMIT_VALID_READ].Add(1)
			}
			goto TRANSACTION_RETRY
		}
	}

	VM.rev++
	rev = VM.rev

	// Publish transaction-local copies while no speculative reader can observe
	// partially committed object state.
	for original, local := range t.objectSet {
		if original != local {
			original.copyObjectFrom(local)
			original.rev = rev
		}
	}

	VM.stateLock.Unlock()

	env.threadStack = copyFrames(t.transactionStack, false)
	env.transactionPC = nil

	if VM.logTransaction {
		VM.stats.status[LOG_COMMIT_SUCCESS].Add(1)
	}

	return true

TRANSACTION_RETRY:
	t.initTransaction(env, t.instList)
	VM.stateLock.RLock()
	t.rev = VM.rev
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
		if VM.hasErrors() {
			VM.abortTransaction(env)
			return
		}
		t = env.transactionPC
		if t == nil {
			t = VM.transactionBegin(env, instList[i:])
		}
		// Validate the read-set
		for orig_obj := range t.objectSet {
			if orig_obj.rev > t.rev {
				// Retry
				if VM.logTransaction {
					VM.stats.status[LOG_EXEC_VALID_READ].Add(1)
				}
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
				if VM.transactionEnd(env) == false {
					return
				}
				VM.wg.Wait()
				VM.transactionBegin(env, []Instruction{})
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
			if VM.hasErrors() {
				VM.abortTransaction(env)
				return
			}
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
	VM.stateLock.RUnlock()
	t.initTransaction(env, t.instList)
	VM.stateLock.RLock()
	t.rev = VM.rev
	VM.executeBytecodes(t.instList, env)
}

func (VM *GobiesVM) executeBlock(env *ThreadEnv, block *RObject, args []*RObject) {
	if env.transactionPC != nil { // Before
		VM.transactionEnd(env)
	}

	VM.stateLock.RLock()
	definition := slices.Clone(block.methods["def"].def)
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
	VM.stateLock.RUnlock()

	VM.transactionBegin(env, definition)

	// Execute block definition
	VM.executeBytecodes(definition, env)
	env.threadStack = env.threadStack[0 : len(env.threadStack)-1]

	VM.transactionBegin(env, []Instruction{})
}
