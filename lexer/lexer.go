package lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/sneaky-potato/goof/constants"
    "github.com/sneaky-potato/goof/model"
	"github.com/sneaky-potato/goof/util"
)


var memory = make(map[string]int)

func ParseTokenAsOp(token model.Token) model.Operation {
    if constants.COUNT_OPS != 46 {
        panic("Exhaustive handling in parseTokenAsOp")
    }

    if token.TokenWord.Type == constants.TOKEN_WORD {
        if val, ok := constants.BUILTIN_WORDS[token.TokenWord.Value.(string)]; ok {
            if val == constants.OP_HERE {
                msg := token.FilePath + ":" + strconv.Itoa(token.Row)
                return model.Operation{ constants.OP_PUSH_STR, msg, -1, token.FilePath, token.Row }
            }
            return model.Operation{ val, 0, -1, token.FilePath, token.Row }
        } else if val, ok := memory[token.TokenWord.Value.(string)]; ok {
            return model.Operation{ constants.OP_PUSH_PTR, val, -1, token.FilePath, token.Row }
        } else {
            errorString := fmt.Sprintf("undefined token %s", token.TokenWord.Value)
            util.TerminateWithError(token.FilePath, token.Row, errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_INT {
        val, ok := token.TokenWord.Value.(int64)
        if ok {
            return model.Operation{ constants.OP_PUSH_INT, val, -1, token.FilePath, token.Row }
        } else {
            errorString := fmt.Sprintf("undefined token %s", token.TokenWord.Value)
            util.TerminateWithError(token.FilePath, token.Row, errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_STR {
        val, ok := token.TokenWord.Value.(string)
        if ok {
            return model.Operation{ constants.OP_PUSH_STR, val, -1, token.FilePath, token.Row }
        } else {
            errorString := fmt.Sprintf("undefined token %s", token.TokenWord.Value)
            util.TerminateWithError(token.FilePath, token.Row, errorString)
        }
    } else if token.TokenWord.Type == constants.TOKEN_CHAR {
        val, ok := token.TokenWord.Value.(string)
        valByte := []byte(val)[0]
        if ok {
            return model.Operation{ constants.OP_PUSH_INT, valByte, -1, token.FilePath, token.Row }
        } else {
            errorString := fmt.Sprintf("undefined token %s", token.TokenWord.Value)
            util.TerminateWithError(token.FilePath, token.Row, errorString)
        }
    }
    panic("Unreachable code")
}

func expandMacro(macroTokens []model.Token, expanded int) []model.Token {
    for i := range macroTokens {
        macroTokens[i].TokenWord.Expanded = expanded
    }
    return macroTokens
}


func compileTokenList(tokenList []model.Token) []model.Operation {
    var stack = new(util.Stack[int])
    var elifStack = new(util.Stack[int])
    var program []model.Operation
    macros := make(map[string][]model.Token)
    memoryPtr := 0

    if constants.COUNT_OPS != 46 {
        panic("Exhaustive handling inside compileTokenList")
    }

    ip := 0
    var token model.Token
    for len(tokenList) > 0 {
        token, tokenList = tokenList[0], tokenList[1:]

        if token.TokenWord.Type == constants.TOKEN_WORD {
            if val, ok := macros[token.TokenWord.Value.(string)]; ok {

                if token.TokenWord.Expanded >= 35 {
                    errorString := fmt.Sprintf("Reached recursion limit %s", token.TokenWord.Value)
                    util.TerminateWithError(token.FilePath, token.Row, errorString)
                }

                tokenList = append(expandMacro(val, token.TokenWord.Expanded + 1), tokenList...)
                continue
            }
        }

        op := ParseTokenAsOp(token)
        if op.Op == constants.OP_INCLUDE {
            if len(tokenList) == 0 {
                util.TerminateWithError(token.FilePath, token.Row, "expected include file, found nothing")
            }
            token, tokenList = tokenList[0], tokenList[1:]
            if token.TokenWord.Type != constants.TOKEN_STR {
                errorString := fmt.Sprintf("expected include file to be string, found %+v", token.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }
            if token.TokenWord.Expanded >= 35 {
                errorString := fmt.Sprintf("Reached recursion limit %s", token.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }
            includedOperations := lexFile(token.TokenWord.Value.(string), token.TokenWord.Expanded + 1)
            tokenList = append(includedOperations, tokenList...)
        }
        if op.Op != constants.OP_MACRO {
            program = append(program, op)
        }
        if op.Op == constants.OP_IF {
            stack.Push(ip)
        } else if op.Op == constants.OP_ELIF {
            if stack.Size() == 0 {
                util.TerminateWithError(token.FilePath, token.Row, "`elif` can only be used after `do`")
            }
            do_ip := stack.Pop()
            if program[do_ip].Op != constants.OP_DO {
                util.TerminateWithError(token.FilePath, token.Row, "`elif` can only be used after `do` block")
            }
            // # ip + 1 so that it doesn't jump to elif but rather body of elif
            program[do_ip].Jump = ip + 1
            stack.Push(ip)
            elifStack.Push(ip)
        } else if op.Op == constants.OP_ELSE {
            if stack.Size() == 0 {
                util.TerminateWithError(token.FilePath, token.Row, "`else` can only be used after `if-do` block")
            }
            do_ip := stack.Pop()
            if program[do_ip].Op != constants.OP_DO {
                util.TerminateWithError(token.FilePath, token.Row, "`else` can only be used after `do` block")
            }
            // # ip + 1 so that it doesn't jump to else but rather body of else
            program[do_ip].Jump = ip + 1
            stack.Push(ip)
        } else if op.Op == constants.OP_END {
            if stack.Size() == 0 {
                util.TerminateWithError(token.FilePath, token.Row, "`end` can only be used after `if` `else` `do`")
            }
            block_ip := stack.Pop()
            if program[block_ip].Op == constants.OP_ELSE {
                program[block_ip].Jump = ip
                program[ip].Jump = ip + 1
                for elifStack.Size() > 0 {
                    program[elifStack.Pop()].Jump = ip
                }
            } else if program[block_ip].Op == constants.OP_DO {
                if program[program[block_ip].Jump].Op == constants.OP_WHILE {
                    program[ip].Jump = program[block_ip].Jump
                    program[block_ip].Jump = ip + 1
                } else {
                    program[ip].Jump = ip + 1
                    program[block_ip].Jump = ip
                    for elifStack.Size() > 0 {
                        program[elifStack.Pop()].Jump = ip
                    }
                }
            } else {
                util.TerminateWithError(token.FilePath, token.Row, "`end` can only close `if` `else` `do` blocks for now")
            }
        } else if op.Op == constants.OP_WHILE {
            stack.Push(ip)
        } else if op.Op == constants.OP_DO {
            if stack.Size() == 0 {
                util.TerminateWithError(token.FilePath, token.Row, "`do` can only be used after `while`, `if`, `elif`")
            }
            while_if_ip := stack.Pop()

            if program[while_if_ip].Op == constants.OP_IF || program[while_if_ip].Op == constants.OP_ELIF {
                program[ip].Jump = while_if_ip
            } else if program[while_if_ip].Op == constants.OP_WHILE {
                program[ip].Jump = while_if_ip
            } else {
                fmt.Printf("debug: %+v\nwhile_ip=%+v\n", program[while_if_ip], while_if_ip)
                util.TerminateWithError(token.FilePath, token.Row, "`do` can only be used after `while`, `if`, `elif`")
            }

            stack.Push(ip)
        } else if op.Op == constants.OP_MEMORY {
            token, tokenList = tokenList[0], tokenList[1:]
            var memoryName = token

            if memoryName.TokenWord.Type != constants.TOKEN_WORD {
                errorString := fmt.Sprintf("expected memory name to be a word but found %+v", memoryName.TokenWord.Value)
                util.TerminateWithError(memoryName.FilePath, memoryName.Row, errorString)
            }

            memoryNameString := memoryName.TokenWord.Value.(string)

            _, ok := constants.BUILTIN_WORDS[memoryNameString]
            if ok {
                errorString := fmt.Sprintf("redefinition of builtin word %+v", memoryName.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }

            _, ok = memory[memoryNameString]
            if ok {
                errorString := fmt.Sprintf("redefinition of memory region %+v", memoryName.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }

            _, ok = macros[memoryNameString]
            if ok {
                errorString := fmt.Sprintf("redefinition of macro as a memory region %+v", memoryName.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }

            var nextToken model.Token
            if len(tokenList) == 0 {
                errorString := fmt.Sprintf("memory definition incomplete %+v", memoryName.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }

            memoryStack := new(util.Stack[int])

            for len(tokenList) > 0 {
                nextToken, tokenList = tokenList[0], tokenList[1:]

                if nextToken.TokenWord.Type == constants.TOKEN_WORD {
                    nextTokenString := nextToken.TokenWord.Value.(string)
                    if nextTokenString == "end" {
                        break
                    } else {
                        errorString := fmt.Sprintf("unsupported keyword inside memory definition %+v", memoryName.TokenWord.Value)
                        util.TerminateWithError(token.FilePath, token.Row, errorString)
                    }
                } else if nextToken.TokenWord.Type == constants.TOKEN_INT {
                    val := nextToken.TokenWord.Value.(int64)
                    memoryStack.Push(int(val))
                } else if nextToken.TokenWord.Type == constants.TOKEN_WORD {
                    tokenValue := nextToken.TokenWord.Value.(string)
                    if tokenValue == "+" {
                        if memoryStack.Size() < 2 {
                            errorString := fmt.Sprintf("operation + requires 2 arguments but found %d\n", memoryStack.Size())
                            util.TerminateWithError(token.FilePath, token.Row, errorString)
                        }
                        a := memoryStack.Pop() 
                        b := memoryStack.Pop()
                        memoryStack.Push(a + b)
                    } else if tokenValue == "*" {
                        if memoryStack.Size() < 2 {
                            errorString := fmt.Sprintf("operation * requires 2 arguments but found %d\n", memoryStack.Size())
                            util.TerminateWithError(token.FilePath, token.Row, errorString)
                        }
                        a := memoryStack.Pop() 
                        b := memoryStack.Pop()
                        memoryStack.Push(a * b)
                    } else {
                        errorString := fmt.Sprintf("unsupported operation inside memory definition %+v", memoryName.TokenWord.Value)
                        util.TerminateWithError(token.FilePath, token.Row, errorString)
                    }
                } else {
                    errorString := fmt.Sprintf("unsupported token inside memory definition %+v", memoryName.TokenWord.Value)
                    util.TerminateWithError(token.FilePath, token.Row, errorString)
                }

                if memoryStack.Size() < 1 {
                    errorString := fmt.Sprintf("memory definition leaves no value after evaluation")
                    util.TerminateWithError(token.FilePath, token.Row, errorString)
                }
                memorySize := memoryStack.Pop()
                memory[memoryNameString] = memoryPtr
                memoryPtr += int(memorySize)
            }
            ip -= 1

        } else if op.Op == constants.OP_MACRO {
            token, tokenList = tokenList[0], tokenList[1:]
            var macroName = token

            if macroName.TokenWord.Type != constants.TOKEN_WORD {
                errorString := fmt.Sprintf("expected macro name to be a word but found %+v", macroName.TokenWord.Value)
                util.TerminateWithError(macroName.FilePath, macroName.Row, errorString)
            }

            macroNameString := macroName.TokenWord.Value.(string)

            _, ok := constants.BUILTIN_WORDS[macroNameString]
            if ok {
                errorString := fmt.Sprintf("redefinition of builtin word %+v", macroName.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }

            _, ok = macros[macroNameString]
            if ok {
                errorString := fmt.Sprintf("redefinition of macro %+v", macroName.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }

            _, ok = macros[macroNameString]
            if ok {
                errorString := fmt.Sprintf("redefinition of memory region as a macro %+v", macroName.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }

            var nextToken model.Token
            if len(tokenList) == 0 {
                errorString := fmt.Sprintf("macro definition incomplete %+v", macroName.TokenWord.Value)
                util.TerminateWithError(token.FilePath, token.Row, errorString)
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
                util.TerminateWithError(token.FilePath, token.Row, errorString)
            }
            ip -= 1
        }

        ip += 1
    }
    return program
}

func lexWord(filePath string, row int, tokenWord string) model.Word {
    if intValue, err := strconv.ParseInt(tokenWord, 10, 64); err == nil {
        return model.Word{ constants.TOKEN_INT, intValue, 0 }
    }
    n := len(tokenWord)
    if n > 1 {
        first := tokenWord[0]
        last := tokenWord[n - 1]
        if first == '"' && last == '"' {
            return model.Word{ constants.TOKEN_STR, tokenWord[1:n-1], 0 }
        }

        if first == '\'' && last == '\'' {
            return model.Word{ constants.TOKEN_CHAR, tokenWord[1:n-1], 0 }
        }

        if  (first == '"' || first == '\'') && first != last {
            errorString := fmt.Sprintf("unclosed literal %s", tokenWord)
            util.TerminateWithError(filePath, row, errorString)
        }
    }
    return model.Word{ constants.TOKEN_WORD, tokenWord, 0 }
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

func lexFile(filePath string, expanded int) []model.Token {
    var tokenList []model.Token
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    row := 1

    for scanner.Scan() {
        text := scanner.Text()
        text = strings.Split(text, "//")[0]
        words := strings.FieldsFunc(text, splitProgramWithStrings)
        for _, word := range words {
            tokenWord := lexWord(filePath, row, word)
            tokenWord.Expanded = expanded
            tokenList = append(tokenList, model.Token{ filePath, row, tokenWord })
        }
        row += 1
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return tokenList
}
func LoadProgramFromFile(filePath string) []model.Operation {
    tokenList := lexFile(filePath, 0)
    program := compileTokenList(tokenList)
    return program
}
