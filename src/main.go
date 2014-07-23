package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)

	// Flag initialization
	var printAST, printInst bool
	flag.BoolVar(&printAST, "ast", false, "Print abstract syntax tree structure")
	flag.BoolVar(&printInst, "bytecode", false, "Print comprehensive bytecode instructions")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("FILE: the .rb file to execute")
	}

	file := flag.Args()[0]

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	p := &Leg{Buffer: string(buffer)}
	p.Init()
	if err := p.Parse(); err != nil {
		log.Fatal(err)
	}

	p.Execute()

	if printAST {
		Traverse(rootAST)
	}

	vm := initVM()

	vm.compile(rootAST)

	if printInst {
		printInstructions(vm.instList, true)
	}

	vm.executeBytecode(nil)
}

func printInstructions(inst []Instruction, blocks bool) {
	BytecodeInString := []string{
		"BC_PUTSELF",
		"BC_PUTNIL",
		"BC_PUTOBJ",
		"BC_PUTTRUE",
		"BC_PUTFALSE",
		"BC_SETLOCAL",
		"BC_GETLOCAL",
		"BC_SETGLOBAL",
		"BC_GETGLOBAL",
		"BC_SETSYMBOL",
		"BC_GETSYMBOL",
		"BC_SETCONST",
		"BC_GETCONST",
		"BC_SETIVAR",
		"BC_GETIVAR",
		"BC_SETCVAR",
		"BC_GETCVAR",
		"BC_SEND",
		"BC_JUMP",
	}
	for _, v := range inst {
		fmt.Println(BytecodeInString[v.inst_type], v)
		fmt.Print("\t")
		fmt.Println(v.obj)

		if blocks && (v.inst_type == BC_PUTOBJ) && v.obj.(*RObject).name == "RBlock" {
			printInstructions(v.obj.(*RObject).methods["def"].def, blocks)
		}
	}
}
