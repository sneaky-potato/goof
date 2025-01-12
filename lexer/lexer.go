package lexer

import (
	"fmt"
	"strconv"
    "bufio"
    "log"
    "os"
    "strings"

	"github.com/sneaky-potato/g4th/constants"
)

type Operation struct {
    Op    int
    Value int
    Jump  int
}

type Token struct {
    FilePath string
    Row      int
    Word     string
}

func pushInteger(value string) (Operation, error) {
    var intValue int; var err error

    if intValue, err = strconv.Atoi(value); err == nil {
        return Operation{ constants.OP_PUSH, intValue, -1 }, nil
    }
    return Operation{}, err
}

func ParseTokenAsOp(token Token) Operation {
    if constants.COUNT_OPS != 26 {
        panic("Exhaustive handling in parseTokenAsOp")
    }
    switch token.Word {
    case "+":
        return Operation{ constants.OP_PLUS, 0, -1 }
    case "-":
        return Operation{ constants.OP_MINUS, 0, -1 }
    case "=":
        return Operation{ constants.OP_EQUAL, 0, -1 }
    case "if":
        return Operation{ constants.OP_IF, 0, -1 }
    case "else":
        return Operation{ constants.OP_ELSE, 0, -1 }
    case "end":
        return Operation{ constants.OP_END, 0, -1 }
    case "dump":
        return Operation{ constants.OP_DUMP, 0, -1 }
    case "dup":
        return Operation{ constants.OP_DUP, 0, -1 }
    case "2dup":
        return Operation{ constants.OP_2DUP, 0, -1 }
    case "swap":
        return Operation{ constants.OP_SWAP, 0, -1 }
    case "rot":
        return Operation{ constants.OP_ROT, 0, -1 }
    case "drop":
        return Operation{ constants.OP_DROP, 0, -1 }
    case "over":
        return Operation{ constants.OP_OVER, 0, -1 }
    case "shl":
        return Operation{ constants.OP_SHL, 0, -1 }
    case "shr":
        return Operation{ constants.OP_SHR, 0, -1 }
    case "or":
        return Operation{ constants.OP_OR, 0, -1 }
    case "and":
        return Operation{ constants.OP_AND, 0, -1 }
    case "<":
        return Operation{ constants.OP_LT, 0, -1 }
    case ">":
        return Operation{ constants.OP_GT, 0, -1 }
    case "while":
        return Operation{ constants.OP_WHILE, 0, -1 }
    case "do":
        return Operation{ constants.OP_DO, 0, -1 }
    case "mem":
        return Operation{ constants.OP_MEM, 0, -1 }
    case ",":
        return Operation{ constants.OP_LOAD, 0, -1 }
    case ".":
        return Operation{ constants.OP_STORE, 0, -1 }
    case "syscall3":
        return Operation{ constants.OP_SYSCALL3, 0, -1 }
    default:
        operation, err := pushInteger(token.Word)
        if err == nil {
            return operation
        }
        errorString := fmt.Sprintf("%s:%d: %s", token.FilePath, token.Row, err)
		panic(errorString)
    }
}


func crossreferenceBlocks(program []Operation) []Operation {
    var stack []int
    var n int = 0
    if constants.COUNT_OPS != 26 {
        panic("Exhaustive handling inside crossreference")
    }
    ip := 0
    for ip < len(program) {
        op := program[ip]
        if op.Op == constants.OP_IF {
            stack = append(stack, ip)
            n += 1
        } else if op.Op == constants.OP_ELSE {
            if_ip := stack[n - 1]
            stack = stack[:n - 1]
            n -= 1
            if program[if_ip].Op != constants.OP_IF {
                panic("`else` can only be used inside `if` blocks")
            }
            // # ip + 1 so that it doesn't jump to else but rather body of else
            program[if_ip] = Operation{ constants.OP_IF, 0, ip + 1 }
            stack = append(stack, ip)
            n += 1
        } else if op.Op == constants.OP_END {
            block_ip := stack[n - 1]
            stack = stack[:n - 1]
            n -= 1
            if program[block_ip].Op == constants.OP_IF || program[block_ip].Op == constants.OP_ELSE {
                program[block_ip] = Operation{ program[block_ip].Op, 0, ip }
                program[ip] = Operation{ constants.OP_END, 0, ip + 1 }
            } else if program[block_ip].Op == constants.OP_DO {
                program[ip] = Operation{ constants.OP_END, 0, program[block_ip].Jump }
                program[block_ip] = Operation{ constants.OP_DO, 0, ip + 1 }
            } else {
                panic("end can only close `if` `else` `do` blocks for now")
            }
        } else if op.Op == constants.OP_WHILE {
            stack = append(stack, ip)
            n += 1
        } else if op.Op == constants.OP_DO {
            while_ip := stack[n - 1]
            stack = stack[:n - 1]
            n -= 1
            program[ip] = Operation{ constants.OP_DO, 0, while_ip }
            stack = append(stack, ip)
            n += 1
        }
        ip += 1
    }
    return program
}

func LoadProgramFromFile(filePath string) []Operation {
    var program []Operation
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    row := 0

    for scanner.Scan() {
        text := scanner.Text()
        text = strings.Split(text, "//")[0]
        words := strings.Fields(text)
        for _, word := range words {
            operation := ParseTokenAsOp(Token{ filePath, row, word })
            program = append(program, operation)
        }
        row += 1
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    program = crossreferenceBlocks(program)

    return program
}
