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
    if constants.COUNT_OPS != 8 {
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
    default:
        operation, err := pushInteger(token.Word)
        if err == nil {
            return operation
        }
        errorString := fmt.Sprintf("%s:%d: %s", token.FilePath, token.Row, err)
		panic(errorString)
    }
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
            operation := ParseTokenAsOp(Token{filePath, row, word})
            program = append(program, operation)
        }
        row += 1
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return program
}
