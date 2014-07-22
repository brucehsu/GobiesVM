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
	findTransactions(vm.instList)

	if printInst {
		printInstructions(vm.instList, true)
	}

	vm.executeBytecode(nil)
}

func printInstructions(inst []Instruction, blocks bool) {
	for _, v := range inst {
		fmt.Println(v)
		fmt.Print("\t")
		fmt.Println(v.obj)
		if (v.inst_type == BC_PUTOBJ || v.inst_type == BC_INITTRANS) && v.obj.(*RObject).name == "RBlock" {
			printInstructions(v.obj.(*RObject).methods["def"].def, blocks)
		}
	}
}
