package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"time"
)

const (
	LOG_COMMIT_SUCCESS int = iota
	LOG_COMMIT_VALID_READ
	LOG_COMMIT_LOCK_WRITE
	LOG_COMMIT_REVALID_READ
	LOG_EXEC_VALID_READ
)

var log_transaction bool
var log_status [5]int32

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Flag initialization
	var printAST, printInst bool
	var bench int
	flag.BoolVar(&printAST, "ast", false, "Print abstract syntax tree structure")
	flag.BoolVar(&printInst, "bytecode", false, "Print comprehensive bytecode instructions")
	flag.BoolVar(&log_transaction, "log", false, "Log transaction status")
	flag.IntVar(&bench, "bench", 0, "Benchmark script execution speed (without parsing stage)")

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

	if bench > 0 {
		start := time.Now()
		for i := 0; i < bench; i++ {
			vm := initVM()
			vm.compile(rootAST)
			vm.execute()
		}

		if log_transaction {
			printTransactionLog()
		}

		total := float64(time.Since(start).Nanoseconds()) / float64(1000000000)
		fmt.Printf("Total time in %d iterations: %v s\n", bench, total)
		fmt.Printf("Average Time in %d iterations: %v s\n", bench, total/float64(bench))
		return
	}

	vm := initVM()

	vm.compile(rootAST)

	if printInst {
		printInstructions(vm.instList, true)
	}

	vm.execute()

	if log_transaction {
		printTransactionLog()
	}
}

func printInstructions(inst []Instruction, blocks bool) {
	for _, v := range inst {
		fmt.Println(v)
		fmt.Print("\t")
		fmt.Println(v.obj)
		if blocks && (v.inst_type == BC_PUTOBJ || v.inst_type == BC_INIT_THREAD) && v.obj.(*RObject).name == "RBlock" {
			printInstructions(v.obj.(*RObject).methods["def"].def, blocks)
		}
	}
}

func printTransactionLog() {
	total := int64(0)
	fmt.Println("\nTransaction Logging Status")
	fmt.Println("Success:", log_status[LOG_COMMIT_SUCCESS])
	fmt.Println("Fail:")
	fmt.Println("\tCommit")
	fmt.Println("\t\tREAD:", log_status[LOG_COMMIT_VALID_READ])
	fmt.Println("\t\tREAD (revalid):", log_status[LOG_COMMIT_REVALID_READ])
	fmt.Println("\t\tWRITE:", log_status[LOG_COMMIT_LOCK_WRITE])
	fmt.Println("\tExecution")
	fmt.Println("\t\tREAD:", log_status[LOG_EXEC_VALID_READ])

	for _, log_cnt := range log_status {
		total += int64(log_cnt)
	}

	fmt.Println("Total:", total)
}
