GobiesVM
=====

[![Go](https://github.com/brucehsu/GobiesVM/actions/workflows/ci.yml/badge.svg)](https://github.com/brucehsu/GobiesVM/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://github.com/brucehsu/GobiesVM/blob/master/LICENSE)

A Ruby VM written in Go aims to exploit parallelism via Software Transactional Memory.

# Requirements

[Go 1.26.5](https://go.dev/doc/devel/release) or newer.

# Build

```
git clone https://github.com/brucehsu/GobiesVM.git
cd GobiesVM
make build
```

The executable will be placed at `bin/gobiesvm`.

# Development checks

```
make test
make race
make vet
make vuln
```

Run the complete local check suite with `make check`.

# Usage

```
bin/gobiesvm [OPTIONS] RBFILE

Options:
  -ast: Print abstract syntax tree structure
  -bench=[0]: Benchmark script execution speed (without parsing stage)
  -bytecode: Print comprehensive bytecode instructions
  -log: Log transaction status
```

# What is supported

- Object manipulations (creation, method calling)
- Variable assignment
- Integer
- Float
- String (single-quoted)
- Basic string operations
- Blocks
- Iterators
- Arrays
- Hashes
- Threads
- Basic IO

# What is NOT supported (yet)

- Conditional statements
- Class/method definition
- Native extension interface
- Rubygems

# Related Paper

> [Software Transactional Memory in a Dynamic Programming Language at Virtual Machine Level](http://handle.ncl.edu.tw/11296/ndltd/49813811049358903093)
