package main

type NodeType uint8

// Resembles TinyRb AST node types
const (
	NODE_ROOT NodeType = iota
	NODE_LIST
	NODE_BLOCK
	NODE_VALUE
	NODE_ASTVAL
	NODE_STRING
	NODE_ASSIGN
	NODE_ARG
	NODE_SEND
	NODE_MSG
	NODE_IF
	NODE_UNLESS
	NODE_AND
	NODE_OR
	NODE_WHILE
	NODE_UNTIL
	NODE_TRUE
	NODE_FALSE
	NODE_NIL
	NODE_SELF
	NODE_LEAVE
	NODE_RETURN
	NODE_BREAK
	NODE_YIELD
	NODE_DEF
	NODE_METHOD
	NODE_PARAM
	NODE_CLASS
	NODE_MODULE
	NODE_CONST
	NODE_SETCONST
	NODE_ARRAY
	NODE_HASH
	NODE_RANGE
	NODE_GETIVAR
	NODE_SETIVAR
	NODE_GETCVAR
	NODE_SETCVAR
	NODE_GETGLOBAL
	NODE_SETGLOBAL
	NODE_ADD
	NODE_SUB
	NODE_LT
	NODE_NEG
	NODE_NOT
)

var node_type_str = [...]string{
	"NODE_ROOT",
	"NODE_LIST",
	"NODE_BLOCK",
	"NODE_VALUE",
	"NODE_ASTVAL",
	"NODE_STRING",
	"NODE_ASSIGN",
	"NODE_ARG",
	"NODE_SEND",
	"NODE_MSG",
	"NODE_IF",
	"NODE_UNLESS",
	"NODE_AND",
	"NODE_OR",
	"NODE_WHILE",
	"NODE_UNTIL",
	"NODE_TRUE",
	"NODE_FALSE",
	"NODE_NIL",
	"NODE_SELF",
	"NODE_LEAVE",
	"NODE_RETURN",
	"NODE_BREAK",
	"NODE_YIELD",
	"NODE_DEF",
	"NODE_METHOD",
	"NODE_PARAM",
	"NODE_CLASS",
	"NODE_MODULE",
	"NODE_CONST",
	"NODE_SETCONST",
	"NODE_ARRAY",
	"NODE_HASH",
	"NODE_RANGE",
	"NODE_GETIVAR",
	"NODE_SETIVAR",
	"NODE_GETCVAR",
	"NODE_SETCVAR",
	"NODE_GETGLOBAL",
	"NODE_SETGLOBAL",
	"NODE_ADD",
	"NODE_SUB",
	"NODE_LT",
	"NODE_NEG",
	"NODE_NOT",
}

type ASTVal struct {
	str     string
	numeric int64
}

type AST struct {
	Type   NodeType
	length int
	line   int
	prev   *AST
	next   *AST
	head   *AST
	tail   *AST
	args   [3]*AST
	value  ASTVal
}

func (node *AST) PushFront(child *AST) *AST {
	if child == nil {
		return node
	}
	node.length += 1
	if node.tail == nil {
		node.tail = child
	}
	child.next = node.head
	if node.head != nil {
		node.head.prev = child
	}
	node.head = child
	return node
}

func (node *AST) PushBack(child *AST) *AST {
	if child == nil {
		return node
	}
	if node.head == nil {
		node.head = child
	}
	node.length += 1
	child.prev = node.tail
	if node.tail != nil {
		node.tail.next = child
	}
	node.tail = child
	return node
}

func (node *AST) AddArgs(first *AST, second *AST, third *AST) *AST {
	node.args[0] = first
	node.args[1] = second
	node.args[2] = third
	return node
}

func (node *AST) PopFront() *AST {
	node.length -= 1
	front := node.head
	node.head = front.next
	if node.head != nil {
		node.head.prev = nil
	}
	front.next = nil
	return front
}

func (node *AST) PopBack() *AST {
	node.length -= 1
	back := node.tail
	node.tail = back.prev
	if node.tail != nil {
		node.tail.next = nil
	}
	back.prev = nil
	return back
}

func (node *AST) Front() *AST {
	return node.head
}

func (node *AST) Back() *AST {
	return node.tail
}

func MakeASTNode(ntype NodeType, arg1 *AST, arg2 *AST, arg3 *AST, line int) *AST {
	node := &AST{Type: ntype, line: line, length: 1}
	node.AddArgs(arg1, arg2, arg3)
	return node
}
