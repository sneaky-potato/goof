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
- [x] Compiled
- [x] Native
- [x] Turing complete
- [ ] Add support for defining and calling functions
- [ ] Add library builtin functions
- [ ] Achieve bootstrapped compiler

## Language Reference

The language implements the following constructs

### Literals

#### Integer

Currently a sequence of digits which may optionally start with a dash (-) is interpreted as an integer.

```pascal
10 1 +
```

The code above pushes 10 and 1 to the stack, pop them, sums them up and then pushes the result on top of stack.

#### String

A string is a sequence of characters sandwiched between double quotes (").

```pascal
include "std.g4th"

"Hello World\n" write
```

When the compiler encounters a string the following happens:
1. the size of the string in bytes is pushed onto the stack,
2. the bytes of the string are copied somewhere into the memory,
3. the pointer to the beginning of the string is pushed onto the data stack.

### Intrinsics

#### stack

| Name    | Signature        | Description                                                                                  |
| ---     | ---              | ---                                                                                          |
| `dup`   | `a -- a a`       | duplicate an element on top of the stack.                                                    |
| `swap`  | `a b -- b a`     | swap 2 elements on the top of the stack.                                                     |
| `drop`  | `a b -- a`       | drops the top element of the stack.                                                          |
| `dump`  | `a b -- a`       | print the element on top of the stack in a free form to stdout and remove it from the stack. |
| `over`  | `a b -- a b a`   | copy the element below the top of the stack                                                  |
| `rot`   | `a b c -- b c a` | rotate the top three stack elements.                                                         |

#### Comparison

| Name | Signature                              | Description                                                  |
| ---  | ---                                    | ---                                                          |
| `= ` | `[a: int] [b: int] -- [a == b : bool]` | checks if two elements on top of the stack are equal.        |
| `> ` | `[a: int] [b: int] -- [a > b  : bool]` | applies the greater comparison on top two elements.          |
| `< ` | `[a: int] [b: int] -- [a < b  : bool]` | applies the less comparison on top two elements.             |

#### Arithmetic

| Name     | Signature                                        | Description                                                                                                              |
| ---      | ---                                              | ---                                                                                                                      |
| `+`      | `[a: int] [b: int] -- [a + b: int]`              | sums up two elements on the top of the stack.                                                                            |
| `-`      | `[a: int] [b: int] -- [a - b: int]`              | subtracts two elements on the top of the stack                                                                           |

#### Bitwise

| Name  | Signature                            | Description                   |
| ---   | ---                                  | ---                           |
| `shr` | `[a: int] [b: int] -- [a >> b: int]` | right **unsigned** bit shift. |
| `shl` | `[a: int] [b: int] -- [a << b: int]` | light bit shift.              |
| `or`  | `[a: int] [b: int] -- [a \| b: int]` | bit `or`.                     |
| `and` | `[a: int] [b: int] -- [a & b: int]`  | bit `and`.                    |

#### System

- `syscall<n>` - perform a syscall with n arguments where n is in range `[0..6]`. (`syscall1`, `syscall2`, etc)

```porth
syscall_number = pop()
<move syscall_number to the corresponding register>
for i in range(n):
    arg = pop()
    <move arg to i-th register according to the call convention>
<perform the syscall>
```
