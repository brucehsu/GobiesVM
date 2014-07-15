package main

import (
	/*"bytes"*/
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

const END_SYMBOL rune = 4

var traverse_level int = 0

func Traverse(root *AST) {
	for i := 0; i < traverse_level; i++ {
		fmt.Printf("  ")
	}
	traverse_level += 1
	fmt.Printf("%s: %+v\n", node_type_str[root.Type], root)

	root_head := root.head

	for idx, _ := range root.args {
		if root.args[idx] != nil {
			Traverse(root.args[idx])
		}
	}

	traverse_level -= 1

	for root_head != nil {
		Traverse(root_head)
		root_head = root_head.next
	}

}

var currentLine int = 0
var rootAST *AST
var parseStr []string

/* The rule types inferred from the grammar are below. */
type Rule uint8

const (
	RuleUnknown Rule = iota
	RuleRoot
	RuleStmts
	RuleOptStmts
	RuleStmt
	RuleExpr
	RuleComment
	RuleCall
	RuleAsgnCall
	RuleReceiver
	RuleSpecCall
	RuleBinOp
	RuleUnaryOp
	RuleMessage
	RuleArgs
	RuleBlock
	RuleAssign
	RuleWhile
	RuleUntil
	RuleIf
	RuleUnless
	RuleElse
	RuleMethod
	RuleDef
	RuleParams
	RuleParam
	RuleClass
	RuleModule
	RuleRange
	RuleYield
	RuleReturn
	RuleBreak
	RuleValue
	RuleAryItems
	RuleHashItems
	RuleKEYWORD
	RuleNAME
	RuleID
	RuleCONST
	RuleBINOP
	RuleUNOP
	RuleMETHOD
	RuleASSIGN
	RuleIVAR
	RuleCVAR
	RuleGLOBAL
	RuleNUMBER
	RuleSYMBOL
	RuleSTRING1
	Rule_
	RuleSPACE
	RuleEOL
	RuleEOF
	RuleSEP
	RuleAction0
	RuleAction1
	RuleAction2
	RuleAction3
	RuleAction4
	RuleAction5
	RuleAction6
	RuleAction7
	RuleAction8
	RuleAction9
	RuleAction10
	RuleAction11
	RuleAction12
	RuleAction13
	RuleAction14
	RuleAction15
	RuleAction16
	RuleAction17
	RuleAction18
	RuleAction19
	RuleAction20
	RuleAction21
	RuleAction22
	RuleAction23
	RuleAction24
	RuleAction25
	RuleAction26
	RuleAction27
	RuleAction28
	RuleAction29
	RuleAction30
	RuleAction31
	RuleAction32
	RuleAction33
	RuleAction34
	RuleAction35
	RuleAction36
	RuleAction37
	RuleAction38
	RuleAction39
	RuleAction40
	RuleAction41
	RuleAction42
	RuleAction43
	RuleAction44
	RuleAction45
	RuleAction46
	RuleAction47
	RuleAction48
	RuleAction49
	RuleAction50
	RuleAction51
	RuleAction52
	RuleAction53
	RuleAction54
	RuleAction55
	RuleAction56
	RuleAction57
	RuleAction58
	RuleAction59
	RuleAction60
	RuleAction61
	RuleAction62
	RuleAction63
	RuleAction64
	RuleAction65
	RuleAction66
	RuleAction67
	RuleAction68
	RuleAction69
	RuleAction70
	RuleAction71
	RuleAction72
	RuleAction73
	RuleAction74
	RuleAction75
	RuleAction76
	RuleAction77
	RuleAction78
	RuleAction79
	RuleAction80
	RuleAction81
	RuleAction82
	RuleAction83
	RuleAction84
	RuleAction85
	RuleAction86
	RuleAction87
	RuleAction88
	RuleAction89
	RuleAction90
	RuleAction91
	RuleAction92
	RuleAction93
	RulePegText
	RuleAction94
	RuleAction95
	RuleAction96
	RuleAction97
	RuleAction98
	RuleAction99
	RuleAction100
	RuleAction101
	RuleAction102
	RuleAction103
	RuleAction104
	RuleAction105
	RuleAction106
	RuleAction107
	RuleAction108
	RuleAction109
	RuleAction110

	RuleActionPush
	RuleActionPop
	RuleActionSet
	RulePre_
	Rule_In_
	Rule_Suf
)

var Rul3s = [...]string{
	"Unknown",
	"Root",
	"Stmts",
	"OptStmts",
	"Stmt",
	"Expr",
	"Comment",
	"Call",
	"AsgnCall",
	"Receiver",
	"SpecCall",
	"BinOp",
	"UnaryOp",
	"Message",
	"Args",
	"Block",
	"Assign",
	"While",
	"Until",
	"If",
	"Unless",
	"Else",
	"Method",
	"Def",
	"Params",
	"Param",
	"Class",
	"Module",
	"Range",
	"Yield",
	"Return",
	"Break",
	"Value",
	"AryItems",
	"HashItems",
	"KEYWORD",
	"NAME",
	"ID",
	"CONST",
	"BINOP",
	"UNOP",
	"METHOD",
	"ASSIGN",
	"IVAR",
	"CVAR",
	"GLOBAL",
	"NUMBER",
	"SYMBOL",
	"STRING1",
	"_",
	"SPACE",
	"EOL",
	"EOF",
	"SEP",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",
	"Action36",
	"Action37",
	"Action38",
	"Action39",
	"Action40",
	"Action41",
	"Action42",
	"Action43",
	"Action44",
	"Action45",
	"Action46",
	"Action47",
	"Action48",
	"Action49",
	"Action50",
	"Action51",
	"Action52",
	"Action53",
	"Action54",
	"Action55",
	"Action56",
	"Action57",
	"Action58",
	"Action59",
	"Action60",
	"Action61",
	"Action62",
	"Action63",
	"Action64",
	"Action65",
	"Action66",
	"Action67",
	"Action68",
	"Action69",
	"Action70",
	"Action71",
	"Action72",
	"Action73",
	"Action74",
	"Action75",
	"Action76",
	"Action77",
	"Action78",
	"Action79",
	"Action80",
	"Action81",
	"Action82",
	"Action83",
	"Action84",
	"Action85",
	"Action86",
	"Action87",
	"Action88",
	"Action89",
	"Action90",
	"Action91",
	"Action92",
	"Action93",
	"PegText",
	"Action94",
	"Action95",
	"Action96",
	"Action97",
	"Action98",
	"Action99",
	"Action100",
	"Action101",
	"Action102",
	"Action103",
	"Action104",
	"Action105",
	"Action106",
	"Action107",
	"Action108",
	"Action109",
	"Action110",

	"RuleActionPush",
	"RuleActionPop",
	"RuleActionSet",
	"Pre_",
	"_In_",
	"_Suf",
}

type TokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule Rule, begin, end, next, depth int)
	Expand(index int) TokenTree
	Tokens() <-chan token32
	Error() []token32
	trim(length int)
}

/* ${@} bit structure for abstract syntax tree */
type token16 struct {
	Rule
	begin, end, next int16
}

func (t *token16) isZero() bool {
	return t.Rule == RuleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token16) isParentOf(u token16) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token16) GetToken32() token32 {
	return token32{Rule: t.Rule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token16) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", Rul3s[t.Rule], t.begin, t.end, t.next)
}

type tokens16 struct {
	tree    []token16
	ordered [][]token16
}

