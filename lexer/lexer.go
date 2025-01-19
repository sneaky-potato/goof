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

func ParseTokenAsOp(token Token) Operation {
    if constants.COUNT_OPS != 32 {
        panic("Exhaustive handling in parseTokenAsOp")
    }
    if token.TokenWord.Type == constants.TOKEN_WORD {
        val, ok := constants.BUILTIN_WORDS[token.TokenWord.Value.(string)]
        if ok {
            return Operation{ val, 0, -1 }
        } else {
            errorString := fmt.Sprintf("%s:%d -- undefined token %s", token.FilePath, token.Row, token.TokenWord.Value)
            panic(errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_INT {
        val, ok := token.TokenWord.Value.(int)
        if ok {
            return Operation{ constants.OP_PUSH_INT, val, -1 }
        } else {
            errorString := fmt.Sprintf("%s:%d -- undefined token %s", token.FilePath, token.Row, token.TokenWord.Value)
            panic(errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_STR {
        val, ok := token.TokenWord.Value.(string)
        if ok {
            return Operation{ constants.OP_PUSH_STR, val, -1 }
        } else {
            errorString := fmt.Sprintf("%s:%d -- undefined token %s", token.FilePath, token.Row, token.TokenWord.Value)
            panic(errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_CHAR {
        val, ok := token.TokenWord.Value.(string)
        valByte := []byte(val)[0]
        if ok {
            return Operation{ constants.OP_PUSH_INT, valByte, -1 }
        } else {
            errorString := fmt.Sprintf("%s:%d -- undefined token %s", token.FilePath, token.Row, token.TokenWord.Value)
            panic(errorString)
        }
    } else {
        panic("Unreachable code")
    }
}

func expandMacro(macroTokens []Token, expanded int) []Token {
    if expanded >= 35 {
        panic("Reached recursion limit")
    }
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
                tokenList = append(expandMacro(val, token.TokenWord.Expanded + 1), tokenList...)
                continue
            }
        }

        op := ParseTokenAsOp(token)
        if op.Op == constants.OP_INCLUDE {
            if len(tokenList) == 0 {
                panic(fmt.Sprintf("%s:%d -- expected include file found nothing", token.FilePath, token.Row))
            }
            token, tokenList = tokenList[0], tokenList[1:]
            if token.TokenWord.Type != constants.TOKEN_STR {
                panic(fmt.Sprintf("%s:%d -- expected include file to be string found %+v", token.FilePath, token.Row, token.TokenWord.Value))
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
        } else if op.Op == constants.OP_MACRO {
            token, tokenList = tokenList[0], tokenList[1:]
            var macroName = token

            if macroName.TokenWord.Type != constants.TOKEN_WORD {
                panic(fmt.Sprintf("%s:%d -- expected macro name to be a word but found %+v", macroName.FilePath, macroName.Row, macroName.TokenWord.Value))
            }

            _, ok := constants.BUILTIN_WORDS[macroName.TokenWord.Value.(string)]
            if ok {
                panic(fmt.Sprintf("%s:%d -- redefinition of builtin word %+v", macroName.FilePath, macroName.Row, macroName.TokenWord.Value))
            }

            _, ok = macros[macroName.TokenWord.Value.(string)]
            if ok {
                panic(fmt.Sprintf("%s:%d -- redefinition of macro %+v", macroName.FilePath, macroName.Row, macroName.TokenWord.Value))
            }

            var nextToken Token
            for len(tokenList) > 0 {
                nextToken, tokenList = tokenList[0], tokenList[1:]
                if nextToken.TokenWord.Type == constants.TOKEN_WORD && nextToken.TokenWord.Value.(string) == "end" {
                    break
                } else {
                    macros[macroName.TokenWord.Value.(string)] = append(macros[macroName.TokenWord.Value.(string)], nextToken)
                }
            }

            if nextToken.TokenWord.Type != constants.TOKEN_WORD || nextToken.TokenWord.Value.(string) != "end" {
                panic(fmt.Sprintf("%s:%d -- macro definition incomplete", macroName.FilePath, macroName.Row))
            }
            ip -= 1
        }

        ip += 1
    }
    return program
}

func lexWord(tokenWord string) Word {
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
    if expanded >= 35 {
        panic("Reached recursion limit")
    }
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
            tokenWord := lexWord(word)
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
