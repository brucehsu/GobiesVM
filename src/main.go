package main

import (
	"flag"
	// "fmt"
	"io/ioutil"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("FILE: the .rb file to execute")
	}
	file := flag.Arg(0)

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

	// Traverse(rootAST)

	vm := initVM()

	vm.compile(rootAST)
	vm.executeBytecode()

	// fmt.Println("")
	// fmt.Println(len(vm.instList))
	// for _, v := range vm.instList {
	// 	fmt.Println(v)
	// 	fmt.Print("\t")
	// 	fmt.Println(v.obj)
	// }
}
