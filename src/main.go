package main

import (
    "flag"
    "io/ioutil"
    "log"
    "runtime"
)

func main() {
    runtime.GOMAXPROCS(2)
    flag.Parse()

    if flag.NArg() != 1 {
        flag.Usage()
        log.Fatalf("FILE: the leg file to compile")
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
    
    Traverse(rootAST)
}
