package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

const (
	LOG_COMMIT_SUCCESS int = iota
	LOG_COMMIT_VALID_READ
	LOG_COMMIT_LOCK_WRITE
	LOG_COMMIT_REVALID_READ
	LOG_EXEC_VALID_READ
)

type transactionStats struct {
	status [5]atomic.Int64
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("gobiesvm", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var printAST, printInst, logTransaction bool
	var bench int
	flags.BoolVar(&printAST, "ast", false, "Print abstract syntax tree structure")
	flags.BoolVar(&printInst, "bytecode", false, "Print comprehensive bytecode instructions")
	flags.BoolVar(&logTransaction, "log", false, "Log transaction status")
	flags.IntVar(&bench, "bench", 0, "Benchmark script execution speed (without parsing stage)")

	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		flags.Usage()
		return fmt.Errorf("FILE: expected one .rb file to execute")
	}

	buffer, err := os.ReadFile(flags.Arg(0))
	if err != nil {
		return fmt.Errorf("read %q: %w", flags.Arg(0), err)
	}

	resetParserState()
	p := &Leg{Buffer: string(buffer)}
	p.Init()
	if err := p.Parse(); err != nil {
		return fmt.Errorf("parse %q: %w", flags.Arg(0), err)
	}
	p.Execute()

	if printAST {
		Traverse(stdout, rootAST)
	}

	stats := &transactionStats{}
	if bench > 0 {
		start := time.Now()
		for range bench {
			vm := initVM(stdout, logTransaction, stats)
			vm.compile(rootAST)
			if err := vm.execute(); err != nil {
				return err
			}
		}

		if logTransaction {
			printTransactionLog(stdout, stats)
		}

		total := time.Since(start).Seconds()
		fmt.Fprintf(stdout, "Total time in %d iterations: %v s\n", bench, total)
		fmt.Fprintf(stdout, "Average Time in %d iterations: %v s\n", bench, total/float64(bench))
		return nil
	}

	vm := initVM(stdout, logTransaction, stats)
	vm.compile(rootAST)

	if printInst {
		printInstructions(stdout, vm.instList, true)
	}

	if err := vm.execute(); err != nil {
		return err
	}
	if logTransaction {
		printTransactionLog(stdout, stats)
	}
	return nil
}

func resetParserState() {
	traverse_level = 0
	currentLine = 0
	rootAST = nil
	parseStr = nil
}

func printInstructions(out io.Writer, inst []Instruction, blocks bool) {
	bytecodeNames := []string{
		"BC_PUTSELF",
		"BC_PUTNIL",
		"BC_PUTOBJ",
		"BC_PUTFIXNUM",
		"BC_PUTSTRING",
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
		"BC_INIT_THREAD",
	}
	for _, instruction := range inst {
		fmt.Fprintln(out, bytecodeNames[instruction.inst_type], instruction)
		fmt.Fprint(out, "\t")
		fmt.Fprintln(out, instruction.obj)
		if blocks && (instruction.inst_type == BC_PUTOBJ || instruction.inst_type == BC_INIT_THREAD) {
			if obj, ok := instruction.obj.(*RObject); ok && obj.name == "RBlock" {
				printInstructions(out, obj.methods["def"].def, blocks)
			}
		}
	}
}

func printTransactionLog(out io.Writer, stats *transactionStats) {
	fmt.Fprintln(out, "\nTransaction Logging Status")
	fmt.Fprintln(out, "Success:", stats.status[LOG_COMMIT_SUCCESS].Load())
	fmt.Fprintln(out, "Fail:")
	fmt.Fprintln(out, "\tCommit")
	fmt.Fprintln(out, "\t\tREAD:", stats.status[LOG_COMMIT_VALID_READ].Load())
	fmt.Fprintln(out, "\t\tREAD (revalid):", stats.status[LOG_COMMIT_REVALID_READ].Load())
	fmt.Fprintln(out, "\t\tWRITE:", stats.status[LOG_COMMIT_LOCK_WRITE].Load())
	fmt.Fprintln(out, "\tExecution")
	fmt.Fprintln(out, "\t\tREAD:", stats.status[LOG_EXEC_VALID_READ].Load())

	var total int64
	for i := range stats.status {
		total += stats.status[i].Load()
	}
	fmt.Fprintln(out, "Total:", total)
}
