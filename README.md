GobiesVM
=====

A Ruby VM written in Go aims to exploit parallelism via Software Transactional Memory.

# Build
```
git clone https://github.com/brucehsu/GobiesVM.git
cd GobiesVM
make goenv
source goenv
make
```

The executable will be placed inside ``bin/``.

# Usage
```
gobiesvm [OPTIONS] RBFILE

Options:
  -ast: Print abstract syntax tree structure
  -bench=[0]: Benchmark script execution speed (without parsing stage)
  -bytecode: Print comprehensive bytecode instructions
  -log: Log transaction status
```