func (t *tokens16) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens16) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens16) Order() [][]token16 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int16, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.Rule == RuleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token16, len(depths)), make([]token16, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int16(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type State16 struct {
	token16
	depths []int16
	leaf   bool
}

func (t *tokens16) PreOrder() (<-chan State16, [][]token16) {
	s, ordered := make(chan State16, 6), t.Order()
	go func() {
		var states [8]State16
		for i, _ := range states {
			states[i].depths = make([]int16, len(ordered))
		}
		depths, state, depth := make([]int16, len(ordered)), 0, 1
		write := func(t token16, leaf bool) {
			S := states[state]
			state, S.Rule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.Rule, t.begin, t.end, int16(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token16 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token16{Rule: Rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token16{Rule: RulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.Rule != RuleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.Rule != RuleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token16{Rule: Rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens16) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", Rul3s[ordered[i][depths[i]-1].Rule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", Rul3s[token.Rule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", Rul3s[ordered[i][depths[i]-1].Rule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", Rul3s[token.Rule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", Rul3s[ordered[i][depths[i]-1].Rule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", Rul3s[token.Rule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens16) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", Rul3s[token.Rule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens16) Add(rule Rule, begin, end, depth, index int) {
	t.tree[index] = token16{Rule: rule, begin: int16(begin), end: int16(end), next: int16(depth)}
}

func (t *tokens16) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.GetToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens16) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].GetToken32()
		}
	}
	return tokens
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	Rule
	begin, end, next int32
}

func (t *token32) isZero() bool {
	return t.Rule == RuleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) GetToken32() token32 {
	return token32{Rule: t.Rule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", Rul3s[t.Rule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.Rule == RuleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type State32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) PreOrder() (<-chan State32, [][]token32) {
	s, ordered := make(chan State32, 6), t.Order()
	go func() {
		var states [8]State32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.Rule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.Rule, t.begin, t.end, int32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{Rule: Rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{Rule: RulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.Rule != RuleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.Rule != RuleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{Rule: Rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", Rul3s[ordered[i][depths[i]-1].Rule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", Rul3s[token.Rule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", Rul3s[ordered[i][depths[i]-1].Rule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", Rul3s[token.Rule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", Rul3s[ordered[i][depths[i]-1].Rule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", Rul3s[token.Rule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", Rul3s[token.Rule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens32) Add(rule Rule, begin, end, depth, index int) {
	t.tree[index] = token32{Rule: rule, begin: int32(begin), end: int32(end), next: int32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.GetToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].GetToken32()
		}
	}
	return tokens
}

func (t *tokens16) Expand(index int) TokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		for i, v := range tree {
			expanded[i] = v.GetToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}

func (t *tokens32) Expand(index int) TokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type Leg struct {
	Buffer string
	buffer []rune
	rules  [166]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	TokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer[0:] {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p *Leg
}

func (e *parseError) Error() string {
	tokens, error := e.p.TokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			Rul3s[token.Rule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *Leg) PrintSyntaxTree() {
	p.TokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *Leg) Highlighter() {
	p.TokenTree.PrintSyntax()
}

func (p *Leg) Execute() {
	buffer, begin, end := p.Buffer, 0, 0

	var yy *AST
	stack := make([]*AST, 1024)
	stack_idx := 0

	for token := range p.TokenTree.Tokens() {
		switch token.Rule {
		case RulePegText:
			begin, end = int(token.begin), int(token.end)
		case RuleAction0:
			rootAST = MakeASTNode(NODE_ROOT, stack[stack_idx-0], nil, nil, currentLine)
		case RuleAction1:
			//Stmts:0
		case RuleAction2:
			stack[stack_idx-1].PushBack(stack[stack_idx-0]) //Stmts:1
		case RuleAction3:
			yy = stack[stack_idx-1] //Stmts:2
		case RuleAction4:
			yy = nil //Stmts:3
		case RuleAction5:
			yy = nil //OptStmts:0
		case RuleAction6:
			stack[stack_idx-3] = nil
			stack[stack_idx-0] = nil //Call:0
		case RuleAction7:
			currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-3], stack[stack_idx-2], nil, currentLine)
			stack[stack_idx-3] = currentAST //Call:1
		case RuleAction8:
			currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-3], stack[stack_idx-1], stack[stack_idx-0], currentLine)
			yy = currentAST //Call:2
		case RuleAction9:
			stack[stack_idx-4] = nil //AsgnCall:0
		case RuleAction10:
			currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-4], stack[stack_idx-3], nil, currentLine)
			stack[stack_idx-4] = currentAST //AsgnCall:1
		case RuleAction11:
			argAST := MakeASTNode(NODE_ARG, stack[stack_idx-0], nil, nil, currentLine)
			msgVal := []string{stack[stack_idx-2].value.str, "="}
			stack[stack_idx-2].value.str = strings.Join(msgVal, "")
			msgAST := MakeASTNode(NODE_MSG, stack[stack_idx-2], argAST, nil, currentLine)
			yy = MakeASTNode(NODE_SEND, stack[stack_idx-4], msgAST, nil, currentLine) //AsgnCall:2
		case RuleAction12:
			stack[stack_idx-0] = nil //Receiver:0
		case RuleAction13:
			yy = stack[stack_idx-0] //Receiver:1
		case RuleAction14:
			yy = stack[stack_idx-0] //Receiver:2
		case RuleAction15:
			currentAST := MakeASTNode(NODE_ARG, stack[stack_idx-0], nil, nil, currentLine)
			stack[stack_idx-1].PushBack(currentAST)
			methodAST := &AST{Type: NODE_ASTVAL}
			methodAST.value.str = "[]="
			msgAST := MakeASTNode(NODE_MSG, methodAST, stack[stack_idx-1], nil, currentLine)
			currentAST = MakeASTNode(NODE_SEND, stack[stack_idx-2], msgAST, nil, currentLine)
			yy = currentAST //SpecCall:0
		case RuleAction16:
			methodAST := &AST{Type: NODE_ASTVAL}
			methodAST.value.str = "[]"
			msgAST := MakeASTNode(NODE_MSG, methodAST, stack[stack_idx-1], nil, currentLine)
			currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-2], msgAST, nil, currentLine)
			yy = currentAST //SpecCall:1
		case RuleAction17:
			currentAST := MakeASTNode(NODE_AND, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine)
			yy = currentAST //BinOp:0 &&
		case RuleAction18:
			currentAST := MakeASTNode(NODE_OR, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine)
			yy = currentAST //BinOp:1 ||
		case RuleAction19:
			currentAST := MakeASTNode(NODE_ADD, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine)
			yy = currentAST //BinOp:2 +
		case RuleAction20:
			currentAST := MakeASTNode(NODE_SUB, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine)
			yy = currentAST //BinOp:3 -
		case RuleAction21:
			currentAST := MakeASTNode(NODE_LT, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine)
			yy = currentAST //BinOp:4 <
		case RuleAction22:
			argAST := MakeASTNode(NODE_ARG, stack[stack_idx-1], nil, nil, currentLine)
			msgAST := MakeASTNode(NODE_MSG, stack[stack_idx-0], argAST, nil, currentLine)
			currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-2], msgAST, nil, currentLine)
			yy = currentAST //BinOp:5 BINOP
		case RuleAction23:
			currentAST := MakeASTNode(NODE_NEG, stack[stack_idx-0], nil, nil, currentLine)
			yy = currentAST //UnaryOp:0 NODE_NEG -
		case RuleAction24:
			currentAST := MakeASTNode(NODE_NOT, stack[stack_idx-0], nil, nil, currentLine)
			yy = currentAST //UnaryOp:0 NODE_NOT !
		case RuleAction25:
			stack[stack_idx-0] = nil //Message:0
		case RuleAction26:
			currentAST := MakeASTNode(NODE_MSG, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Message:1
		case RuleAction27:
			currentAST := MakeASTNode(NODE_ARG, stack[stack_idx-2], nil, nil, currentLine)
			stack[stack_idx-2] = currentAST //Args:0
		case RuleAction28:
			stack[stack_idx-2].args[0].PushBack(stack[stack_idx-1]) //Args:1
		case RuleAction29:
			//Args:2 No going to support :p
		case RuleAction30:
			yy = stack[stack_idx-2] //Args:3
		case RuleAction31:
			currentAST := MakeASTNode(NODE_BLOCK, stack[stack_idx-1], nil, nil, currentLine)
			yy = currentAST //Block:0
		case RuleAction32:
			currentAST := MakeASTNode(NODE_BLOCK, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Block:1 with stack[stack_idx-0]
		case RuleAction33:
			currentAST := MakeASTNode(NODE_ASSIGN, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Assign:0
		case RuleAction34:
			currentAST := MakeASTNode(NODE_SETCONST, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Assign:1
		case RuleAction35:
			currentAST := MakeASTNode(NODE_SETIVAR, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Assign:2
		case RuleAction36:
			currentAST := MakeASTNode(NODE_SETCVAR, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Assign:3
		case RuleAction37:
			currentAST := MakeASTNode(NODE_SETGLOBAL, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Assign:4
		case RuleAction38:
			currentAST := MakeASTNode(NODE_WHILE, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //While NODE_WHILE
		case RuleAction39:
			currentAST := MakeASTNode(NODE_UNTIL, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Until NODE_UNTIL
		case RuleAction40:
			stack[stack_idx-0] = nil //If:0 stack[stack_idx-0] = 0
		case RuleAction41:
			currentAST := MakeASTNode(NODE_IF, stack[stack_idx-2], stack[stack_idx-1], stack[stack_idx-0], currentLine)
			yy = currentAST //If:1
		case RuleAction42:
			currentAST := MakeASTNode(NODE_IF, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine)
			yy = currentAST //If:2 No ELSE
		case RuleAction43:
			stack[stack_idx-0] = nil //Unless:0 stack[stack_idx-0] = 0
		case RuleAction44:
			currentAST := MakeASTNode(NODE_UNLESS, stack[stack_idx-2], stack[stack_idx-1], stack[stack_idx-0], currentLine)
			yy = currentAST //Unless:1
		case RuleAction45:
			currentAST := MakeASTNode(NODE_UNLESS, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine)
			yy = currentAST //Unless:2 No ELSE
		case RuleAction46:
			yy = stack[stack_idx-0] //Else:0
		case RuleAction47:
			msgAST := MakeASTNode(NODE_MSG, stack[stack_idx-1], nil, nil, currentLine)
			sendAST := MakeASTNode(NODE_SEND, nil, msgAST, nil, currentLine)
			currentAST := MakeASTNode(NODE_METHOD, sendAST, stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Method:0 NODE_METHOD
		case RuleAction48:
			currentAST := MakeASTNode(NODE_METHOD, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Method:1
		case RuleAction49:
			currentAST := MakeASTNode(NODE_METHOD, nil, stack[stack_idx-0], nil, currentLine)
			yy = currentAST //Method:2 Call top self
		case RuleAction50:
			stack[stack_idx-1] = nil //Def:0 stack[stack_idx-1]=0
		case RuleAction51:
			currentAST := MakeASTNode(NODE_DEF, stack[stack_idx-2], stack[stack_idx-1], stack[stack_idx-0], currentLine)
			yy = currentAST //Def:1 stack[stack_idx-1]>0
		case RuleAction52:
			//Params:0 stack[stack_idx-1]
		case RuleAction53:
			stack[stack_idx-1].PushBack(stack[stack_idx-0]) //Params:1
		case RuleAction54:
			yy = stack[stack_idx-1] //Params:2 yy = stack[stack_idx-1]
		case RuleAction55:
			yy = MakeASTNode(NODE_PARAM, stack[stack_idx-1], nil, stack[stack_idx-0], currentLine) //Param:0
		case RuleAction56:
			yy = MakeASTNode(NODE_PARAM, stack[stack_idx-1], nil, nil, currentLine) //Param:1
		case RuleAction57:
			//Param:2 splat TODO
		case RuleAction58:
			//Class:0 stack[stack_idx-1] = 0
		case RuleAction59:
			yy = MakeASTNode(NODE_CLASS, stack[stack_idx-2], stack[stack_idx-1], stack[stack_idx-0], currentLine) //Class:1
		case RuleAction60:
			yy = MakeASTNode(NODE_MODULE, stack[stack_idx-1], nil, stack[stack_idx-0], currentLine) //Module:0
		case RuleAction61:
			//Range:0 TODO
		case RuleAction62:
			//Range:1 TODO
		case RuleAction63:
			yy = MakeASTNode(NODE_YIELD, stack[stack_idx-0], nil, nil, currentLine) //Yield:0
		case RuleAction64:
			yy = MakeASTNode(NODE_YIELD, stack[stack_idx-0], nil, nil, currentLine) //Yield:1
		case RuleAction65:
			yy = MakeASTNode(NODE_YIELD, nil, nil, nil, currentLine) //Yield:2
		case RuleAction66:
			yy = MakeASTNode(NODE_RETURN, stack[stack_idx-1], nil, nil, currentLine) //Return:0
		case RuleAction67:
			yy = MakeASTNode(NODE_RETURN, stack[stack_idx-1], nil, nil, currentLine) //Return:1
		case RuleAction68:
			yy = MakeASTNode(NODE_RETURN, MakeASTNode(NODE_ARRAY, stack[stack_idx-0], nil, nil, currentLine), nil, nil, currentLine) //Return:2
		case RuleAction69:
			yy = MakeASTNode(NODE_RETURN, MakeASTNode(NODE_ARRAY, stack[stack_idx-0], nil, nil, currentLine), nil, nil, currentLine) //Return:3
		case RuleAction70:
			yy = MakeASTNode(NODE_RETURN, nil, nil, nil, currentLine) //Return:4
		case RuleAction71:
			//Break:0
		case RuleAction72:
			yy = MakeASTNode(NODE_VALUE, stack[stack_idx-2], nil, nil, currentLine) //Value:0 NUMBER
		case RuleAction73:
			yy = MakeASTNode(NODE_VALUE, stack[stack_idx-2], nil, nil, currentLine) //Value:1 SYMBOL
		case RuleAction74:
			yy = MakeASTNode(NODE_STRING, stack[stack_idx-2], nil, nil, currentLine) //Value:3 STRING1
		case RuleAction75:
			yy = MakeASTNode(NODE_CONST, stack[stack_idx-2], nil, nil, currentLine) //Value:5 CONST
		case RuleAction76:
			yy = MakeASTNode(NODE_NIL, nil, nil, nil, currentLine) //Value:6 nil
		case RuleAction77:
			yy = MakeASTNode(NODE_TRUE, nil, nil, nil, currentLine) //Value:7 true
		case RuleAction78:
			yy = MakeASTNode(NODE_FALSE, nil, nil, nil, currentLine) //Value:8 false
		case RuleAction79:
			yy = MakeASTNode(NODE_SELF, nil, nil, nil, currentLine) //Value:9 self
		case RuleAction80:
			yy = MakeASTNode(NODE_GETIVAR, stack[stack_idx-1], nil, nil, currentLine) //Value:10 IVAR
		case RuleAction81:
			yy = MakeASTNode(NODE_GETCVAR, stack[stack_idx-1], nil, nil, currentLine) //Value:11 CVAR
		case RuleAction82:
			yy = MakeASTNode(NODE_GETGLOBAL, stack[stack_idx-1], nil, nil, currentLine) //Value:12 GLOBAL
		case RuleAction83:
			yy = MakeASTNode(NODE_ARRAY, nil, nil, nil, currentLine) //Value:13 []
		case RuleAction84:
			yy = MakeASTNode(NODE_ARRAY, stack[stack_idx-0], nil, nil, currentLine) //Value:14 [AryItems]
		case RuleAction85:
			yy = MakeASTNode(NODE_HASH, nil, nil, nil, currentLine) //Value:15 {}
		case RuleAction86:
			yy = MakeASTNode(NODE_HASH, stack[stack_idx-0], nil, nil, currentLine) //Value:16 {HashItems}
		case RuleAction87:
			//AryItems:0
		case RuleAction88:
			stack[stack_idx-1].PushBack(stack[stack_idx-0]) //AryItems:1
		case RuleAction89:
			yy = stack[stack_idx-1] //AryItems:2
		case RuleAction90:
			stack[stack_idx-2].PushBack(stack[stack_idx-1]) //HashItems:0
		case RuleAction91:
			stack[stack_idx-2].PushBack(stack[stack_idx-0]) //HashItems:1 stack[stack_idx-0]
		case RuleAction92:
			stack[stack_idx-2].PushBack(stack[stack_idx-1]) //HashItems:2 stack[stack_idx-1]
		case RuleAction93:
			yy = stack[stack_idx-2] //HashItems:3 yy = stack[stack_idx-2]
		case RuleAction94:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //ID:0 KEYWORD.([
		case RuleAction95:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //ID:1 KEYWORD NAME
		case RuleAction96:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //ID:2 simply ID
		case RuleAction97:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //CONST:0
		case RuleAction98:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //BINOP:0
		case RuleAction99:
			//UNOP:0
		case RuleAction100:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //ASSIGN:0
		case RuleAction101:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //IVAR:0
		case RuleAction102:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //CVAR:0
		case RuleAction103:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //GLOBAL:0
		case RuleAction104:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			tempInt, _ := strconv.Atoi(buffer[begin:end])
			currentAST.value.numeric = int64(tempInt)
			yy = currentAST //NUMBER:0
		case RuleAction105:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = buffer[begin:end]
			yy = currentAST //SYMBOL:0
		case RuleAction106:
			parseStr = make([]string, 12) //STRING1:0 STRING_START
		case RuleAction107:
			parseStr = append(parseStr, "\\'") //STRING1:1 escaped \'
		case RuleAction108:
			parseStr = append(parseStr, buffer[begin:end]) //STRING1:2 content
		case RuleAction109:
			currentAST := &AST{Type: NODE_ASTVAL, line: currentLine}
			currentAST.value.str = strings.Join(parseStr, "")
			yy = currentAST //STRING1:3
		case RuleAction110:
			currentLine += 1 //EOL:0 count line

		case RuleActionPush:
			stack_idx += 1
		case RuleActionPop:
			stack_idx -= 1
		case RuleActionSet:
			stack[stack_idx] = yy

		}
	}
}

func (p *Leg) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != END_SYMBOL {
		p.buffer = append(p.buffer, END_SYMBOL)
	}

	var tree TokenTree = &tokens16{tree: make([]token16, math.MaxInt16)}
	position, depth, tokenIndex, buffer, rules := 0, 0, 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.TokenTree = tree
		if matches {
			p.TokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule Rule, begin int) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	matchDot := func() bool {
		if buffer[position] != END_SYMBOL {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
	    if buffer[position] == c {
	        position++
	        return true
	    }
	    return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
	    if c := buffer[position]; c >= lower && c <= upper {
	        position++
	        return true
	    }
	    return false
	}*/

	rules = [...]func() bool{
		nil,
		/* 0 Root <- <(Stmts EOF Action0)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 1
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position1 := position
				depth++
				if !rules[RuleStmts]() {
					goto l0
				}
				variableIdx = 0
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleEOF]() {
					goto l0
				}
				if !rules[RuleAction0]() {
					goto l0
				}
				depth--
				add(RuleRoot, position1)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Stmts <- <((SEP* _ Stmt Comment? Action1 ((SEP _ Stmt Comment? Action2) / (SEP _ Comment))* SEP? Action3) / (SEP+ Action4))> */
		func() bool {
			position2, tokenIndex2, depth2 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position3 := position
				depth++
				{

					position4, tokenIndex4, depth4 := position, tokenIndex, depth
				l6:
					{

						position7, tokenIndex7, depth7 := position, tokenIndex, depth
						if !rules[RuleSEP]() {
							goto l7
						}
						goto l6
					l7:
						position, tokenIndex, depth = position7, tokenIndex7, depth7
					}
					if !rules[Rule_]() {
						goto l5
					}
					if !rules[RuleStmt]() {
						goto l5
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					{

						position8, tokenIndex8, depth8 := position, tokenIndex, depth
						if !rules[RuleComment]() {
							goto l8
						}
						goto l9
					l8:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
					}
				l9:
					if !rules[RuleAction1]() {
						goto l5
					}
				l10:
					{

						position11, tokenIndex11, depth11 := position, tokenIndex, depth
						{

							position12, tokenIndex12, depth12 := position, tokenIndex, depth
							if !rules[RuleSEP]() {
								goto l13
							}
							if !rules[Rule_]() {
								goto l13
							}
							if !rules[RuleStmt]() {
								goto l13
							}
							variableIdx = 0
							for i := 0; i < variableIdx; i++ {
								add(RuleActionPop, position)
							}
							add(RuleActionSet, position)
							for i := 0; i < variableIdx; i++ {
								add(RuleActionPush, position)
							}
							{

								position14, tokenIndex14, depth14 := position, tokenIndex, depth
								if !rules[RuleComment]() {
									goto l14
								}
								goto l15
							l14:
								position, tokenIndex, depth = position14, tokenIndex14, depth14
							}
						l15:
							if !rules[RuleAction2]() {
								goto l13
							}
							goto l12
						l13:
							position, tokenIndex, depth = position12, tokenIndex12, depth12
							if !rules[RuleSEP]() {
								goto l11
							}
							if !rules[Rule_]() {
								goto l11
							}
							if !rules[RuleComment]() {
								goto l11
							}
						}
					l12:
						goto l10
					l11:
						position, tokenIndex, depth = position11, tokenIndex11, depth11
					}
					{

						position16, tokenIndex16, depth16 := position, tokenIndex, depth
						if !rules[RuleSEP]() {
							goto l16
						}
						goto l17
					l16:
						position, tokenIndex, depth = position16, tokenIndex16, depth16
					}
				l17:
					if !rules[RuleAction3]() {
						goto l5
					}
					goto l4
				l5:
					position, tokenIndex, depth = position4, tokenIndex4, depth4
					if !rules[RuleSEP]() {
						goto l2
					}
				l18:
					{

						position19, tokenIndex19, depth19 := position, tokenIndex, depth
						if !rules[RuleSEP]() {
							goto l19
						}
						goto l18
					l19:
						position, tokenIndex, depth = position19, tokenIndex19, depth19
					}
					if !rules[RuleAction4]() {
						goto l2
					}
				}
			l4:
				depth--
				add(RuleStmts, position3)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l2:
			position, tokenIndex, depth = position2, tokenIndex2, depth2
			return false
		},
		/* 2 OptStmts <- <(Stmts / (_ SEP? Action5))> */
		func() bool {
			position20, tokenIndex20, depth20 := position, tokenIndex, depth
			{

				position21 := position
				depth++
				{

					position22, tokenIndex22, depth22 := position, tokenIndex, depth
					if !rules[RuleStmts]() {
						goto l23
					}
					goto l22
				l23:
					position, tokenIndex, depth = position22, tokenIndex22, depth22
					if !rules[Rule_]() {
						goto l20
					}
					{

						position24, tokenIndex24, depth24 := position, tokenIndex, depth
						if !rules[RuleSEP]() {
							goto l24
						}
						goto l25
					l24:
						position, tokenIndex, depth = position24, tokenIndex24, depth24
					}
				l25:
					if !rules[RuleAction5]() {
						goto l20
					}
				}
			l22:
				depth--
				add(RuleOptStmts, position21)
			}
			return true
		l20:
			position, tokenIndex, depth = position20, tokenIndex20, depth20
			return false
		},
		/* 3 Stmt <- <(While / Until / If / Unless / Def / Class / Module / Expr)> */
		func() bool {
			position26, tokenIndex26, depth26 := position, tokenIndex, depth
			{

				position27 := position
				depth++
				{

					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if !rules[RuleWhile]() {
						goto l29
					}
					goto l28
				l29:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !rules[RuleUntil]() {
						goto l30
					}
					goto l28
				l30:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !rules[RuleIf]() {
						goto l31
					}
					goto l28
				l31:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !rules[RuleUnless]() {
						goto l32
					}
					goto l28
				l32:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !rules[RuleDef]() {
						goto l33
					}
					goto l28
				l33:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !rules[RuleClass]() {
						goto l34
					}
					goto l28
				l34:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !rules[RuleModule]() {
						goto l35
					}
					goto l28
				l35:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if !rules[RuleExpr]() {
						goto l26
					}
				}
			l28:
				depth--
				add(RuleStmt, position27)
			}
			return true
		l26:
			position, tokenIndex, depth = position26, tokenIndex26, depth26
			return false
		},
		/* 4 Expr <- <(Assign / AsgnCall / UnaryOp / BinOp / SpecCall / Call / Range / Yield / Return / Break / Value)> */
		func() bool {
			position36, tokenIndex36, depth36 := position, tokenIndex, depth
			{

				position37 := position
				depth++
				{

					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if !rules[RuleAssign]() {
						goto l39
					}
					goto l38
				l39:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleAsgnCall]() {
						goto l40
					}
					goto l38
				l40:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleUnaryOp]() {
						goto l41
					}
					goto l38
				l41:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleBinOp]() {
						goto l42
					}
					goto l38
				l42:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleSpecCall]() {
						goto l43
					}
					goto l38
				l43:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleCall]() {
						goto l44
					}
					goto l38
				l44:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleRange]() {
						goto l45
					}
					goto l38
				l45:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleYield]() {
						goto l46
					}
					goto l38
				l46:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleReturn]() {
						goto l47
					}
					goto l38
				l47:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleBreak]() {
						goto l48
					}
					goto l38
				l48:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !rules[RuleValue]() {
						goto l36
					}
				}
			l38:
				depth--
				add(RuleExpr, position37)
			}
			return true
		l36:
			position, tokenIndex, depth = position36, tokenIndex36, depth36
			return false
		},
		/* 5 Comment <- <(_ '#' (!EOL .)*)> */
		func() bool {
			position49, tokenIndex49, depth49 := position, tokenIndex, depth
			{

				position50 := position
				depth++
				if !rules[Rule_]() {
					goto l49
				}
				if buffer[position] != rune('#') {
					goto l49
				}
				position++
			l51:
				{

					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					{

						position53, tokenIndex53, depth53 := position, tokenIndex, depth
						if !rules[RuleEOL]() {
							goto l53
						}
						goto l52
					l53:
						position, tokenIndex, depth = position53, tokenIndex53, depth53
					}
					if !matchDot() {
						goto l52
					}
					goto l51
				l52:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
				}
				depth--
				add(RuleComment, position50)
			}
			return true
		l49:
			position, tokenIndex, depth = position49, tokenIndex49, depth49
			return false
		},
		/* 6 Call <- <(Action6 (Value '.')? (Message '.' Action7)* Message _ Block? Action8)> */
		func() bool {
			position54, tokenIndex54, depth54 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 4
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position55 := position
				depth++
				if !rules[RuleAction6]() {
					goto l54
				}
				{

					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if !rules[RuleValue]() {
						goto l56
					}
					variableIdx = 3
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('.') {
						goto l56
					}
					position++
					goto l57
				l56:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
				}
			l57:
			l58:
				{

					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					if !rules[RuleMessage]() {
						goto l59
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('.') {
						goto l59
					}
					position++
					if !rules[RuleAction7]() {
						goto l59
					}
					goto l58
				l59:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
				}
				if !rules[RuleMessage]() {
					goto l54
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l54
				}
				{

					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					if !rules[RuleBlock]() {
						goto l60
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					goto l61
				l60:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
				}
			l61:
				if !rules[RuleAction8]() {
					goto l54
				}
				depth--
				add(RuleCall, position55)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l54:
			position, tokenIndex, depth = position54, tokenIndex54, depth54
			return false
		},
		/* 7 AsgnCall <- <(Action9 (Value '.')? (Message '.' Action10)* ID _ ASSIGN _ Stmt Action11)> */
		func() bool {
			position62, tokenIndex62, depth62 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 5
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position63 := position
				depth++
				if !rules[RuleAction9]() {
					goto l62
				}
				{

					position64, tokenIndex64, depth64 := position, tokenIndex, depth
					if !rules[RuleValue]() {
						goto l64
					}
					variableIdx = 4
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('.') {
						goto l64
					}
					position++
					goto l65
				l64:
					position, tokenIndex, depth = position64, tokenIndex64, depth64
				}
			l65:
			l66:
				{

					position67, tokenIndex67, depth67 := position, tokenIndex, depth
					if !rules[RuleMessage]() {
						goto l67
					}
					variableIdx = 3
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('.') {
						goto l67
					}
					position++
					if !rules[RuleAction10]() {
						goto l67
					}
					goto l66
				l67:
					position, tokenIndex, depth = position67, tokenIndex67, depth67
				}
				if !rules[RuleID]() {
					goto l62
				}
				variableIdx = 2
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l62
				}
				if !rules[RuleASSIGN]() {
					goto l62
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l62
				}
				if !rules[RuleStmt]() {
					goto l62
				}
				variableIdx = 0
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleAction11]() {
					goto l62
				}
				depth--
				add(RuleAsgnCall, position63)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l62:
			position, tokenIndex, depth = position62, tokenIndex62, depth62
			return false
		},
		/* 8 Receiver <- <(Action12 ((Call Action13) / (Value Action14)))> */
		func() bool {
			position68, tokenIndex68, depth68 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 1
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position69 := position
				depth++
				if !rules[RuleAction12]() {
					goto l68
				}
				{

					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if !rules[RuleCall]() {
						goto l71
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction13]() {
						goto l71
					}
					goto l70
				l71:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
					if !rules[RuleValue]() {
						goto l68
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction14]() {
						goto l68
					}
				}
			l70:
				depth--
				add(RuleReceiver, position69)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l68:
			position, tokenIndex, depth = position68, tokenIndex68, depth68
			return false
		},
		/* 9 SpecCall <- <((Receiver '[' Args ']' _ ASSIGN _ Stmt Action15) / (Receiver '[' Args ']' Action16))> */
		func() bool {
			position72, tokenIndex72, depth72 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position73 := position
				depth++
				{

					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if !rules[RuleReceiver]() {
						goto l75
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('[') {
						goto l75
					}
					position++
					if !rules[RuleArgs]() {
						goto l75
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune(']') {
						goto l75
					}
					position++
					if !rules[Rule_]() {
						goto l75
					}
					if !rules[RuleASSIGN]() {
						goto l75
					}
					if !rules[Rule_]() {
						goto l75
					}
					if !rules[RuleStmt]() {
						goto l75
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction15]() {
						goto l75
					}
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					if !rules[RuleReceiver]() {
						goto l72
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('[') {
						goto l72
					}
					position++
					if !rules[RuleArgs]() {
						goto l72
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune(']') {
						goto l72
					}
					position++
					if !rules[RuleAction16]() {
						goto l72
					}
				}
			l74:
				depth--
				add(RuleSpecCall, position73)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l72:
			position, tokenIndex, depth = position72, tokenIndex72, depth72
			return false
		},
		/* 10 BinOp <- <((SpecCall / Receiver) _ (('&' '&' _ Expr Action17) / ('|' '|' _ Expr Action18) / ('+' _ Expr Action19) / ('-' _ Expr Action20) / ('<' _ Expr Action21) / (BINOP _ Expr Action22)))> */
		func() bool {
			position76, tokenIndex76, depth76 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position77 := position
				depth++
				{

					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if !rules[RuleSpecCall]() {
						goto l79
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if !rules[RuleReceiver]() {
						goto l76
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
				}
			l78:
				if !rules[Rule_]() {
					goto l76
				}
				{

					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('&') {
						goto l81
					}
					position++
					if buffer[position] != rune('&') {
						goto l81
					}
					position++
					if !rules[Rule_]() {
						goto l81
					}
					if !rules[RuleExpr]() {
						goto l81
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction17]() {
						goto l81
					}
					goto l80
				l81:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('|') {
						goto l82
					}
					position++
					if buffer[position] != rune('|') {
						goto l82
					}
					position++
					if !rules[Rule_]() {
						goto l82
					}
					if !rules[RuleExpr]() {
						goto l82
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction18]() {
						goto l82
					}
					goto l80
				l82:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('+') {
						goto l83
					}
					position++
					if !rules[Rule_]() {
						goto l83
					}
					if !rules[RuleExpr]() {
						goto l83
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction19]() {
						goto l83
					}
					goto l80
				l83:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('-') {
						goto l84
					}
					position++
					if !rules[Rule_]() {
						goto l84
					}
					if !rules[RuleExpr]() {
						goto l84
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction20]() {
						goto l84
					}
					goto l80
				l84:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('<') {
						goto l85
					}
					position++
					if !rules[Rule_]() {
						goto l85
					}
					if !rules[RuleExpr]() {
						goto l85
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction21]() {
						goto l85
					}
					goto l80
				l85:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if !rules[RuleBINOP]() {
						goto l76
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l76
					}
					if !rules[RuleExpr]() {
						goto l76
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction22]() {
						goto l76
					}
				}
			l80:
				depth--
				add(RuleBinOp, position77)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l76:
			position, tokenIndex, depth = position76, tokenIndex76, depth76
			return false
		},
		/* 11 UnaryOp <- <(('-' Expr Action23) / ('!' Expr Action24))> */
		func() bool {
			position86, tokenIndex86, depth86 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 1
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position87 := position
				depth++
				{

					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l89
					}
					position++
					if !rules[RuleExpr]() {
						goto l89
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction23]() {
						goto l89
					}
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if buffer[position] != rune('!') {
						goto l86
					}
					position++
					if !rules[RuleExpr]() {
						goto l86
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction24]() {
						goto l86
					}
				}
			l88:
				depth--
				add(RuleUnaryOp, position87)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l86:
			position, tokenIndex, depth = position86, tokenIndex86, depth86
			return false
		},
		/* 12 Message <- <(ID Action25 (('(' Args? ')') / (SPACE Args))? Action26)> */
		func() bool {
			position90, tokenIndex90, depth90 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position91 := position
				depth++
				if !rules[RuleID]() {
					goto l90
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleAction25]() {
					goto l90
				}
				{

					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					{

						position94, tokenIndex94, depth94 := position, tokenIndex, depth
						if buffer[position] != rune('(') {
							goto l95
						}
						position++
						{

							position96, tokenIndex96, depth96 := position, tokenIndex, depth
							if !rules[RuleArgs]() {
								goto l96
							}
							variableIdx = 0
							for i := 0; i < variableIdx; i++ {
								add(RuleActionPop, position)
							}
							add(RuleActionSet, position)
							for i := 0; i < variableIdx; i++ {
								add(RuleActionPush, position)
							}
							goto l97
						l96:
							position, tokenIndex, depth = position96, tokenIndex96, depth96
						}
					l97:
						if buffer[position] != rune(')') {
							goto l95
						}
						position++
						goto l94
					l95:
						position, tokenIndex, depth = position94, tokenIndex94, depth94
						if !rules[RuleSPACE]() {
							goto l92
						}
						if !rules[RuleArgs]() {
							goto l92
						}
						variableIdx = 0
						for i := 0; i < variableIdx; i++ {
							add(RuleActionPop, position)
						}
						add(RuleActionSet, position)
						for i := 0; i < variableIdx; i++ {
							add(RuleActionPush, position)
						}
					}
				l94:
					goto l93
				l92:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
				}
			l93:
				if !rules[RuleAction26]() {
					goto l90
				}
				depth--
				add(RuleMessage, position91)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l90:
			position, tokenIndex, depth = position90, tokenIndex90, depth90
			return false
		},
		/* 13 Args <- <(_ Expr _ Action27 (',' _ Expr _ Action28)* (',' _ '*' Expr _ Action29)? Action30)> */
		func() bool {
			position98, tokenIndex98, depth98 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position99 := position
				depth++
				if !rules[Rule_]() {
					goto l98
				}
				if !rules[RuleExpr]() {
					goto l98
				}
				variableIdx = 2
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l98
				}
				if !rules[RuleAction27]() {
					goto l98
				}
			l100:
				{

					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l101
					}
					position++
					if !rules[Rule_]() {
						goto l101
					}
					if !rules[RuleExpr]() {
						goto l101
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l101
					}
					if !rules[RuleAction28]() {
						goto l101
					}
					goto l100
				l101:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
				}
				{

					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l102
					}
					position++
					if !rules[Rule_]() {
						goto l102
					}
					if buffer[position] != rune('*') {
						goto l102
					}
					position++
					if !rules[RuleExpr]() {
						goto l102
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l102
					}
					if !rules[RuleAction29]() {
						goto l102
					}
					goto l103
				l102:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
				}
			l103:
				if !rules[RuleAction30]() {
					goto l98
				}
				depth--
				add(RuleArgs, position99)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l98:
			position, tokenIndex, depth = position98, tokenIndex98, depth98
			return false
		},
		/* 14 Block <- <(('d' 'o' SEP _ OptStmts _ ('e' 'n' 'd') Action31) / ('d' 'o' _ '|' Params '|' SEP _ OptStmts _ ('e' 'n' 'd') Action32))> */
		func() bool {
			position104, tokenIndex104, depth104 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position105 := position
				depth++
				{

					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l107
					}
					position++
					if buffer[position] != rune('o') {
						goto l107
					}
					position++
					if !rules[RuleSEP]() {
						goto l107
					}
					if !rules[Rule_]() {
						goto l107
					}
					if !rules[RuleOptStmts]() {
						goto l107
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l107
					}
					if buffer[position] != rune('e') {
						goto l107
					}
					position++
					if buffer[position] != rune('n') {
						goto l107
					}
					position++
					if buffer[position] != rune('d') {
						goto l107
					}
					position++
					if !rules[RuleAction31]() {
						goto l107
					}
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if buffer[position] != rune('d') {
						goto l104
					}
					position++
					if buffer[position] != rune('o') {
						goto l104
					}
					position++
					if !rules[Rule_]() {
						goto l104
					}
					if buffer[position] != rune('|') {
						goto l104
					}
					position++
					if !rules[RuleParams]() {
						goto l104
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('|') {
						goto l104
					}
					position++
					if !rules[RuleSEP]() {
						goto l104
					}
					if !rules[Rule_]() {
						goto l104
					}
					if !rules[RuleOptStmts]() {
						goto l104
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l104
					}
					if buffer[position] != rune('e') {
						goto l104
					}
					position++
					if buffer[position] != rune('n') {
						goto l104
					}
					position++
					if buffer[position] != rune('d') {
						goto l104
					}
					position++
					if !rules[RuleAction32]() {
						goto l104
					}
				}
			l106:
				depth--
				add(RuleBlock, position105)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l104:
			position, tokenIndex, depth = position104, tokenIndex104, depth104
			return false
		},
		/* 15 Assign <- <((ID _ ASSIGN _ Stmt Action33) / (CONST _ ASSIGN _ Stmt Action34) / (IVAR _ ASSIGN _ Stmt Action35) / (CVAR _ ASSIGN _ Stmt Action36) / (GLOBAL _ ASSIGN _ Stmt Action37))> */
		func() bool {
			position108, tokenIndex108, depth108 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position109 := position
				depth++
				{

					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if !rules[RuleID]() {
						goto l111
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l111
					}
					if !rules[RuleASSIGN]() {
						goto l111
					}
					if !rules[Rule_]() {
						goto l111
					}
					if !rules[RuleStmt]() {
						goto l111
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction33]() {
						goto l111
					}
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if !rules[RuleCONST]() {
						goto l112
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l112
					}
					if !rules[RuleASSIGN]() {
						goto l112
					}
					if !rules[Rule_]() {
						goto l112
					}
					if !rules[RuleStmt]() {
						goto l112
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction34]() {
						goto l112
					}
					goto l110
				l112:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if !rules[RuleIVAR]() {
						goto l113
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l113
					}
					if !rules[RuleASSIGN]() {
						goto l113
					}
					if !rules[Rule_]() {
						goto l113
					}
					if !rules[RuleStmt]() {
						goto l113
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction35]() {
						goto l113
					}
					goto l110
				l113:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if !rules[RuleCVAR]() {
						goto l114
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l114
					}
					if !rules[RuleASSIGN]() {
						goto l114
					}
					if !rules[Rule_]() {
						goto l114
					}
					if !rules[RuleStmt]() {
						goto l114
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction36]() {
						goto l114
					}
					goto l110
				l114:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if !rules[RuleGLOBAL]() {
						goto l108
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l108
					}
					if !rules[RuleASSIGN]() {
						goto l108
					}
					if !rules[Rule_]() {
						goto l108
					}
					if !rules[RuleStmt]() {
						goto l108
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction37]() {
						goto l108
					}
				}
			l110:
				depth--
				add(RuleAssign, position109)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l108:
			position, tokenIndex, depth = position108, tokenIndex108, depth108
			return false
		},
		/* 16 While <- <('w' 'h' 'i' 'l' 'e' SPACE Expr SEP Stmts _ ('e' 'n' 'd') Action38)> */
		func() bool {
			position115, tokenIndex115, depth115 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position116 := position
				depth++
				if buffer[position] != rune('w') {
					goto l115
				}
				position++
				if buffer[position] != rune('h') {
					goto l115
				}
				position++
				if buffer[position] != rune('i') {
					goto l115
				}
				position++
				if buffer[position] != rune('l') {
					goto l115
				}
				position++
				if buffer[position] != rune('e') {
					goto l115
				}
				position++
				if !rules[RuleSPACE]() {
					goto l115
				}
				if !rules[RuleExpr]() {
					goto l115
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleSEP]() {
					goto l115
				}
				if !rules[RuleStmts]() {
					goto l115
				}
				variableIdx = 0
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l115
				}
				if buffer[position] != rune('e') {
					goto l115
				}
				position++
				if buffer[position] != rune('n') {
					goto l115
				}
				position++
				if buffer[position] != rune('d') {
					goto l115
				}
				position++
				if !rules[RuleAction38]() {
					goto l115
				}
				depth--
				add(RuleWhile, position116)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l115:
			position, tokenIndex, depth = position115, tokenIndex115, depth115
			return false
		},
		/* 17 Until <- <('u' 'n' 't' 'i' 'l' SPACE Expr SEP Stmts _ ('e' 'n' 'd') Action39)> */
		func() bool {
			position117, tokenIndex117, depth117 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position118 := position
				depth++
				if buffer[position] != rune('u') {
					goto l117
				}
				position++
				if buffer[position] != rune('n') {
					goto l117
				}
				position++
				if buffer[position] != rune('t') {
					goto l117
				}
				position++
				if buffer[position] != rune('i') {
					goto l117
				}
				position++
				if buffer[position] != rune('l') {
					goto l117
				}
				position++
				if !rules[RuleSPACE]() {
					goto l117
				}
				if !rules[RuleExpr]() {
					goto l117
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleSEP]() {
					goto l117
				}
				if !rules[RuleStmts]() {
					goto l117
				}
				variableIdx = 0
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l117
				}
				if buffer[position] != rune('e') {
					goto l117
				}
				position++
				if buffer[position] != rune('n') {
					goto l117
				}
				position++
				if buffer[position] != rune('d') {
					goto l117
				}
				position++
				if !rules[RuleAction39]() {
					goto l117
				}
				depth--
				add(RuleUntil, position118)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l117:
			position, tokenIndex, depth = position117, tokenIndex117, depth117
			return false
		},
		/* 18 If <- <(('i' 'f' SPACE Expr SEP Action40 Stmts _ Else? ('e' 'n' 'd') Action41) / (Expr _ ('i' 'f') _ Expr Action42))> */
		func() bool {
			position119, tokenIndex119, depth119 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position120 := position
				depth++
				{

					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l122
					}
					position++
					if buffer[position] != rune('f') {
						goto l122
					}
					position++
					if !rules[RuleSPACE]() {
						goto l122
					}
					if !rules[RuleExpr]() {
						goto l122
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleSEP]() {
						goto l122
					}
					if !rules[RuleAction40]() {
						goto l122
					}
					if !rules[RuleStmts]() {
						goto l122
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l122
					}
					{

						position123, tokenIndex123, depth123 := position, tokenIndex, depth
						if !rules[RuleElse]() {
							goto l123
						}
						variableIdx = 0
						for i := 0; i < variableIdx; i++ {
							add(RuleActionPop, position)
						}
						add(RuleActionSet, position)
						for i := 0; i < variableIdx; i++ {
							add(RuleActionPush, position)
						}
						goto l124
					l123:
						position, tokenIndex, depth = position123, tokenIndex123, depth123
					}
				l124:
					if buffer[position] != rune('e') {
						goto l122
					}
					position++
					if buffer[position] != rune('n') {
						goto l122
					}
					position++
					if buffer[position] != rune('d') {
						goto l122
					}
					position++
					if !rules[RuleAction41]() {
						goto l122
					}
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					if !rules[RuleExpr]() {
						goto l119
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l119
					}
					if buffer[position] != rune('i') {
						goto l119
					}
					position++
					if buffer[position] != rune('f') {
						goto l119
					}
					position++
					if !rules[Rule_]() {
						goto l119
					}
					if !rules[RuleExpr]() {
						goto l119
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction42]() {
						goto l119
					}
				}
			l121:
				depth--
				add(RuleIf, position120)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l119:
			position, tokenIndex, depth = position119, tokenIndex119, depth119
			return false
		},
		/* 19 Unless <- <(('u' 'n' 'l' 'e' 's' 's' SPACE Expr SEP Action43 Stmts _ Else? ('e' 'n' 'd') Action44) / (Expr _ ('u' 'n' 'l' 'e' 's' 's') _ Expr Action45))> */
		func() bool {
			position125, tokenIndex125, depth125 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position126 := position
				depth++
				{

					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l128
					}
					position++
					if buffer[position] != rune('n') {
						goto l128
					}
					position++
					if buffer[position] != rune('l') {
						goto l128
					}
					position++
					if buffer[position] != rune('e') {
						goto l128
					}
					position++
					if buffer[position] != rune('s') {
						goto l128
					}
					position++
					if buffer[position] != rune('s') {
						goto l128
					}
					position++
					if !rules[RuleSPACE]() {
						goto l128
					}
					if !rules[RuleExpr]() {
						goto l128
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleSEP]() {
						goto l128
					}
					if !rules[RuleAction43]() {
						goto l128
					}
					if !rules[RuleStmts]() {
						goto l128
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l128
					}
					{

						position129, tokenIndex129, depth129 := position, tokenIndex, depth
						if !rules[RuleElse]() {
							goto l129
						}
						variableIdx = 0
						for i := 0; i < variableIdx; i++ {
							add(RuleActionPop, position)
						}
						add(RuleActionSet, position)
						for i := 0; i < variableIdx; i++ {
							add(RuleActionPush, position)
						}
						goto l130
					l129:
						position, tokenIndex, depth = position129, tokenIndex129, depth129
					}
				l130:
					if buffer[position] != rune('e') {
						goto l128
					}
					position++
					if buffer[position] != rune('n') {
						goto l128
					}
					position++
					if buffer[position] != rune('d') {
						goto l128
					}
					position++
					if !rules[RuleAction44]() {
						goto l128
					}
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if !rules[RuleExpr]() {
						goto l125
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l125
					}
					if buffer[position] != rune('u') {
						goto l125
					}
					position++
					if buffer[position] != rune('n') {
						goto l125
					}
					position++
					if buffer[position] != rune('l') {
						goto l125
					}
					position++
					if buffer[position] != rune('e') {
						goto l125
					}
					position++
					if buffer[position] != rune('s') {
						goto l125
					}
					position++
					if buffer[position] != rune('s') {
						goto l125
					}
					position++
					if !rules[Rule_]() {
						goto l125
					}
					if !rules[RuleExpr]() {
						goto l125
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction45]() {
						goto l125
					}
				}
			l127:
				depth--
				add(RuleUnless, position126)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l125:
			position, tokenIndex, depth = position125, tokenIndex125, depth125
			return false
		},
		/* 20 Else <- <('e' 'l' 's' 'e' SEP _ Stmts _ Action46)> */
		func() bool {
			position131, tokenIndex131, depth131 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 1
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position132 := position
				depth++
				if buffer[position] != rune('e') {
					goto l131
				}
				position++
				if buffer[position] != rune('l') {
					goto l131
				}
				position++
				if buffer[position] != rune('s') {
					goto l131
				}
				position++
				if buffer[position] != rune('e') {
					goto l131
				}
				position++
				if !rules[RuleSEP]() {
					goto l131
				}
				if !rules[Rule_]() {
					goto l131
				}
				if !rules[RuleStmts]() {
					goto l131
				}
				variableIdx = 0
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l131
				}
				if !rules[RuleAction46]() {
					goto l131
				}
				depth--
				add(RuleElse, position132)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l131:
			position, tokenIndex, depth = position131, tokenIndex131, depth131
			return false
		},
		/* 21 Method <- <((ID '.' METHOD Action47) / (Value '.' METHOD Action48) / (METHOD Action49))> */
		func() bool {
			position133, tokenIndex133, depth133 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position134 := position
				depth++
				{

					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if !rules[RuleID]() {
						goto l136
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('.') {
						goto l136
					}
					position++
					if !rules[RuleMETHOD]() {
						goto l136
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction47]() {
						goto l136
					}
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if !rules[RuleValue]() {
						goto l137
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune('.') {
						goto l137
					}
					position++
					if !rules[RuleMETHOD]() {
						goto l137
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction48]() {
						goto l137
					}
					goto l135
				l137:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if !rules[RuleMETHOD]() {
						goto l133
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction49]() {
						goto l133
					}
				}
			l135:
				depth--
				add(RuleMethod, position134)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l133:
			position, tokenIndex, depth = position133, tokenIndex133, depth133
			return false
		},
		/* 22 Def <- <('d' 'e' 'f' SPACE Method Action50 (_ '(' Params? ')')? SEP OptStmts _ ('e' 'n' 'd') Action51)> */
		func() bool {
			position138, tokenIndex138, depth138 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position139 := position
				depth++
				if buffer[position] != rune('d') {
					goto l138
				}
				position++
				if buffer[position] != rune('e') {
					goto l138
				}
				position++
				if buffer[position] != rune('f') {
					goto l138
				}
				position++
				if !rules[RuleSPACE]() {
					goto l138
				}
				if !rules[RuleMethod]() {
					goto l138
				}
				variableIdx = 2
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleAction50]() {
					goto l138
				}
				{

					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if !rules[Rule_]() {
						goto l140
					}
					if buffer[position] != rune('(') {
						goto l140
					}
					position++
					{

						position142, tokenIndex142, depth142 := position, tokenIndex, depth
						if !rules[RuleParams]() {
							goto l142
						}
						variableIdx = 1
						for i := 0; i < variableIdx; i++ {
							add(RuleActionPop, position)
						}
						add(RuleActionSet, position)
						for i := 0; i < variableIdx; i++ {
							add(RuleActionPush, position)
						}
						goto l143
					l142:
						position, tokenIndex, depth = position142, tokenIndex142, depth142
					}
				l143:
					if buffer[position] != rune(')') {
						goto l140
					}
					position++
					goto l141
				l140:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
				}
			l141:
				if !rules[RuleSEP]() {
					goto l138
				}
				if !rules[RuleOptStmts]() {
					goto l138
				}
				variableIdx = 0
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l138
				}
				if buffer[position] != rune('e') {
					goto l138
				}
				position++
				if buffer[position] != rune('n') {
					goto l138
				}
				position++
				if buffer[position] != rune('d') {
					goto l138
				}
				position++
				if !rules[RuleAction51]() {
					goto l138
				}
				depth--
				add(RuleDef, position139)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l138:
			position, tokenIndex, depth = position138, tokenIndex138, depth138
			return false
		},
		/* 23 Params <- <(Param Action52 (',' Param Action53)* Action54)> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position145 := position
				depth++
				if !rules[RuleParam]() {
					goto l144
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleAction52]() {
					goto l144
				}
			l146:
				{

					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l147
					}
					position++
					if !rules[RuleParam]() {
						goto l147
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction53]() {
						goto l147
					}
					goto l146
				l147:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
				}
				if !rules[RuleAction54]() {
					goto l144
				}
				depth--
				add(RuleParams, position145)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 24 Param <- <((_ ID _ '=' _ Expr Action55) / (_ ID _ Action56) / (_ '*' ID _ Action57))> */
		func() bool {
			position148, tokenIndex148, depth148 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position149 := position
				depth++
				{

					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if !rules[Rule_]() {
						goto l151
					}
					if !rules[RuleID]() {
						goto l151
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l151
					}
					if buffer[position] != rune('=') {
						goto l151
					}
					position++
					if !rules[Rule_]() {
						goto l151
					}
					if !rules[RuleExpr]() {
						goto l151
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction55]() {
						goto l151
					}
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if !rules[Rule_]() {
						goto l152
					}
					if !rules[RuleID]() {
						goto l152
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l152
					}
					if !rules[RuleAction56]() {
						goto l152
					}
					goto l150
				l152:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if !rules[Rule_]() {
						goto l148
					}
					if buffer[position] != rune('*') {
						goto l148
					}
					position++
					if !rules[RuleID]() {
						goto l148
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l148
					}
					if !rules[RuleAction57]() {
						goto l148
					}
				}
			l150:
				depth--
				add(RuleParam, position149)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l148:
			position, tokenIndex, depth = position148, tokenIndex148, depth148
			return false
		},
		/* 25 Class <- <('c' 'l' 'a' 's' 's' SPACE CONST Action58 (_ '<' _ CONST)? SEP OptStmts _ ('e' 'n' 'd') Action59)> */
		func() bool {
			position153, tokenIndex153, depth153 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position154 := position
				depth++
				if buffer[position] != rune('c') {
					goto l153
				}
				position++
				if buffer[position] != rune('l') {
					goto l153
				}
				position++
				if buffer[position] != rune('a') {
					goto l153
				}
				position++
				if buffer[position] != rune('s') {
					goto l153
				}
				position++
				if buffer[position] != rune('s') {
					goto l153
				}
				position++
				if !rules[RuleSPACE]() {
					goto l153
				}
				if !rules[RuleCONST]() {
					goto l153
				}
				variableIdx = 2
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleAction58]() {
					goto l153
				}
				{

					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if !rules[Rule_]() {
						goto l155
					}
					if buffer[position] != rune('<') {
						goto l155
					}
					position++
					if !rules[Rule_]() {
						goto l155
					}
					if !rules[RuleCONST]() {
						goto l155
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					goto l156
				l155:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
				}
			l156:
				if !rules[RuleSEP]() {
					goto l153
				}
				if !rules[RuleOptStmts]() {
					goto l153
				}
				variableIdx = 0
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l153
				}
				if buffer[position] != rune('e') {
					goto l153
				}
				position++
				if buffer[position] != rune('n') {
					goto l153
				}
				position++
				if buffer[position] != rune('d') {
					goto l153
				}
				position++
				if !rules[RuleAction59]() {
					goto l153
				}
				depth--
				add(RuleClass, position154)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l153:
			position, tokenIndex, depth = position153, tokenIndex153, depth153
			return false
		},
		/* 26 Module <- <('m' 'o' 'd' 'u' 'l' 'e' SPACE CONST SEP OptStmts _ ('e' 'n' 'd') Action60)> */
		func() bool {
			position157, tokenIndex157, depth157 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position158 := position
				depth++
				if buffer[position] != rune('m') {
					goto l157
				}
				position++
				if buffer[position] != rune('o') {
					goto l157
				}
				position++
				if buffer[position] != rune('d') {
					goto l157
				}
				position++
				if buffer[position] != rune('u') {
					goto l157
				}
				position++
				if buffer[position] != rune('l') {
					goto l157
				}
				position++
				if buffer[position] != rune('e') {
					goto l157
				}
				position++
				if !rules[RuleSPACE]() {
					goto l157
				}
				if !rules[RuleCONST]() {
					goto l157
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleSEP]() {
					goto l157
				}
				if !rules[RuleOptStmts]() {
					goto l157
				}
				variableIdx = 0
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l157
				}
				if buffer[position] != rune('e') {
					goto l157
				}
				position++
				if buffer[position] != rune('n') {
					goto l157
				}
				position++
				if buffer[position] != rune('d') {
					goto l157
				}
				position++
				if !rules[RuleAction60]() {
					goto l157
				}
				depth--
				add(RuleModule, position158)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l157:
			position, tokenIndex, depth = position157, tokenIndex157, depth157
			return false
		},
		/* 27 Range <- <((Receiver _ ('.' '.') _ Expr Action61) / (Receiver _ ('.' '.' '.') _ Expr Action62))> */
		func() bool {
			position159, tokenIndex159, depth159 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position160 := position
				depth++
				{

					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if !rules[RuleReceiver]() {
						goto l162
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l162
					}
					if buffer[position] != rune('.') {
						goto l162
					}
					position++
					if buffer[position] != rune('.') {
						goto l162
					}
					position++
					if !rules[Rule_]() {
						goto l162
					}
					if !rules[RuleExpr]() {
						goto l162
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction61]() {
						goto l162
					}
					goto l161
				l162:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
					if !rules[RuleReceiver]() {
						goto l159
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l159
					}
					if buffer[position] != rune('.') {
						goto l159
					}
					position++
					if buffer[position] != rune('.') {
						goto l159
					}
					position++
					if buffer[position] != rune('.') {
						goto l159
					}
					position++
					if !rules[Rule_]() {
						goto l159
					}
					if !rules[RuleExpr]() {
						goto l159
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction62]() {
						goto l159
					}
				}
			l161:
				depth--
				add(RuleRange, position160)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l159:
			position, tokenIndex, depth = position159, tokenIndex159, depth159
			return false
		},
		/* 28 Yield <- <(('y' 'i' 'e' 'l' 'd' SPACE AryItems Action63) / ('y' 'i' 'e' 'l' 'd' '(' AryItems ')' Action64) / ('y' 'i' 'e' 'l' 'd' Action65))> */
		func() bool {
			position163, tokenIndex163, depth163 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 1
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position164 := position
				depth++
				{

					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l166
					}
					position++
					if buffer[position] != rune('i') {
						goto l166
					}
					position++
					if buffer[position] != rune('e') {
						goto l166
					}
					position++
					if buffer[position] != rune('l') {
						goto l166
					}
					position++
					if buffer[position] != rune('d') {
						goto l166
					}
					position++
					if !rules[RuleSPACE]() {
						goto l166
					}
					if !rules[RuleAryItems]() {
						goto l166
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction63]() {
						goto l166
					}
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('y') {
						goto l167
					}
					position++
					if buffer[position] != rune('i') {
						goto l167
					}
					position++
					if buffer[position] != rune('e') {
						goto l167
					}
					position++
					if buffer[position] != rune('l') {
						goto l167
					}
					position++
					if buffer[position] != rune('d') {
						goto l167
					}
					position++
					if buffer[position] != rune('(') {
						goto l167
					}
					position++
					if !rules[RuleAryItems]() {
						goto l167
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune(')') {
						goto l167
					}
					position++
					if !rules[RuleAction64]() {
						goto l167
					}
					goto l165
				l167:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('y') {
						goto l163
					}
					position++
					if buffer[position] != rune('i') {
						goto l163
					}
					position++
					if buffer[position] != rune('e') {
						goto l163
					}
					position++
					if buffer[position] != rune('l') {
						goto l163
					}
					position++
					if buffer[position] != rune('d') {
						goto l163
					}
					position++
					if !rules[RuleAction65]() {
						goto l163
					}
				}
			l165:
				depth--
				add(RuleYield, position164)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l163:
			position, tokenIndex, depth = position163, tokenIndex163, depth163
			return false
		},
		/* 29 Return <- <(('r' 'e' 't' 'u' 'r' 'n' SPACE Expr _ !',' Action66) / ('r' 'e' 't' 'u' 'r' 'n' '(' Expr ')' _ !',' Action67) / ('r' 'e' 't' 'u' 'r' 'n' SPACE AryItems Action68) / ('r' 'e' 't' 'u' 'r' 'n' '(' AryItems ')' Action69) / ('r' 'e' 't' 'u' 'r' 'n' Action70))> */
		func() bool {
			position168, tokenIndex168, depth168 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position169 := position
				depth++
				{

					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l171
					}
					position++
					if buffer[position] != rune('e') {
						goto l171
					}
					position++
					if buffer[position] != rune('t') {
						goto l171
					}
					position++
					if buffer[position] != rune('u') {
						goto l171
					}
					position++
					if buffer[position] != rune('r') {
						goto l171
					}
					position++
					if buffer[position] != rune('n') {
						goto l171
					}
					position++
					if !rules[RuleSPACE]() {
						goto l171
					}
					if !rules[RuleExpr]() {
						goto l171
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l171
					}
					{

						position172, tokenIndex172, depth172 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l172
						}
						position++
						goto l171
					l172:
						position, tokenIndex, depth = position172, tokenIndex172, depth172
					}
					if !rules[RuleAction66]() {
						goto l171
					}
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('r') {
						goto l173
					}
					position++
					if buffer[position] != rune('e') {
						goto l173
					}
					position++
					if buffer[position] != rune('t') {
						goto l173
					}
					position++
					if buffer[position] != rune('u') {
						goto l173
					}
					position++
					if buffer[position] != rune('r') {
						goto l173
					}
					position++
					if buffer[position] != rune('n') {
						goto l173
					}
					position++
					if buffer[position] != rune('(') {
						goto l173
					}
					position++
					if !rules[RuleExpr]() {
						goto l173
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune(')') {
						goto l173
					}
					position++
					if !rules[Rule_]() {
						goto l173
					}
					{

						position174, tokenIndex174, depth174 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l174
						}
						position++
						goto l173
					l174:
						position, tokenIndex, depth = position174, tokenIndex174, depth174
					}
					if !rules[RuleAction67]() {
						goto l173
					}
					goto l170
				l173:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('r') {
						goto l175
					}
					position++
					if buffer[position] != rune('e') {
						goto l175
					}
					position++
					if buffer[position] != rune('t') {
						goto l175
					}
					position++
					if buffer[position] != rune('u') {
						goto l175
					}
					position++
					if buffer[position] != rune('r') {
						goto l175
					}
					position++
					if buffer[position] != rune('n') {
						goto l175
					}
					position++
					if !rules[RuleSPACE]() {
						goto l175
					}
					if !rules[RuleAryItems]() {
						goto l175
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction68]() {
						goto l175
					}
					goto l170
				l175:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('r') {
						goto l176
					}
					position++
					if buffer[position] != rune('e') {
						goto l176
					}
					position++
					if buffer[position] != rune('t') {
						goto l176
					}
					position++
					if buffer[position] != rune('u') {
						goto l176
					}
					position++
					if buffer[position] != rune('r') {
						goto l176
					}
					position++
					if buffer[position] != rune('n') {
						goto l176
					}
					position++
					if buffer[position] != rune('(') {
						goto l176
					}
					position++
					if !rules[RuleAryItems]() {
						goto l176
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if buffer[position] != rune(')') {
						goto l176
					}
					position++
					if !rules[RuleAction69]() {
						goto l176
					}
					goto l170
				l176:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('r') {
						goto l168
					}
					position++
					if buffer[position] != rune('e') {
						goto l168
					}
					position++
					if buffer[position] != rune('t') {
						goto l168
					}
					position++
					if buffer[position] != rune('u') {
						goto l168
					}
					position++
					if buffer[position] != rune('r') {
						goto l168
					}
					position++
					if buffer[position] != rune('n') {
						goto l168
					}
					position++
					if !rules[RuleAction70]() {
						goto l168
					}
				}
			l170:
				depth--
				add(RuleReturn, position169)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l168:
			position, tokenIndex, depth = position168, tokenIndex168, depth168
			return false
		},
		/* 30 Break <- <('b' 'r' 'e' 'a' 'k' Action71)> */
		func() bool {
			position177, tokenIndex177, depth177 := position, tokenIndex, depth
			{

				position178 := position
				depth++
				if buffer[position] != rune('b') {
					goto l177
				}
				position++
				if buffer[position] != rune('r') {
					goto l177
				}
				position++
				if buffer[position] != rune('e') {
					goto l177
				}
				position++
				if buffer[position] != rune('a') {
					goto l177
				}
				position++
				if buffer[position] != rune('k') {
					goto l177
				}
				position++
				if !rules[RuleAction71]() {
					goto l177
				}
				depth--
				add(RuleBreak, position178)
			}
			return true
		l177:
			position, tokenIndex, depth = position177, tokenIndex177, depth177
			return false
		},
		/* 31 Value <- <((NUMBER Action72) / (SYMBOL Action73) / (STRING1 Action74) / (CONST Action75) / ('n' 'i' 'l' Action76) / ('t' 'r' 'u' 'e' Action77) / ('f' 'a' 'l' 's' 'e' Action78) / ('s' 'e' 'l' 'f' Action79) / (IVAR Action80) / (CVAR Action81) / (GLOBAL Action82) / ('[' _ ']' Action83) / ('[' _ AryItems _ ']' Action84) / ('{' _ '}' Action85) / ('{' _ HashItems _ '}' Action86) / ('(' _ Expr _ ')'))> */
		func() bool {
			position179, tokenIndex179, depth179 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position180 := position
				depth++
				{

					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if !rules[RuleNUMBER]() {
						goto l182
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction72]() {
						goto l182
					}
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if !rules[RuleSYMBOL]() {
						goto l183
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction73]() {
						goto l183
					}
					goto l181
				l183:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if !rules[RuleSTRING1]() {
						goto l184
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction74]() {
						goto l184
					}
					goto l181
				l184:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if !rules[RuleCONST]() {
						goto l185
					}
					variableIdx = 2
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction75]() {
						goto l185
					}
					goto l181
				l185:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('n') {
						goto l186
					}
					position++
					if buffer[position] != rune('i') {
						goto l186
					}
					position++
					if buffer[position] != rune('l') {
						goto l186
					}
					position++
					if !rules[RuleAction76]() {
						goto l186
					}
					goto l181
				l186:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('t') {
						goto l187
					}
					position++
					if buffer[position] != rune('r') {
						goto l187
					}
					position++
					if buffer[position] != rune('u') {
						goto l187
					}
					position++
					if buffer[position] != rune('e') {
						goto l187
					}
					position++
					if !rules[RuleAction77]() {
						goto l187
					}
					goto l181
				l187:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('f') {
						goto l188
					}
					position++
					if buffer[position] != rune('a') {
						goto l188
					}
					position++
					if buffer[position] != rune('l') {
						goto l188
					}
					position++
					if buffer[position] != rune('s') {
						goto l188
					}
					position++
					if buffer[position] != rune('e') {
						goto l188
					}
					position++
					if !rules[RuleAction78]() {
						goto l188
					}
					goto l181
				l188:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('s') {
						goto l189
					}
					position++
					if buffer[position] != rune('e') {
						goto l189
					}
					position++
					if buffer[position] != rune('l') {
						goto l189
					}
					position++
					if buffer[position] != rune('f') {
						goto l189
					}
					position++
					if !rules[RuleAction79]() {
						goto l189
					}
					goto l181
				l189:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if !rules[RuleIVAR]() {
						goto l190
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction80]() {
						goto l190
					}
					goto l181
				l190:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if !rules[RuleCVAR]() {
						goto l191
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction81]() {
						goto l191
					}
					goto l181
				l191:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if !rules[RuleGLOBAL]() {
						goto l192
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction82]() {
						goto l192
					}
					goto l181
				l192:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('[') {
						goto l193
					}
					position++
					if !rules[Rule_]() {
						goto l193
					}
					if buffer[position] != rune(']') {
						goto l193
					}
					position++
					if !rules[RuleAction83]() {
						goto l193
					}
					goto l181
				l193:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('[') {
						goto l194
					}
					position++
					if !rules[Rule_]() {
						goto l194
					}
					if !rules[RuleAryItems]() {
						goto l194
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l194
					}
					if buffer[position] != rune(']') {
						goto l194
					}
					position++
					if !rules[RuleAction84]() {
						goto l194
					}
					goto l181
				l194:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('{') {
						goto l195
					}
					position++
					if !rules[Rule_]() {
						goto l195
					}
					if buffer[position] != rune('}') {
						goto l195
					}
					position++
					if !rules[RuleAction85]() {
						goto l195
					}
					goto l181
				l195:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('{') {
						goto l196
					}
					position++
					if !rules[Rule_]() {
						goto l196
					}
					if !rules[RuleHashItems]() {
						goto l196
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l196
					}
					if buffer[position] != rune('}') {
						goto l196
					}
					position++
					if !rules[RuleAction86]() {
						goto l196
					}
					goto l181
				l196:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('(') {
						goto l179
					}
					position++
					if !rules[Rule_]() {
						goto l179
					}
					if !rules[RuleExpr]() {
						goto l179
					}
					if !rules[Rule_]() {
						goto l179
					}
					if buffer[position] != rune(')') {
						goto l179
					}
					position++
				}
			l181:
				depth--
				add(RuleValue, position180)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l179:
			position, tokenIndex, depth = position179, tokenIndex179, depth179
			return false
		},
		/* 32 AryItems <- <(_ Expr _ Action87 (',' _ Expr _ Action88)* Action89)> */
		func() bool {
			position197, tokenIndex197, depth197 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 2
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position198 := position
				depth++
				if !rules[Rule_]() {
					goto l197
				}
				if !rules[RuleExpr]() {
					goto l197
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l197
				}
				if !rules[RuleAction87]() {
					goto l197
				}
			l199:
				{

					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l200
					}
					position++
					if !rules[Rule_]() {
						goto l200
					}
					if !rules[RuleExpr]() {
						goto l200
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l200
					}
					if !rules[RuleAction88]() {
						goto l200
					}
					goto l199
				l200:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
				}
				if !rules[RuleAction89]() {
					goto l197
				}
				depth--
				add(RuleAryItems, position198)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l197:
			position, tokenIndex, depth = position197, tokenIndex197, depth197
			return false
		},
		/* 33 HashItems <- <(Expr _ ('=' '>') _ Expr Action90 (_ ',' _ Expr _ Action91 ('=' '>') _ Expr Action92)* Action93)> */
		func() bool {
			position201, tokenIndex201, depth201 := position, tokenIndex, depth
			variableIdx := 0
			variableTotal := 3
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPush, position)
			}
			{

				position202 := position
				depth++
				if !rules[RuleExpr]() {
					goto l201
				}
				variableIdx = 2
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[Rule_]() {
					goto l201
				}
				if buffer[position] != rune('=') {
					goto l201
				}
				position++
				if buffer[position] != rune('>') {
					goto l201
				}
				position++
				if !rules[Rule_]() {
					goto l201
				}
				if !rules[RuleExpr]() {
					goto l201
				}
				variableIdx = 1
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPop, position)
				}
				add(RuleActionSet, position)
				for i := 0; i < variableIdx; i++ {
					add(RuleActionPush, position)
				}
				if !rules[RuleAction90]() {
					goto l201
				}
			l203:
				{

					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if !rules[Rule_]() {
						goto l204
					}
					if buffer[position] != rune(',') {
						goto l204
					}
					position++
					if !rules[Rule_]() {
						goto l204
					}
					if !rules[RuleExpr]() {
						goto l204
					}
					variableIdx = 0
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[Rule_]() {
						goto l204
					}
					if !rules[RuleAction91]() {
						goto l204
					}
					if buffer[position] != rune('=') {
						goto l204
					}
					position++
					if buffer[position] != rune('>') {
						goto l204
					}
					position++
					if !rules[Rule_]() {
						goto l204
					}
					if !rules[RuleExpr]() {
						goto l204
					}
					variableIdx = 1
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPop, position)
					}
					add(RuleActionSet, position)
					for i := 0; i < variableIdx; i++ {
						add(RuleActionPush, position)
					}
					if !rules[RuleAction92]() {
						goto l204
					}
					goto l203
				l204:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
				}
				if !rules[RuleAction93]() {
					goto l201
				}
				depth--
				add(RuleHashItems, position202)
			}
			for i := 0; i < variableTotal; i++ {
				add(RuleActionPop, position)
			}
			return true
		l201:
			position, tokenIndex, depth = position201, tokenIndex201, depth201
			return false
		},
		/* 34 KEYWORD <- <(('w' 'h' 'i' 'l' 'e') / ('u' 'n' 't' 'i' 'l') / ('d' 'o') / ('e' 'n' 'd') / ('i' 'f') / ('u' 'n' 'l' 'e' 's' 's') / ('e' 'l' 's' 'e') / ('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e') / ('n' 'i' 'l') / ('s' 'e' 'l' 'f') / ('c' 'l' 'a' 's' 's') / ('m' 'o' 'd' 'u' 'l' 'e') / ('d' 'e' 'f') / ('y' 'i' 'e' 'l' 'd') / ('r' 'e' 't' 'u' 'r' 'n') / ('b' 'r' 'e' 'a' 'k'))> */
		func() bool {
			position205, tokenIndex205, depth205 := position, tokenIndex, depth
			{

				position206 := position
				depth++
				{

					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l208
					}
					position++
					if buffer[position] != rune('h') {
						goto l208
					}
					position++
					if buffer[position] != rune('i') {
						goto l208
					}
					position++
					if buffer[position] != rune('l') {
						goto l208
					}
					position++
					if buffer[position] != rune('e') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('u') {
						goto l209
					}
					position++
					if buffer[position] != rune('n') {
						goto l209
					}
					position++
					if buffer[position] != rune('t') {
						goto l209
					}
					position++
					if buffer[position] != rune('i') {
						goto l209
					}
					position++
					if buffer[position] != rune('l') {
						goto l209
					}
					position++
					goto l207
				l209:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('d') {
						goto l210
					}
					position++
					if buffer[position] != rune('o') {
						goto l210
					}
					position++
					goto l207
				l210:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('e') {
						goto l211
					}
					position++
					if buffer[position] != rune('n') {
						goto l211
					}
					position++
					if buffer[position] != rune('d') {
						goto l211
					}
					position++
					goto l207
				l211:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('i') {
						goto l212
					}
					position++
					if buffer[position] != rune('f') {
						goto l212
					}
					position++
					goto l207
				l212:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('u') {
						goto l213
					}
					position++
					if buffer[position] != rune('n') {
						goto l213
					}
					position++
					if buffer[position] != rune('l') {
						goto l213
					}
					position++
					if buffer[position] != rune('e') {
						goto l213
					}
					position++
					if buffer[position] != rune('s') {
						goto l213
					}
					position++
					if buffer[position] != rune('s') {
						goto l213
					}
					position++
					goto l207
				l213:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('e') {
						goto l214
					}
					position++
					if buffer[position] != rune('l') {
						goto l214
					}
					position++
					if buffer[position] != rune('s') {
						goto l214
					}
					position++
					if buffer[position] != rune('e') {
						goto l214
					}
					position++
					goto l207
				l214:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('t') {
						goto l215
					}
					position++
					if buffer[position] != rune('r') {
						goto l215
					}
					position++
					if buffer[position] != rune('u') {
						goto l215
					}
					position++
					if buffer[position] != rune('e') {
						goto l215
					}
					position++
					goto l207
				l215:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('f') {
						goto l216
					}
					position++
					if buffer[position] != rune('a') {
						goto l216
					}
					position++
					if buffer[position] != rune('l') {
						goto l216
					}
					position++
					if buffer[position] != rune('s') {
						goto l216
					}
					position++
					if buffer[position] != rune('e') {
						goto l216
					}
					position++
					goto l207
				l216:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('n') {
						goto l217
					}
					position++
					if buffer[position] != rune('i') {
						goto l217
					}
					position++
					if buffer[position] != rune('l') {
						goto l217
					}
					position++
					goto l207
				l217:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('s') {
						goto l218
					}
					position++
					if buffer[position] != rune('e') {
						goto l218
					}
					position++
					if buffer[position] != rune('l') {
						goto l218
					}
					position++
					if buffer[position] != rune('f') {
						goto l218
					}
					position++
					goto l207
				l218:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('c') {
						goto l219
					}
					position++
					if buffer[position] != rune('l') {
						goto l219
					}
					position++
					if buffer[position] != rune('a') {
						goto l219
					}
					position++
					if buffer[position] != rune('s') {
						goto l219
					}
					position++
					if buffer[position] != rune('s') {
						goto l219
					}
					position++
					goto l207
				l219:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('m') {
						goto l220
					}
					position++
					if buffer[position] != rune('o') {
						goto l220
					}
					position++
					if buffer[position] != rune('d') {
						goto l220
					}
					position++
					if buffer[position] != rune('u') {
						goto l220
					}
					position++
					if buffer[position] != rune('l') {
						goto l220
					}
					position++
					if buffer[position] != rune('e') {
						goto l220
					}
					position++
					goto l207
				l220:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('d') {
						goto l221
					}
					position++
					if buffer[position] != rune('e') {
						goto l221
					}
					position++
					if buffer[position] != rune('f') {
						goto l221
					}
					position++
					goto l207
				l221:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('y') {
						goto l222
					}
					position++
					if buffer[position] != rune('i') {
						goto l222
					}
					position++
					if buffer[position] != rune('e') {
						goto l222
					}
					position++
					if buffer[position] != rune('l') {
						goto l222
					}
					position++
					if buffer[position] != rune('d') {
						goto l222
					}
					position++
					goto l207
				l222:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('r') {
						goto l223
					}
					position++
					if buffer[position] != rune('e') {
						goto l223
					}
					position++
					if buffer[position] != rune('t') {
						goto l223
					}
					position++
					if buffer[position] != rune('u') {
						goto l223
					}
					position++
					if buffer[position] != rune('r') {
						goto l223
					}
					position++
					if buffer[position] != rune('n') {
						goto l223
					}
					position++
					goto l207
				l223:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('b') {
						goto l205
					}
					position++
					if buffer[position] != rune('r') {
						goto l205
					}
					position++
					if buffer[position] != rune('e') {
						goto l205
					}
					position++
					if buffer[position] != rune('a') {
						goto l205
					}
					position++
					if buffer[position] != rune('k') {
						goto l205
					}
					position++
				}
			l207:
				depth--
				add(RuleKEYWORD, position206)
			}
			return true
		l205:
			position, tokenIndex, depth = position205, tokenIndex205, depth205
			return false
		},
		/* 35 NAME <- <([a-z] / [A-Z] / [0-9] / '_')+> */
		func() bool {
			position224, tokenIndex224, depth224 := position, tokenIndex, depth
			{

				position225 := position
				depth++
				{

					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l230
					}
					position++
					goto l228
				l230:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l231
					}
					position++
					goto l228
				l231:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('_') {
						goto l224
					}
					position++
				}
			l228:
			l226:
				{

					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					{

						position232, tokenIndex232, depth232 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l233
						}
						position++
						goto l232
					l233:
						position, tokenIndex, depth = position232, tokenIndex232, depth232
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l234
						}
						position++
						goto l232
					l234:
						position, tokenIndex, depth = position232, tokenIndex232, depth232
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l235
						}
						position++
						goto l232
					l235:
						position, tokenIndex, depth = position232, tokenIndex232, depth232
						if buffer[position] != rune('_') {
							goto l227
						}
						position++
					}
				l232:
					goto l226
				l227:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
				}
				depth--
				add(RuleNAME, position225)
			}
			return true
		l224:
			position, tokenIndex, depth = position224, tokenIndex224, depth224
			return false
		},
		/* 36 ID <- <((!('s' 'e' 'l' 'f') <KEYWORD> &('.' / '(' / '[') Action94) / (<(KEYWORD NAME)> Action95) / (!KEYWORD <(([a-z] / '_') NAME? (('=' &'(') / '!' / '?')?)> Action96))> */
		func() bool {
			position236, tokenIndex236, depth236 := position, tokenIndex, depth
			{

				position237 := position
				depth++
				{

					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					{

						position240, tokenIndex240, depth240 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l240
						}
						position++
						if buffer[position] != rune('e') {
							goto l240
						}
						position++
						if buffer[position] != rune('l') {
							goto l240
						}
						position++
						if buffer[position] != rune('f') {
							goto l240
						}
						position++
						goto l239
					l240:
						position, tokenIndex, depth = position240, tokenIndex240, depth240
					}
					{

						position241 := position
						depth++
						if !rules[RuleKEYWORD]() {
							goto l239
						}
						depth--
						add(RulePegText, position241)
					}
					{

						position242, tokenIndex242, depth242 := position, tokenIndex, depth
						{

							position243, tokenIndex243, depth243 := position, tokenIndex, depth
							if buffer[position] != rune('.') {
								goto l244
							}
							position++
							goto l243
						l244:
							position, tokenIndex, depth = position243, tokenIndex243, depth243
							if buffer[position] != rune('(') {
								goto l245
							}
							position++
							goto l243
						l245:
							position, tokenIndex, depth = position243, tokenIndex243, depth243
							if buffer[position] != rune('[') {
								goto l239
							}
							position++
						}
					l243:
						position, tokenIndex, depth = position242, tokenIndex242, depth242
					}
					if !rules[RuleAction94]() {
						goto l239
					}
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					{

						position247 := position
						depth++
						if !rules[RuleKEYWORD]() {
							goto l246
						}
						if !rules[RuleNAME]() {
							goto l246
						}
						depth--
						add(RulePegText, position247)
					}
					if !rules[RuleAction95]() {
						goto l246
					}
					goto l238
				l246:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					{

						position248, tokenIndex248, depth248 := position, tokenIndex, depth
						if !rules[RuleKEYWORD]() {
							goto l248
						}
						goto l236
					l248:
						position, tokenIndex, depth = position248, tokenIndex248, depth248
					}
					{

						position249 := position
						depth++
						{

							position250, tokenIndex250, depth250 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l251
							}
							position++
							goto l250
						l251:
							position, tokenIndex, depth = position250, tokenIndex250, depth250
							if buffer[position] != rune('_') {
								goto l236
							}
							position++
						}
					l250:
						{

							position252, tokenIndex252, depth252 := position, tokenIndex, depth
							if !rules[RuleNAME]() {
								goto l252
							}
							goto l253
						l252:
							position, tokenIndex, depth = position252, tokenIndex252, depth252
						}
					l253:
						{

							position254, tokenIndex254, depth254 := position, tokenIndex, depth
							{

								position256, tokenIndex256, depth256 := position, tokenIndex, depth
								if buffer[position] != rune('=') {
									goto l257
								}
								position++
								{

									position258, tokenIndex258, depth258 := position, tokenIndex, depth
									if buffer[position] != rune('(') {
										goto l257
									}
									position++
									position, tokenIndex, depth = position258, tokenIndex258, depth258
								}
								goto l256
							l257:
								position, tokenIndex, depth = position256, tokenIndex256, depth256
								if buffer[position] != rune('!') {
									goto l259
								}
								position++
								goto l256
							l259:
								position, tokenIndex, depth = position256, tokenIndex256, depth256
								if buffer[position] != rune('?') {
									goto l254
								}
								position++
							}
						l256:
							goto l255
						l254:
							position, tokenIndex, depth = position254, tokenIndex254, depth254
						}
					l255:
						depth--
						add(RulePegText, position249)
					}
					if !rules[RuleAction96]() {
						goto l236
					}
				}
			l238:
				depth--
				add(RuleID, position237)
			}
			return true
		l236:
			position, tokenIndex, depth = position236, tokenIndex236, depth236
			return false
		},
		/* 37 CONST <- <(<([A-Z] NAME?)> Action97)> */
		func() bool {
			position260, tokenIndex260, depth260 := position, tokenIndex, depth
			{

				position261 := position
				depth++
				{

					position262 := position
					depth++
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l260
					}
					position++
					{

						position263, tokenIndex263, depth263 := position, tokenIndex, depth
						if !rules[RuleNAME]() {
							goto l263
						}
						goto l264
					l263:
						position, tokenIndex, depth = position263, tokenIndex263, depth263
					}
				l264:
					depth--
					add(RulePegText, position262)
				}
				if !rules[RuleAction97]() {
					goto l260
				}
				depth--
				add(RuleCONST, position261)
			}
			return true
		l260:
			position, tokenIndex, depth = position260, tokenIndex260, depth260
			return false
		},
		/* 38 BINOP <- <(<(('*' '*') / '^' / '&' / '|' / '~' / '+' / '-' / '*' / '/' / '%' / ('<' '=' '>') / ('<' '<') / ('>' '>') / ('=' '=') / ('=' '~') / ('!' '=') / ('=' '=' '=') / '<' / '>' / ('<' '=') / ('>' '='))> Action98)> */
		func() bool {
			position265, tokenIndex265, depth265 := position, tokenIndex, depth
			{

				position266 := position
				depth++
				{

					position267 := position
					depth++
					{

						position268, tokenIndex268, depth268 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l269
						}
						position++
						if buffer[position] != rune('*') {
							goto l269
						}
						position++
						goto l268
					l269:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('^') {
							goto l270
						}
						position++
						goto l268
					l270:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('&') {
							goto l271
						}
						position++
						goto l268
					l271:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('|') {
							goto l272
						}
						position++
						goto l268
					l272:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('~') {
							goto l273
						}
						position++
						goto l268
					l273:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('+') {
							goto l274
						}
						position++
						goto l268
					l274:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('-') {
							goto l275
						}
						position++
						goto l268
					l275:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('*') {
							goto l276
						}
						position++
						goto l268
					l276:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('/') {
							goto l277
						}
						position++
						goto l268
					l277:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('%') {
							goto l278
						}
						position++
						goto l268
					l278:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('<') {
							goto l279
						}
						position++
						if buffer[position] != rune('=') {
							goto l279
						}
						position++
						if buffer[position] != rune('>') {
							goto l279
						}
						position++
						goto l268
					l279:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('<') {
							goto l280
						}
						position++
						if buffer[position] != rune('<') {
							goto l280
						}
						position++
						goto l268
					l280:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('>') {
							goto l281
						}
						position++
						if buffer[position] != rune('>') {
							goto l281
						}
						position++
						goto l268
					l281:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('=') {
							goto l282
						}
						position++
						if buffer[position] != rune('=') {
							goto l282
						}
						position++
						goto l268
					l282:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('=') {
							goto l283
						}
						position++
						if buffer[position] != rune('~') {
							goto l283
						}
						position++
						goto l268
					l283:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('!') {
							goto l284
						}
						position++
						if buffer[position] != rune('=') {
							goto l284
						}
						position++
						goto l268
					l284:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('=') {
							goto l285
						}
						position++
						if buffer[position] != rune('=') {
							goto l285
						}
						position++
						if buffer[position] != rune('=') {
							goto l285
						}
						position++
						goto l268
					l285:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('<') {
							goto l286
						}
						position++
						goto l268
					l286:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('>') {
							goto l287
						}
						position++
						goto l268
					l287:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('<') {
							goto l288
						}
						position++
						if buffer[position] != rune('=') {
							goto l288
						}
						position++
						goto l268
					l288:
						position, tokenIndex, depth = position268, tokenIndex268, depth268
						if buffer[position] != rune('>') {
							goto l265
						}
						position++
						if buffer[position] != rune('=') {
							goto l265
						}
						position++
					}
				l268:
					depth--
					add(RulePegText, position267)
				}
				if !rules[RuleAction98]() {
					goto l265
				}
				depth--
				add(RuleBINOP, position266)
			}
			return true
		l265:
			position, tokenIndex, depth = position265, tokenIndex265, depth265
			return false
		},
		/* 39 UNOP <- <(<(('-' '@') / '!')> Action99)> */
		func() bool {
			position289, tokenIndex289, depth289 := position, tokenIndex, depth
			{

				position290 := position
				depth++
				{

					position291 := position
					depth++
					{

						position292, tokenIndex292, depth292 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l293
						}
						position++
						if buffer[position] != rune('@') {
							goto l293
						}
						position++
						goto l292
					l293:
						position, tokenIndex, depth = position292, tokenIndex292, depth292
						if buffer[position] != rune('!') {
							goto l289
						}
						position++
					}
				l292:
					depth--
					add(RulePegText, position291)
				}
				if !rules[RuleAction99]() {
					goto l289
				}
				depth--
				add(RuleUNOP, position290)
			}
			return true
		l289:
			position, tokenIndex, depth = position289, tokenIndex289, depth289
			return false
		},
		/* 40 METHOD <- <(ID / UNOP / BINOP)> */
		func() bool {
			position294, tokenIndex294, depth294 := position, tokenIndex, depth
			{

				position295 := position
				depth++
				{

					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if !rules[RuleID]() {
						goto l297
					}
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if !rules[RuleUNOP]() {
						goto l298
					}
					goto l296
				l298:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if !rules[RuleBINOP]() {
						goto l294
					}
				}
			l296:
				depth--
				add(RuleMETHOD, position295)
			}
			return true
		l294:
			position, tokenIndex, depth = position294, tokenIndex294, depth294
			return false
		},
		/* 41 ASSIGN <- <(<'='> &!'=' Action100)> */
		func() bool {
			position299, tokenIndex299, depth299 := position, tokenIndex, depth
			{

				position300 := position
				depth++
				{

					position301 := position
					depth++
					if buffer[position] != rune('=') {
						goto l299
					}
					position++
					depth--
					add(RulePegText, position301)
				}
				{

					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					{

						position303, tokenIndex303, depth303 := position, tokenIndex, depth
						if buffer[position] != rune('=') {
							goto l303
						}
						position++
						goto l299
					l303:
						position, tokenIndex, depth = position303, tokenIndex303, depth303
					}
					position, tokenIndex, depth = position302, tokenIndex302, depth302
				}
				if !rules[RuleAction100]() {
					goto l299
				}
				depth--
				add(RuleASSIGN, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 42 IVAR <- <(<('@' NAME)> Action101)> */
		func() bool {
			position304, tokenIndex304, depth304 := position, tokenIndex, depth
			{

				position305 := position
				depth++
				{

					position306 := position
					depth++
					if buffer[position] != rune('@') {
						goto l304
					}
					position++
					if !rules[RuleNAME]() {
						goto l304
					}
					depth--
					add(RulePegText, position306)
				}
				if !rules[RuleAction101]() {
					goto l304
				}
				depth--
				add(RuleIVAR, position305)
			}
			return true
		l304:
			position, tokenIndex, depth = position304, tokenIndex304, depth304
			return false
		},
		/* 43 CVAR <- <(<('@' '@' NAME)> Action102)> */
		func() bool {
			position307, tokenIndex307, depth307 := position, tokenIndex, depth
			{

				position308 := position
				depth++
				{

					position309 := position
					depth++
					if buffer[position] != rune('@') {
						goto l307
					}
					position++
					if buffer[position] != rune('@') {
						goto l307
					}
					position++
					if !rules[RuleNAME]() {
						goto l307
					}
					depth--
					add(RulePegText, position309)
				}
				if !rules[RuleAction102]() {
					goto l307
				}
				depth--
				add(RuleCVAR, position308)
			}
			return true
		l307:
			position, tokenIndex, depth = position307, tokenIndex307, depth307
			return false
		},
		/* 44 GLOBAL <- <(<('$' NAME)> Action103)> */
		func() bool {
			position310, tokenIndex310, depth310 := position, tokenIndex, depth
			{

				position311 := position
				depth++
				{

					position312 := position
					depth++
					if buffer[position] != rune('$') {
						goto l310
					}
					position++
					if !rules[RuleNAME]() {
						goto l310
					}
					depth--
					add(RulePegText, position312)
				}
				if !rules[RuleAction103]() {
					goto l310
				}
				depth--
				add(RuleGLOBAL, position311)
			}
			return true
		l310:
			position, tokenIndex, depth = position310, tokenIndex310, depth310
			return false
		},
		/* 45 NUMBER <- <(<[0-9]+> Action104)> */
		func() bool {
			position313, tokenIndex313, depth313 := position, tokenIndex, depth
			{

				position314 := position
				depth++
				{

					position315 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l313
					}
					position++
				l316:
					{

						position317, tokenIndex317, depth317 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l317
						}
						position++
						goto l316
					l317:
						position, tokenIndex, depth = position317, tokenIndex317, depth317
					}
					depth--
					add(RulePegText, position315)
				}
				if !rules[RuleAction104]() {
					goto l313
				}
				depth--
				add(RuleNUMBER, position314)
			}
			return true
		l313:
			position, tokenIndex, depth = position313, tokenIndex313, depth313
			return false
		},
		/* 46 SYMBOL <- <(':' <(NAME / KEYWORD)> Action105)> */
		func() bool {
			position318, tokenIndex318, depth318 := position, tokenIndex, depth
			{

				position319 := position
				depth++
				if buffer[position] != rune(':') {
					goto l318
				}
				position++
				{

					position320 := position
					depth++
					{

						position321, tokenIndex321, depth321 := position, tokenIndex, depth
						if !rules[RuleNAME]() {
							goto l322
						}
						goto l321
					l322:
						position, tokenIndex, depth = position321, tokenIndex321, depth321
						if !rules[RuleKEYWORD]() {
							goto l318
						}
					}
				l321:
					depth--
					add(RulePegText, position320)
				}
				if !rules[RuleAction105]() {
					goto l318
				}
				depth--
				add(RuleSYMBOL, position319)
			}
			return true
		l318:
			position, tokenIndex, depth = position318, tokenIndex318, depth318
			return false
		},
		/* 47 STRING1 <- <('\'' Action106 (('\\' '\'' Action107) / (<(!'\'' .)> Action108))* '\'' Action109)> */
		func() bool {
			position323, tokenIndex323, depth323 := position, tokenIndex, depth
			{

				position324 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l323
				}
				position++
				if !rules[RuleAction106]() {
					goto l323
				}
			l325:
				{

					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					{

						position327, tokenIndex327, depth327 := position, tokenIndex, depth
						if buffer[position] != rune('\\') {
							goto l328
						}
						position++
						if buffer[position] != rune('\'') {
							goto l328
						}
						position++
						if !rules[RuleAction107]() {
							goto l328
						}
						goto l327
					l328:
						position, tokenIndex, depth = position327, tokenIndex327, depth327
						{

							position329 := position
							depth++
							{

								position330, tokenIndex330, depth330 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l330
								}
								position++
								goto l326
							l330:
								position, tokenIndex, depth = position330, tokenIndex330, depth330
							}
							if !matchDot() {
								goto l326
							}
							depth--
							add(RulePegText, position329)
						}
						if !rules[RuleAction108]() {
							goto l326
						}
					}
				l327:
					goto l325
				l326:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
				}
				if buffer[position] != rune('\'') {
					goto l323
				}
				position++
				if !rules[RuleAction109]() {
					goto l323
				}
				depth--
				add(RuleSTRING1, position324)
			}
			return true
		l323:
			position, tokenIndex, depth = position323, tokenIndex323, depth323
			return false
		},
		/* 48 _ <- <(' ' / '\t')*> */
		func() bool {
			{

				position332 := position
				depth++
			l333:
				{

					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					{

						position335, tokenIndex335, depth335 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l336
						}
						position++
						goto l335
					l336:
						position, tokenIndex, depth = position335, tokenIndex335, depth335
						if buffer[position] != rune('\t') {
							goto l334
						}
						position++
					}
				l335:
					goto l333
				l334:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
				}
				depth--
				add(Rule_, position332)
			}
			return true
		},
		/* 49 SPACE <- <' '+> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{

				position338 := position
				depth++
				if buffer[position] != rune(' ') {
					goto l337
				}
				position++
			l339:
				{

					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l340
					}
					position++
					goto l339
				l340:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
				}
				depth--
				add(RuleSPACE, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 50 EOL <- <(('\n' / ('\r' '\n') / '\r') Action110)> */
		func() bool {
			position341, tokenIndex341, depth341 := position, tokenIndex, depth
			{

				position342 := position
				depth++
				{

					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if buffer[position] != rune('\n') {
						goto l344
					}
					position++
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('\r') {
						goto l345
					}
					position++
					if buffer[position] != rune('\n') {
						goto l345
					}
					position++
					goto l343
				l345:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('\r') {
						goto l341
					}
					position++
				}
			l343:
				if !rules[RuleAction110]() {
					goto l341
				}
				depth--
				add(RuleEOL, position342)
			}
			return true
		l341:
			position, tokenIndex, depth = position341, tokenIndex341, depth341
			return false
		},
		/* 51 EOF <- <!.> */
		func() bool {
			position346, tokenIndex346, depth346 := position, tokenIndex, depth
			{

				position347 := position
				depth++
				{

					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if !matchDot() {
						goto l348
					}
					goto l346
				l348:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
				}
				depth--
				add(RuleEOF, position347)
			}
			return true
		l346:
			position, tokenIndex, depth = position346, tokenIndex346, depth346
			return false
		},
		/* 52 SEP <- <(_ Comment? (EOL / ';'))+> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{

				position350 := position
				depth++
				if !rules[Rule_]() {
					goto l349
				}
				{

					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if !rules[RuleComment]() {
						goto l353
					}
					goto l354
				l353:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
				}
			l354:
				{

					position355, tokenIndex355, depth355 := position, tokenIndex, depth
					if !rules[RuleEOL]() {
						goto l356
					}
					goto l355
				l356:
					position, tokenIndex, depth = position355, tokenIndex355, depth355
					if buffer[position] != rune(';') {
						goto l349
					}
					position++
				}
			l355:
			l351:
				{

					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if !rules[Rule_]() {
						goto l352
					}
					{

						position357, tokenIndex357, depth357 := position, tokenIndex, depth
						if !rules[RuleComment]() {
							goto l357
						}
						goto l358
					l357:
						position, tokenIndex, depth = position357, tokenIndex357, depth357
					}
				l358:
					{

						position359, tokenIndex359, depth359 := position, tokenIndex, depth
						if !rules[RuleEOL]() {
							goto l360
						}
						goto l359
					l360:
						position, tokenIndex, depth = position359, tokenIndex359, depth359
						if buffer[position] != rune(';') {
							goto l352
						}
						position++
					}
				l359:
					goto l351
				l352:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
				}
				depth--
				add(RuleSEP, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 54 Action0 <- <{ rootAST = MakeASTNode(NODE_ROOT, stack[stack_idx-0], nil, nil, currentLine) }> */
		func() bool {
			{

				add(RuleAction0, position)
			}
			return true
		},
		/* 55 Action1 <- <{ //Stmts:0}> */
		func() bool {
			{

				add(RuleAction1, position)
			}
			return true
		},
		/* 56 Action2 <- <{ stack[stack_idx-1].PushBack(stack[stack_idx-0]) //Stmts:1 }> */
		func() bool {
			{

				add(RuleAction2, position)
			}
			return true
		},
		/* 57 Action3 <- <{ yy = stack[stack_idx-1] //Stmts:2}> */
		func() bool {
			{

				add(RuleAction3, position)
			}
			return true
		},
		/* 58 Action4 <- <{ yy = nil //Stmts:3}> */
		func() bool {
			{

				add(RuleAction4, position)
			}
			return true
		},
		/* 59 Action5 <- <{ yy = nil //OptStmts:0}> */
		func() bool {
			{

				add(RuleAction5, position)
			}
			return true
		},
		/* 60 Action6 <- <{ stack[stack_idx-3] = nil; stack[stack_idx-0] = nil;//Call:0 }> */
		func() bool {
			{

				add(RuleAction6, position)
			}
			return true
		},
		/* 61 Action7 <- <{ currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-3], stack[stack_idx-2], nil, currentLine); stack[stack_idx-3] = currentAST; //Call:1 }> */
		func() bool {
			{

				add(RuleAction7, position)
			}
			return true
		},
		/* 62 Action8 <- <{ currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-3], stack[stack_idx-1], stack[stack_idx-0], currentLine); yy = currentAST; //Call:2 }> */
		func() bool {
			{

				add(RuleAction8, position)
			}
			return true
		},
		/* 63 Action9 <- <{ stack[stack_idx-4] = nil //AsgnCall:0 }> */
		func() bool {
			{

				add(RuleAction9, position)
			}
			return true
		},
		/* 64 Action10 <- <{ currentAST := MakeASTNode(NODE_SEND,stack[stack_idx-4], stack[stack_idx-3], nil, currentLine); stack[stack_idx-4] = currentAST; //AsgnCall:1 }> */
		func() bool {
			{

				add(RuleAction10, position)
			}
			return true
		},
		/* 65 Action11 <- <{ argAST := MakeASTNode(NODE_ARG, stack[stack_idx-0], nil, nil, currentLine); msgVal := []string{stack[stack_idx-2].value.str, "="}; stack[stack_idx-2].value.str = strings.Join(msgVal, ""); msgAST := MakeASTNode(NODE_MSG, stack[stack_idx-2], argAST, nil, currentLine); yy = MakeASTNode(NODE_SEND, stack[stack_idx-4], msgAST, nil, currentLine); //AsgnCall:2}> */
		func() bool {
			{

				add(RuleAction11, position)
			}
			return true
		},
		/* 66 Action12 <- <{ stack[stack_idx-0] = nil //Receiver:0}> */
		func() bool {
			{

				add(RuleAction12, position)
			}
			return true
		},
		/* 67 Action13 <- <{ yy = stack[stack_idx-0] //Receiver:1}> */
		func() bool {
			{

				add(RuleAction13, position)
			}
			return true
		},
		/* 68 Action14 <- <{ yy = stack[stack_idx-0] //Receiver:2}> */
		func() bool {
			{

				add(RuleAction14, position)
			}
			return true
		},
		/* 69 Action15 <- <{ currentAST :=  MakeASTNode(NODE_ARG, stack[stack_idx-0], nil, nil, currentLine); stack[stack_idx-1].PushBack(currentAST); methodAST :=  &AST{Type: NODE_ASTVAL}; methodAST.value.str = "[]="; msgAST := MakeASTNode(NODE_MSG, methodAST, stack[stack_idx-1], nil, currentLine); currentAST = MakeASTNode(NODE_SEND, stack[stack_idx-2], msgAST, nil, currentLine); yy = currentAST; //SpecCall:0 }> */
		func() bool {
			{

				add(RuleAction15, position)
			}
			return true
		},
		/* 70 Action16 <- <{ methodAST :=  &AST{Type: NODE_ASTVAL}; methodAST.value.str = "[]"; msgAST := MakeASTNode(NODE_MSG, methodAST, stack[stack_idx-1], nil, currentLine); currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-2], msgAST, nil, currentLine); yy = currentAST; //SpecCall:1 }> */
		func() bool {
			{

				add(RuleAction16, position)
			}
			return true
		},
		/* 71 Action17 <- <{ currentAST := MakeASTNode(NODE_AND, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine); yy = currentAST; //BinOp:0 && }> */
		func() bool {
			{

				add(RuleAction17, position)
			}
			return true
		},
		/* 72 Action18 <- <{ currentAST := MakeASTNode(NODE_OR, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine); yy = currentAST;  //BinOp:1 || }> */
		func() bool {
			{

				add(RuleAction18, position)
			}
			return true
		},
		/* 73 Action19 <- <{ currentAST := MakeASTNode(NODE_ADD, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine); yy = currentAST;  //BinOp:2 + }> */
		func() bool {
			{

				add(RuleAction19, position)
			}
			return true
		},
		/* 74 Action20 <- <{ currentAST := MakeASTNode(NODE_SUB, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine); yy = currentAST;  //BinOp:3 - }> */
		func() bool {
			{

				add(RuleAction20, position)
			}
			return true
		},
		/* 75 Action21 <- <{ currentAST := MakeASTNode(NODE_LT, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine); yy = currentAST;  //BinOp:4 < }> */
		func() bool {
			{

				add(RuleAction21, position)
			}
			return true
		},
		/* 76 Action22 <- <{ argAST := MakeASTNode(NODE_ARG, stack[stack_idx-1], nil, nil, currentLine); msgAST := MakeASTNode(NODE_MSG, stack[stack_idx-0], argAST, nil, currentLine); currentAST := MakeASTNode(NODE_SEND, stack[stack_idx-2], msgAST, nil, currentLine); yy = currentAST; //BinOp:5 BINOP }> */
		func() bool {
			{

				add(RuleAction22, position)
			}
			return true
		},
		/* 77 Action23 <- <{ currentAST := MakeASTNode(NODE_NEG, stack[stack_idx-0], nil, nil, currentLine); yy = currentAST; //UnaryOp:0 NODE_NEG - }> */
		func() bool {
			{

				add(RuleAction23, position)
			}
			return true
		},
		/* 78 Action24 <- <{ currentAST := MakeASTNode(NODE_NOT, stack[stack_idx-0], nil, nil, currentLine); yy = currentAST; //UnaryOp:0 NODE_NOT ! }> */
		func() bool {
			{

				add(RuleAction24, position)
			}
			return true
		},
		/* 79 Action25 <- <{ stack[stack_idx-0] = nil //Message:0 }> */
		func() bool {
			{

				add(RuleAction25, position)
			}
			return true
		},
		/* 80 Action26 <- <{ currentAST := MakeASTNode(NODE_MSG, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST; //Message:1 }> */
		func() bool {
			{

				add(RuleAction26, position)
			}
			return true
		},
		/* 81 Action27 <- <{ currentAST := MakeASTNode(NODE_ARG, stack[stack_idx-2], nil, nil, currentLine); stack[stack_idx-2] = currentAST; //Args:0 }> */
		func() bool {
			{

				add(RuleAction27, position)
			}
			return true
		},
		/* 82 Action28 <- <{ stack[stack_idx-2].args[0].PushBack(stack[stack_idx-1]) //Args:1 }> */
		func() bool {
			{

				add(RuleAction28, position)
			}
			return true
		},
		/* 83 Action29 <- <{ //Args:2 No going to support :p}> */
		func() bool {
			{

				add(RuleAction29, position)
			}
			return true
		},
		/* 84 Action30 <- <{ yy = stack[stack_idx-2] //Args:3 }> */
		func() bool {
			{

				add(RuleAction30, position)
			}
			return true
		},
		/* 85 Action31 <- <{ currentAST := MakeASTNode(NODE_BLOCK, stack[stack_idx-1], nil, nil, currentLine); yy = currentAST;  //Block:0 }> */
		func() bool {
			{

				add(RuleAction31, position)
			}
			return true
		},
		/* 86 Action32 <- <{ currentAST := MakeASTNode(NODE_BLOCK, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST;  //Block:1 with stack[stack_idx-0] }> */
		func() bool {
			{

				add(RuleAction32, position)
			}
			return true
		},
		/* 87 Action33 <- <{ currentAST := MakeASTNode(NODE_ASSIGN, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST; //Assign:0 }> */
		func() bool {
			{

				add(RuleAction33, position)
			}
			return true
		},
		/* 88 Action34 <- <{ currentAST := MakeASTNode(NODE_SETCONST, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST; //Assign:1 }> */
		func() bool {
			{

				add(RuleAction34, position)
			}
			return true
		},
		/* 89 Action35 <- <{ currentAST := MakeASTNode(NODE_SETIVAR, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST; //Assign:2 }> */
		func() bool {
			{

				add(RuleAction35, position)
			}
			return true
		},
		/* 90 Action36 <- <{ currentAST := MakeASTNode(NODE_SETCVAR, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST; //Assign:3 }> */
		func() bool {
			{

				add(RuleAction36, position)
			}
			return true
		},
		/* 91 Action37 <- <{ currentAST := MakeASTNode(NODE_SETGLOBAL, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST; //Assign:4 }> */
		func() bool {
			{

				add(RuleAction37, position)
			}
			return true
		},
		/* 92 Action38 <- <{ currentAST := MakeASTNode(NODE_WHILE, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST;  //While NODE_WHILE }> */
		func() bool {
			{

				add(RuleAction38, position)
			}
			return true
		},
		/* 93 Action39 <- <{ currentAST := MakeASTNode(NODE_UNTIL, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST;  //Until NODE_UNTIL }> */
		func() bool {
			{

				add(RuleAction39, position)
			}
			return true
		},
		/* 94 Action40 <- <{ stack[stack_idx-0] = nil //If:0 stack[stack_idx-0] = 0}> */
		func() bool {
			{

				add(RuleAction40, position)
			}
			return true
		},
		/* 95 Action41 <- <{ currentAST := MakeASTNode(NODE_IF, stack[stack_idx-2], stack[stack_idx-1], stack[stack_idx-0], currentLine); yy = currentAST;  //If:1 }> */
		func() bool {
			{

				add(RuleAction41, position)
			}
			return true
		},
		/* 96 Action42 <- <{ currentAST := MakeASTNode(NODE_IF, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine); yy = currentAST; //If:2 No ELSE }> */
		func() bool {
			{

				add(RuleAction42, position)
			}
			return true
		},
		/* 97 Action43 <- <{ stack[stack_idx-0] = nil  //Unless:0 stack[stack_idx-0] = 0 }> */
		func() bool {
			{

				add(RuleAction43, position)
			}
			return true
		},
		/* 98 Action44 <- <{ currentAST := MakeASTNode(NODE_UNLESS, stack[stack_idx-2], stack[stack_idx-1], stack[stack_idx-0], currentLine); yy = currentAST; //Unless:1 }> */
		func() bool {
			{

				add(RuleAction44, position)
			}
			return true
		},
		/* 99 Action45 <- <{ currentAST := MakeASTNode(NODE_UNLESS, stack[stack_idx-2], stack[stack_idx-1], nil, currentLine); yy = currentAST; //Unless:2 No ELSE }> */
		func() bool {
			{

				add(RuleAction45, position)
			}
			return true
		},
		/* 100 Action46 <- <{ yy = stack[stack_idx-0] //Else:0 }> */
		func() bool {
			{

				add(RuleAction46, position)
			}
			return true
		},
		/* 101 Action47 <- <{ msgAST := MakeASTNode(NODE_MSG, stack[stack_idx-1], nil, nil, currentLine); sendAST := MakeASTNode(NODE_SEND, nil, msgAST, nil, currentLine); currentAST := MakeASTNode(NODE_METHOD, sendAST, stack[stack_idx-0], nil, currentLine); yy = currentAST; //Method:0 NODE_METHOD }> */
		func() bool {
			{

				add(RuleAction47, position)
			}
			return true
		},
		/* 102 Action48 <- <{ currentAST := MakeASTNode(NODE_METHOD, stack[stack_idx-1], stack[stack_idx-0], nil, currentLine); yy = currentAST; //Method:1 }> */
		func() bool {
			{

				add(RuleAction48, position)
			}
			return true
		},
		/* 103 Action49 <- <{ currentAST := MakeASTNode(NODE_METHOD, nil, stack[stack_idx-0], nil, currentLine); yy = currentAST; //Method:2 Call top self }> */
		func() bool {
			{

				add(RuleAction49, position)
			}
			return true
		},
		/* 104 Action50 <- <{ stack[stack_idx-1] = nil //Def:0 stack[stack_idx-1]=0}> */
		func() bool {
			{

				add(RuleAction50, position)
			}
			return true
		},
		/* 105 Action51 <- <{ currentAST := MakeASTNode(NODE_DEF, stack[stack_idx-2], stack[stack_idx-1], stack[stack_idx-0], currentLine); yy = currentAST;  //Def:1 stack[stack_idx-1]>0}> */
		func() bool {
			{

				add(RuleAction51, position)
			}
			return true
		},
		/* 106 Action52 <- <{ //Params:0 stack[stack_idx-1] }> */
		func() bool {
			{

				add(RuleAction52, position)
			}
			return true
		},
		/* 107 Action53 <- <{ stack[stack_idx-1].PushBack(stack[stack_idx-0]) //Params:1 }> */
		func() bool {
			{

				add(RuleAction53, position)
			}
			return true
		},
		/* 108 Action54 <- <{ yy = stack[stack_idx-1] //Params:2 yy = stack[stack_idx-1] }> */
		func() bool {
			{

				add(RuleAction54, position)
			}
			return true
		},
		/* 109 Action55 <- <{ yy = MakeASTNode(NODE_PARAM, stack[stack_idx-1], nil, stack[stack_idx-0], currentLine) //Param:0 }> */
		func() bool {
			{

				add(RuleAction55, position)
			}
			return true
		},
		/* 110 Action56 <- <{ yy = MakeASTNode(NODE_PARAM, stack[stack_idx-1], nil, nil, currentLine)  //Param:1 }> */
		func() bool {
			{

				add(RuleAction56, position)
			}
			return true
		},
		/* 111 Action57 <- <{ //Param:2 splat TODO}> */
		func() bool {
			{

				add(RuleAction57, position)
			}
			return true
		},
		/* 112 Action58 <- <{ //Class:0 stack[stack_idx-1] = 0 }> */
		func() bool {
			{

				add(RuleAction58, position)
			}
			return true
		},
		/* 113 Action59 <- <{ yy = MakeASTNode(NODE_CLASS, stack[stack_idx-2], stack[stack_idx-1], stack[stack_idx-0], currentLine)  //Class:1 }> */
		func() bool {
			{

				add(RuleAction59, position)
			}
			return true
		},
		/* 114 Action60 <- <{ yy = MakeASTNode(NODE_MODULE, stack[stack_idx-1], nil, stack[stack_idx-0], currentLine)  //Module:0 }> */
		func() bool {
			{

				add(RuleAction60, position)
			}
			return true
		},
		/* 115 Action61 <- <{ //Range:0 TODO}> */
		func() bool {
			{

				add(RuleAction61, position)
			}
			return true
		},
		/* 116 Action62 <- <{ //Range:1 TODO}> */
		func() bool {
			{

				add(RuleAction62, position)
			}
			return true
		},
		/* 117 Action63 <- <{  yy = MakeASTNode(NODE_YIELD, stack[stack_idx-0],nil,nil, currentLine)  //Yield:0 }> */
		func() bool {
			{

				add(RuleAction63, position)
			}
			return true
		},
		/* 118 Action64 <- <{ yy = MakeASTNode(NODE_YIELD, stack[stack_idx-0],nil,nil, currentLine) //Yield:1 }> */
		func() bool {
			{

				add(RuleAction64, position)
			}
			return true
		},
		/* 119 Action65 <- <{ yy = MakeASTNode(NODE_YIELD, nil ,nil,nil, currentLine) //Yield:2 }> */
		func() bool {
			{

				add(RuleAction65, position)
			}
			return true
		},
		/* 120 Action66 <- <{ yy = MakeASTNode(NODE_RETURN, stack[stack_idx-1],nil,nil, currentLine) //Return:0 }> */
		func() bool {
			{

				add(RuleAction66, position)
			}
			return true
		},
		/* 121 Action67 <- <{ yy = MakeASTNode(NODE_RETURN, stack[stack_idx-1],nil,nil, currentLine) //Return:1 }> */
		func() bool {
			{

				add(RuleAction67, position)
			}
			return true
		},
		/* 122 Action68 <- <{ yy = MakeASTNode(NODE_RETURN, MakeASTNode(NODE_ARRAY, stack[stack_idx-0], nil,nil, currentLine) ,nil,nil, currentLine) //Return:2 }> */
		func() bool {
			{

				add(RuleAction68, position)
			}
			return true
		},
		/* 123 Action69 <- <{ yy = MakeASTNode(NODE_RETURN, MakeASTNode(NODE_ARRAY, stack[stack_idx-0], nil,nil, currentLine) ,nil,nil, currentLine) //Return:3 }> */
		func() bool {
			{

				add(RuleAction69, position)
			}
			return true
		},
		/* 124 Action70 <- <{ yy = MakeASTNode(NODE_RETURN, nil, nil,nil, currentLine) //Return:4 }> */
		func() bool {
			{

				add(RuleAction70, position)
			}
			return true
		},
		/* 125 Action71 <- <{ //Break:0 }> */
		func() bool {
			{

				add(RuleAction71, position)
			}
			return true
		},
		/* 126 Action72 <- <{ yy = MakeASTNode(NODE_VALUE, stack[stack_idx-2], nil, nil, currentLine) //Value:0 NUMBER }> */
		func() bool {
			{

				add(RuleAction72, position)
			}
			return true
		},
		/* 127 Action73 <- <{ yy = MakeASTNode(NODE_VALUE, stack[stack_idx-2], nil, nil, currentLine)//Value:1 SYMBOL }> */
		func() bool {
			{

				add(RuleAction73, position)
			}
			return true
		},
		/* 128 Action74 <- <{ yy = MakeASTNode(NODE_STRING, stack[stack_idx-2], nil, nil, currentLine) //Value:3 STRING1 }> */
		func() bool {
			{

				add(RuleAction74, position)
			}
			return true
		},
		/* 129 Action75 <- <{ yy = MakeASTNode(NODE_CONST, stack[stack_idx-2], nil, nil, currentLine) //Value:5 CONST }> */
		func() bool {
			{

				add(RuleAction75, position)
			}
			return true
		},
		/* 130 Action76 <- <{ yy = MakeASTNode(NODE_NIL, nil, nil, nil, currentLine) //Value:6 nil }> */
		func() bool {
			{

				add(RuleAction76, position)
			}
			return true
		},
		/* 131 Action77 <- <{ yy = MakeASTNode(NODE_TRUE, nil, nil, nil, currentLine) //Value:7 true }> */
		func() bool {
			{

				add(RuleAction77, position)
			}
			return true
		},
		/* 132 Action78 <- <{ yy = MakeASTNode(NODE_FALSE, nil, nil, nil, currentLine) //Value:8 false }> */
		func() bool {
			{

				add(RuleAction78, position)
			}
			return true
		},
		/* 133 Action79 <- <{ yy = MakeASTNode(NODE_SELF, nil, nil, nil, currentLine) //Value:9 self }> */
		func() bool {
			{

				add(RuleAction79, position)
			}
			return true
		},
		/* 134 Action80 <- <{ yy = MakeASTNode(NODE_GETIVAR, stack[stack_idx-1], nil, nil, currentLine) //Value:10 IVAR }> */
		func() bool {
			{

				add(RuleAction80, position)
			}
			return true
		},
		/* 135 Action81 <- <{ yy = MakeASTNode(NODE_GETCVAR, stack[stack_idx-1], nil, nil, currentLine) //Value:11 CVAR }> */
		func() bool {
			{

				add(RuleAction81, position)
			}
			return true
		},
		/* 136 Action82 <- <{ yy = MakeASTNode(NODE_GETGLOBAL, stack[stack_idx-1], nil, nil, currentLine) //Value:12 GLOBAL }> */
		func() bool {
			{

				add(RuleAction82, position)
			}
			return true
		},
		/* 137 Action83 <- <{ yy = MakeASTNode(NODE_ARRAY, nil, nil, nil, currentLine) //Value:13 [] }> */
		func() bool {
			{

				add(RuleAction83, position)
			}
			return true
		},
		/* 138 Action84 <- <{ yy = MakeASTNode(NODE_ARRAY, stack[stack_idx-0], nil, nil, currentLine) //Value:14 [AryItems] }> */
		func() bool {
			{

				add(RuleAction84, position)
			}
			return true
		},
		/* 139 Action85 <- <{ yy = MakeASTNode(NODE_HASH, nil, nil, nil, currentLine) //Value:15 {} }> */
		func() bool {
			{

				add(RuleAction85, position)
			}
			return true
		},
		/* 140 Action86 <- <{ yy = MakeASTNode(NODE_HASH, stack[stack_idx-0], nil, nil, currentLine) //Value:16 {HashItems} }> */
		func() bool {
			{

				add(RuleAction86, position)
			}
			return true
		},
		/* 141 Action87 <- <{ //AryItems:0 }> */
		func() bool {
			{

				add(RuleAction87, position)
			}
			return true
		},
		/* 142 Action88 <- <{ stack[stack_idx-1].PushBack(stack[stack_idx-0]) //AryItems:1 }> */
		func() bool {
			{

				add(RuleAction88, position)
			}
			return true
		},
		/* 143 Action89 <- <{ yy = stack[stack_idx-1] //AryItems:2 }> */
		func() bool {
			{

				add(RuleAction89, position)
			}
			return true
		},
		/* 144 Action90 <- <{ stack[stack_idx-2].PushBack(stack[stack_idx-1]) //HashItems:0 }> */
		func() bool {
			{

				add(RuleAction90, position)
			}
			return true
		},
		/* 145 Action91 <- <{ stack[stack_idx-2].PushBack(stack[stack_idx-0]) //HashItems:1 stack[stack_idx-0] }> */
		func() bool {
			{

				add(RuleAction91, position)
			}
			return true
		},
		/* 146 Action92 <- <{ stack[stack_idx-2].PushBack(stack[stack_idx-1]) //HashItems:2 stack[stack_idx-1] }> */
		func() bool {
			{

				add(RuleAction92, position)
			}
			return true
		},
		/* 147 Action93 <- <{ yy = stack[stack_idx-2] //HashItems:3 yy = stack[stack_idx-2] }> */
		func() bool {
			{

				add(RuleAction93, position)
			}
			return true
		},
		nil,
		/* 149 Action94 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //ID:0 KEYWORD.([ }> */
		func() bool {
			{

				add(RuleAction94, position)
			}
			return true
		},
		/* 150 Action95 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //ID:1 KEYWORD NAME }> */
		func() bool {
			{

				add(RuleAction95, position)
			}
			return true
		},
		/* 151 Action96 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //ID:2 simply ID }> */
		func() bool {
			{

				add(RuleAction96, position)
			}
			return true
		},
		/* 152 Action97 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //CONST:0 }> */
		func() bool {
			{

				add(RuleAction97, position)
			}
			return true
		},
		/* 153 Action98 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //BINOP:0 }> */
		func() bool {
			{

				add(RuleAction98, position)
			}
			return true
		},
		/* 154 Action99 <- <{ //UNOP:0 }> */
		func() bool {
			{

				add(RuleAction99, position)
			}
			return true
		},
		/* 155 Action100 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //ASSIGN:0 }> */
		func() bool {
			{

				add(RuleAction100, position)
			}
			return true
		},
		/* 156 Action101 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //IVAR:0 }> */
		func() bool {
			{

				add(RuleAction101, position)
			}
			return true
		},
		/* 157 Action102 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //CVAR:0 }> */
		func() bool {
			{

				add(RuleAction102, position)
			}
			return true
		},
		/* 158 Action103 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //GLOBAL:0 }> */
		func() bool {
			{

				add(RuleAction103, position)
			}
			return true
		},
		/* 159 Action104 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.numeric, _ = strconv.Atoi(buffer[begin:end]); yy = currentAST; //NUMBER:0 }> */
		func() bool {
			{

				add(RuleAction104, position)
			}
			return true
		},
		/* 160 Action105 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = buffer[begin:end]; yy = currentAST; //SYMBOL:0 }> */
		func() bool {
			{

				add(RuleAction105, position)
			}
			return true
		},
		/* 161 Action106 <- <{ parseStr = make([]string, 12) //STRING1:0 STRING_START }> */
		func() bool {
			{

				add(RuleAction106, position)
			}
			return true
		},
		/* 162 Action107 <- <{ parseStr = append(parseStr, "\\'")//STRING1:1 escaped \' }> */
		func() bool {
			{

				add(RuleAction107, position)
			}
			return true
		},
		/* 163 Action108 <- <{ parseStr = append(parseStr, buffer[begin:end])//STRING1:2 content }> */
		func() bool {
			{

				add(RuleAction108, position)
			}
			return true
		},
		/* 164 Action109 <- <{ currentAST :=  &AST{Type: NODE_ASTVAL, line: currentLine}; currentAST.value.str = strings.Join(parseStr, ""); yy = currentAST;  //STRING1:3 }> */
		func() bool {
			{

				add(RuleAction109, position)
			}
			return true
		},
		/* 165 Action110 <- <{ currentLine += 1 //EOL:0 count line }> */
		func() bool {
			{

				add(RuleAction110, position)
			}
			return true
		},
	}
	p.rules = rules
}
