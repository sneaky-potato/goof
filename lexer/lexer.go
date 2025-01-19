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
    Type     int
    Value    interface{}
    Expanded int
}

type Token struct {
    FilePath  string
    Row       int
    TokenWord Word
}

func terminateWithError(filePath string, row int, err string) {
    errorString := fmt.Sprintf("%s:%d -- %s\n", filePath, row, err)
    fmt.Fprintf(os.Stderr, errorString)
    os.Exit(1)
}

func ParseTokenAsOp(token Token) Operation {
    if constants.COUNT_OPS != 32 {
        panic("Exhaustive handling in parseTokenAsOp")
    }

    if token.TokenWord.Type == constants.TOKEN_WORD {
        val, ok := constants.BUILTIN_WORDS[token.TokenWord.Value.(string)]
        if ok {
            return Operation{ val, 0, -1 }
        } else {
            errorString := fmt.Sprintf("undefined token %s", token.TokenWord.Value)
            terminateWithError(token.FilePath, token.Row, errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_INT {
        val, ok := token.TokenWord.Value.(int)
        if ok {
            return Operation{ constants.OP_PUSH_INT, val, -1 }
        } else {
            errorString := fmt.Sprintf("undefined token %s", token.TokenWord.Value)
            terminateWithError(token.FilePath, token.Row, errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_STR {
        val, ok := token.TokenWord.Value.(string)
        if ok {
            return Operation{ constants.OP_PUSH_STR, val, -1 }
        } else {
            errorString := fmt.Sprintf("undefined token %s", token.TokenWord.Value)
            terminateWithError(token.FilePath, token.Row, errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_CHAR {
        val, ok := token.TokenWord.Value.(string)
        valByte := []byte(val)[0]
        if ok {
            return Operation{ constants.OP_PUSH_INT, valByte, -1 }
        } else {
            errorString := fmt.Sprintf("undefined token %s", token.TokenWord.Value)
            terminateWithError(token.FilePath, token.Row, errorString)
        }
    }
    panic("Unreachable code")
}

func expandMacro(macroTokens []Token, expanded int) []Token {
    for i := range macroTokens {
        macroTokens[i].TokenWord.Expanded = expanded
    }
    return macroTokens
}


func compileTokenList(tokenList []Token) []Operation {
    var stack []int
    var n int = 0
    var program []Operation
    macros := make(map[string][]Token)

    if constants.COUNT_OPS != 32 {
        panic("Exhaustive handling inside crossreference")
    }

    ip := 0
    var token Token
    for len(tokenList) > 0 {
        token, tokenList = tokenList[0], tokenList[1:]

        if token.TokenWord.Type == constants.TOKEN_WORD {
            if val, ok := macros[token.TokenWord.Value.(string)]; ok {

                if token.TokenWord.Expanded >= 35 {
                    errorString := fmt.Sprintf("Reached recursion limit %s", token.TokenWord.Value)
                    terminateWithError(token.FilePath, token.Row, errorString)
                }

                tokenList = append(expandMacro(val, token.TokenWord.Expanded + 1), tokenList...)
                continue
            }
        }

        op := ParseTokenAsOp(token)
        if op.Op == constants.OP_INCLUDE {
            if len(tokenList) == 0 {
                terminateWithError(token.FilePath, token.Row, "expected include file found nothing")
            }
            token, tokenList = tokenList[0], tokenList[1:]
            if token.TokenWord.Type != constants.TOKEN_STR {
                    errorString := fmt.Sprintf("expected include file to be string found %+v", token.TokenWord.Value)
                    terminateWithError(token.FilePath, token.Row, errorString)
            }
            if token.TokenWord.Expanded >= 35 {
                errorString := fmt.Sprintf("Reached recursion limit %s", token.TokenWord.Value)
                terminateWithError(token.FilePath, token.Row, errorString)
            }
            includedOperations := lexFile(token.TokenWord.Value.(string), token.TokenWord.Expanded + 1)
            tokenList = append(includedOperations, tokenList...)
        }
        if op.Op != constants.OP_MACRO {
            program = append(program, op)
        }
        if op.Op == constants.OP_IF {
            stack = append(stack, ip)
            n += 1
        } else if op.Op == constants.OP_ELSE {
            if n == 0 {
                terminateWithError(token.FilePath, token.Row, "`else` can only be used after `if`")
            }
            if_ip := stack[n - 1]
            stack = stack[:n - 1]
            n -= 1
            if program[if_ip].Op != constants.OP_IF {
                terminateWithError(token.FilePath, token.Row, "`else` can only be used after `if` block")
            }
            // # ip + 1 so that it doesn't jump to else but rather body of else
            program[if_ip] = Operation{ constants.OP_IF, 0, ip + 1 }
            stack = append(stack, ip)
            n += 1
        } else if op.Op == constants.OP_END {
            if n == 0 {
                terminateWithError(token.FilePath, token.Row, "`end` can only be used after `if` `else` `do`")
            }
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
                terminateWithError(token.FilePath, token.Row, "`end` can only close `if` `else` `do` blocks for now")
            }
        } else if op.Op == constants.OP_WHILE {
            stack = append(stack, ip)
            n += 1
        } else if op.Op == constants.OP_DO {
            if n == 0 {
                terminateWithError(token.FilePath, token.Row, "`do` can only be used after `while`")
            }
            while_ip := stack[n - 1]
            stack = stack[:n - 1]
            n -= 1

            if program[while_ip].Op != constants.OP_WHILE {
                terminateWithError(token.FilePath, token.Row, "`do` can only be used after `while`")
            }

            program[ip] = Operation{ constants.OP_DO, 0, while_ip }
            stack = append(stack, ip)
            n += 1
        } else if op.Op == constants.OP_MACRO {
            token, tokenList = tokenList[0], tokenList[1:]
            var macroName = token

            if macroName.TokenWord.Type != constants.TOKEN_WORD {
                errorString := fmt.Sprintf("expected macro name to be a word but found %+v", macroName.TokenWord.Value)
                terminateWithError(macroName.FilePath, macroName.Row, errorString)
            }

            macroNameString := macroName.TokenWord.Value.(string)

            _, ok := constants.BUILTIN_WORDS[macroNameString]
            if ok {
                errorString := fmt.Sprintf("redefinition of builtin word %+v", macroName.TokenWord.Value)
                terminateWithError(token.FilePath, token.Row, errorString)
            }

            _, ok = macros[macroNameString]
            if ok {
                errorString := fmt.Sprintf("redefinition of macro %+v", macroName.TokenWord.Value)
                terminateWithError(token.FilePath, token.Row, errorString)
            }

            var nextToken Token
            if len(tokenList) == 0 {
                errorString := fmt.Sprintf("macro definition incomplete %+v", macroName.TokenWord.Value)
                terminateWithError(token.FilePath, token.Row, errorString)
            }

            macroStack := 1

            for len(tokenList) > 0 {
                nextToken, tokenList = tokenList[0], tokenList[1:]

                if nextToken.TokenWord.Type == constants.TOKEN_WORD {
                    nextTokenString := nextToken.TokenWord.Value.(string)
                    if nextTokenString == "if" || nextTokenString == "while" {
                        macroStack += 1
                    } else if nextTokenString == "end" {
                        macroStack -= 1
                    }
                }

                if macroStack == 0 {
                    break
                }

                macros[macroName.TokenWord.Value.(string)] = append(macros[macroName.TokenWord.Value.(string)], nextToken)
            }

            if nextToken.TokenWord.Type != constants.TOKEN_WORD || nextToken.TokenWord.Value.(string) != "end" {
                errorString := fmt.Sprintf("macro definition incomplete %+v", macroName.TokenWord.Value)
                terminateWithError(token.FilePath, token.Row, errorString)
            }
            ip -= 1
        }

        ip += 1
    }
    return program
}

func lexWord(filePath string, row int, tokenWord string) Word {
    var intValue int; var err error

    if intValue, err = strconv.Atoi(tokenWord); err == nil {
        return Word{ constants.TOKEN_INT, intValue, 0 }
    }
    n := len(tokenWord)
    if n > 1 {
        first := tokenWord[0]
        last := tokenWord[n - 1]
        if first == '"' && last == '"' {
            return Word{ constants.TOKEN_STR, tokenWord[1:n-1], 0 }
        }

        if first == '\'' && last == '\'' {
            return Word{ constants.TOKEN_CHAR, tokenWord[1:n-1], 0 }
        }

        if  (first == '"' || first == '\'') && first != last {
            errorString := fmt.Sprintf("unclosed literal %s", tokenWord)
            terminateWithError(filePath, row, errorString)
        }
    }
    return Word{ constants.TOKEN_WORD, tokenWord, 0 }
}

var quoted bool = false
var escaped bool = false

func splitProgramWithStrings(r rune) bool {
    if (r == '"' || r == '\'') && !escaped {
        quoted = !quoted
    }

    if r == '\\' {
        escaped = true
    } else {
        escaped = false
    }

    return !quoted && unicode.IsSpace(r)
}

func lexFile(filePath string, expanded int) []Token {
    var tokenList []Token
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
            tokenWord := lexWord(filePath, row, word)
            tokenWord.Expanded = expanded
            tokenList = append(tokenList, Token{ filePath, row, tokenWord })
        }
        row += 1
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return tokenList
}
func LoadProgramFromFile(filePath string) []Operation {
    tokenList := lexFile(filePath, 0)
    program := compileTokenList(tokenList)
    return program
}
