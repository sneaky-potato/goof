package lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/sneaky-potato/g4th/constants"
)

type Operation struct {
    Op    int
    Value interface{}
    Jump  int
}

type Word struct {
    Type  int
    Value interface{}
}

type Token struct {
    FilePath  string
    Row       int
    TokenWord Word
}

func ParseTokenAsOp(token Token) Operation {
    if constants.COUNT_OPS != 27 {
        panic("Exhaustive handling in parseTokenAsOp")
    }
    if token.TokenWord.Type == constants.TOKEN_WORD {
        val, ok := constants.BUILTIN_WORDS[token.TokenWord.Value.(string)]
        if ok {
            return Operation{ val, 0, -1 }
        } else {
            errorString := fmt.Sprintf("%s:%d:%s -- %s", token.FilePath, token.Row, "undefined token", token.TokenWord.Value)
            panic(errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_INT {
        val, ok := token.TokenWord.Value.(int)
        if ok {
            return Operation{ constants.OP_PUSH_INT, val, -1 }
        } else {
            errorString := fmt.Sprintf("%s:%d:%s -- %s", token.FilePath, token.Row, "undefined token", token.TokenWord.Value)
            panic(errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_STR {
        val, ok := token.TokenWord.Value.(string)
        if ok {
            return Operation{ constants.OP_PUSH_STR, val, -1 }
        } else {
            errorString := fmt.Sprintf("%s:%d:%s -- %s", token.FilePath, token.Row, "undefined token", token.TokenWord.Value)
            panic(errorString)
        }
    } else {
        panic("Unreachable code")
    }
}


func crossreferenceBlocks(program []Operation) []Operation {
    var stack []int
    var n int = 0
    if constants.COUNT_OPS != 27 {
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

func lexWord(tokenWord string) Word {
    var intValue int; var err error

    if intValue, err = strconv.Atoi(tokenWord); err == nil {
        return Word{ constants.TOKEN_INT, intValue }
    }
    n := len(tokenWord)
    if n > 1 {
        first := tokenWord[0]
        last := tokenWord[n - 1]
        if first == '"' && last == '"' {
            return Word{ constants.TOKEN_STR, tokenWord[1:n-1]}
        }
    }
    return Word{ constants.TOKEN_WORD, tokenWord }
}

var quoted bool = false
var escaped bool = false

func splitProgramWithStrings(r rune) bool {
    if r == '"' && !escaped {
        quoted = !quoted
    }

    if r == '\\' {
        escaped = true
    } else {
        escaped = false
    }

    return !quoted && unicode.IsSpace(r)
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
        words := strings.FieldsFunc(text, splitProgramWithStrings)
        for _, word := range words {
            tokenWord := lexWord(word)
            operation := ParseTokenAsOp(Token{ filePath, row, tokenWord })
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
