# G4th

> I wanted to talk to the computer using my own language

goForth: a stack-based concatenative programming language inspired by [Forth](https://en.wikipedia.org/wiki/Forth_(programming_language)) and implemented in Go

## Usage

For simulating the program written in test.g4th
```shell
go run main.go sim ./test.g4th
```

For compiling the program written in test.g4th and writing to an ELF executable `output` (you can check the generated assembly in `output.asm`)
```shell
go run main.go com ./test.porth
./output
```

## TODOs
- Make the language turing complete
- Add support for defining and calling functions
- Add library builtin functions
- Achieve bootstrapped compiler
