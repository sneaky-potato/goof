package types

import (
    "fmt"

	"github.com/sneaky-potato/goof/util"
	"github.com/sneaky-potato/goof/model"
	"github.com/sneaky-potato/goof/constants"
)

const (
    TYPE_INT = iota
    TYPE_PTR
    TYPE_BOOL
)

type typedOperand struct {
    typ      int
    filePath string
    row      int
}

type proc struct {
    inputs  []int
    outputs []int
}

var procs = make(map[string]proc)

func stackEqual(s util.Stack[typedOperand], r util.Stack[typedOperand]) bool {
    i := 0
    if s.Size() != r.Size() {
        return false
    }
    for i < s.Size() {
        if s.Peek(i).typ != r.Peek(i).typ {
            return false
        }
        i += 1
    }
    return true
}

type blockStack struct {
    stack util.Stack[typedOperand]
    typ   int
}

func getStackString(stack util.Stack[typedOperand]) string {
    stackString := fmt.Sprintf("[ ")
    for stack.Size() > 0 {
        stackString += fmt.Sprintf("%s ", stack.Pop().getTypedString())
    }
    stackString += fmt.Sprintf(" ]")
    return stackString
}

func (tO typedOperand) getTypedString() string {
    typedString := fmt.Sprintf("%s:%d:", tO.filePath, tO.row)
    switch tO.typ {
    case TYPE_INT:
        typedString += "INT"
    case TYPE_PTR:
        typedString += "PTR"
    case TYPE_BOOL:
        typedString += "BOOL"
    default:
    }
    return typedString
}

func getStringFromOperands(typedOperands ...typedOperand) string {
    stringFromOperands := fmt.Sprintf("[ ")
    for _, tO := range typedOperands {
        stringFromOperands += tO.getTypedString()
        stringFromOperands += " "
    }
    stringFromOperands += "]"
    return stringFromOperands
}

func TypeCheckingProgram(program []model.Operation) {
    var stack = new(util.Stack[typedOperand])
    var blockStacks = new(util.Stack[blockStack])
    if constants.COUNT_OPS != 55 {
        panic("Exhaustive handling inside TypeCheckingProgram")
    }
    var op model.Operation
    for len(program) > 0 {
        op, program = program[0], program[1:]
        switch op.Op {
        case constants.OP_PUSH_INT:
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_PUSH_STR:
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
        case constants.OP_PUSH_PTR:
            stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
        case constants.OP_PLUS:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "+")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == TYPE_INT && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else if a.typ == TYPE_PTR && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
            } else if b.typ == TYPE_PTR && a.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for +\n" + foundArguments)
            }
        case constants.OP_MINUS:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "-")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == TYPE_INT && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else if a.typ == TYPE_PTR && b.typ == TYPE_PTR {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else if a.typ == TYPE_INT && b.typ == TYPE_PTR {
                stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for -\n" + foundArguments)
            }
        case constants.OP_MUL:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "*")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ != TYPE_INT || b.typ != TYPE_INT {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for *\n" + foundArguments)
            }
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_MOD:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "divmod")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ != TYPE_INT || b.typ != TYPE_INT {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for divmod\n" + foundArguments)
            }
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_DUMP:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "dump")
            _ = stack.Pop()
        case constants.OP_EQ:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "=")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_PTR) {
                stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for =\n" + foundArguments)
            }
        case constants.OP_NE:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "!=")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_PTR) {
                stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for !=\n" + foundArguments)
            }
        case constants.OP_IF:
            copyStack := stack.Copy()
            blockStacks.Push(blockStack {*copyStack, op.Op})
        case constants.OP_ELIF:
            var block blockStack
            block = blockStacks.Pop()
            if block.typ != constants.OP_DO {
                panic("elif can only be used after do")
            }
            copyStack := stack.Copy()
            blockStacks.Push(blockStack {*copyStack, op.Op})
            stack.Assign(block.stack)
        case constants.OP_ELSE:
            var block blockStack
            block = blockStacks.Pop()
            if block.typ != constants.OP_DO {
                panic("else can only be used after do")
            }
            copyStack := stack.Copy()
            blockStacks.Push(blockStack {*copyStack, op.Op})
            stack.Assign(block.stack)
        case constants.OP_WHILE:
            copyStack := stack.Copy()
            blockStacks.Push(blockStack {*copyStack, op.Op})
        case constants.OP_DO:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "do")
            var a typedOperand
            a = stack.Pop()
            if a.typ != TYPE_BOOL {
                foundArguments := getStringFromOperands(a)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for do\n" + foundArguments)
            }
            util.CheckNumberOfArguments(blockStacks.Size(), 1, op, "do")
            var block blockStack
            block = blockStacks.Pop()
            if block.typ != constants.OP_WHILE && block.typ != constants.OP_IF && block.typ != constants.OP_ELIF {
                panic("do must be used after if, elif or while")
            }
            isEqual := stackEqual(*stack, block.stack)
            if block.typ == constants.OP_WHILE && !isEqual {
                fmt.Printf("expected: %s\nactual: %s\n", getStackString(block.stack), getStackString(*stack))
                util.TerminateWithError(op.FilePath, op.Row, "while-do condition cannot modify the types on data stack")
            }
            copyStack := stack.Copy()
            blockStacks.Push(blockStack {*copyStack, op.Op})
        case constants.OP_END:
            util.CheckNumberOfArguments(blockStacks.Size(), 1, op, "end")
            var block blockStack
            block = blockStacks.Pop()
            if block.typ == constants.OP_ELSE {
                isEqual := stackEqual(*stack, block.stack)
                if !isEqual {
                    fmt.Printf("expected: %s\nactual: %s\n", getStackString(block.stack), getStackString(*stack))
                    util.TerminateWithError(op.FilePath, op.Row, "all branches of if-else must produce same type arguments on data stack")
                }
            } else if block.typ == constants.OP_DO {
                isEqual := stackEqual(*stack, block.stack)
                if !isEqual {
                    fmt.Printf("expected: %s\nactual: %s\n", getStackString(block.stack), getStackString(*stack))
                    util.TerminateWithError(op.FilePath, op.Row, "do-end block cannot modify the types on data stack")
                }
            } else {
                panic("unreachable")
            }
        case constants.OP_DUP:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "dup")
            var a typedOperand
            a = stack.Pop()
            stack.Push(a)
            stack.Push(a)
        case constants.OP_2DUP:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "2dup")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            stack.Push(b)
            stack.Push(a)
            stack.Push(b)
            stack.Push(a)
        case constants.OP_DROP:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "drop")
            _ = stack.Pop()
        case constants.OP_OVER:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "over")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            stack.Push(b)
            stack.Push(a)
            stack.Push(b)
        case constants.OP_SHR:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "shr")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == TYPE_INT && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for shr\n" + foundArguments)
            }
        case constants.OP_SHL:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "shl")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == TYPE_INT && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for shl\n" + foundArguments)
            }
        case constants.OP_OR:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "|")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_BOOL) {
                stack.Push(typedOperand{ a.typ, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for |\n" + foundArguments)
            }
        case constants.OP_AND:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "&")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_BOOL) {
                stack.Push(typedOperand{ a.typ, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for &\n" + foundArguments)
            }
        case constants.OP_SWAP:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "swap")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            stack.Push(a)
            stack.Push(b)
        case constants.OP_ROT:
            util.CheckNumberOfArguments(stack.Size(), 3, op, "rot")
            var a, b, c typedOperand
            a = stack.Pop()
            b = stack.Pop()
            c = stack.Pop()
            stack.Push(a)
            stack.Push(c)
            stack.Push(b)
        case constants.OP_LT:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "<")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_PTR) {
                stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for <\n" + foundArguments)
            }
        case constants.OP_GT:
            util.CheckNumberOfArguments(stack.Size(), 2, op, ">")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_PTR) {
                stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
            } else {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for >\n" + foundArguments)
            }
        case constants.OP_MEM:
            stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
        case constants.OP_LOAD:
            util.CheckNumberOfArguments(stack.Size(), 1, op, ",")
            a := stack.Pop()
            if a.typ != TYPE_PTR {
                foundArguments := getStringFromOperands(a)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for ,\n" + foundArguments)
            }
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_STORE:
            util.CheckNumberOfArguments(stack.Size(), 2, op, ".")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if (a.typ != TYPE_INT && a.typ != TYPE_PTR) || b.typ != TYPE_PTR {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for .\n" + foundArguments)
            }
        case constants.OP_LOAD64:
            util.CheckNumberOfArguments(stack.Size(), 1, op, ",64")
            a := stack.Pop()
            if a.typ != TYPE_PTR {
                foundArguments := getStringFromOperands(a)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for ,64\n" + foundArguments)
            }
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_STORE64:
            util.CheckNumberOfArguments(stack.Size(), 2, op, ".64")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if (a.typ != TYPE_INT && a.typ != TYPE_PTR) || b.typ != TYPE_PTR {
                foundArguments := getStringFromOperands(a, b)
                util.TerminateWithError(op.FilePath, op.Row, "invalid arguments for .64\n" + foundArguments)
            }
        case constants.OP_CAST_PTR:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "cast(ptr)")
            stack.Pop()
            stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
        case constants.OP_CAST_BOOL:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "cast(bool)")
            stack.Pop()
            stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
        case constants.OP_CAST_INT:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "cast(int)")
            stack.Pop()
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_PREP_PROC:
            var nextOp model.Operation
            var procNameString = op.Value.(string)
            procs[procNameString] = proc{ make([]int, 0), make([]int, 0) }
            for len(program) > 0 {
                nextOp, program = program[0], program[1:]
                if nextOp.Op == constants.OP_PROC_SEP {
                    break
                } else if nextOp.Op == constants.OP_TYPE_BOOL || nextOp.Op == constants.OP_TYPE_INT || nextOp.Op == constants.OP_TYPE_PTR {
                    if entry, ok := procs[procNameString]; ok {
                        entry.inputs = append(entry.inputs, nextOp.Op)
                        procs[procNameString] = entry
                    }
                } else {
                    util.TerminateWithError(op.FilePath, op.Row, "expected type signature for procedure call\n")
                }
            }

            for len(program) > 0 {
                nextOp, program = program[0], program[1:]
                if nextOp.Op == constants.OP_PROC_SEP {
                    break
                } else if nextOp.Op == constants.OP_TYPE_BOOL || nextOp.Op == constants.OP_TYPE_INT || nextOp.Op == constants.OP_TYPE_PTR {
                    if entry, ok := procs[procNameString]; ok {
                        entry.outputs = append(entry.outputs, nextOp.Op)
                        procs[procNameString] = entry
                    }
                } else {
                    util.TerminateWithError(op.FilePath, op.Row, "expected type signature for procedure call\n")
                }
            }

            funcStack := 1

            for len(program) > 0 {
                nextOp, program = program[0], program[1:]
                if nextOp.Op == constants.OP_END {
                    funcStack -= 1
                    if funcStack == 0 {
                        break
                    }
                } else if nextOp.Op == constants.OP_IF || nextOp.Op == constants.OP_WHILE {
                    funcStack += 1
                }
            }
        case constants.OP_CALL:
            util.CheckNumberOfArguments(stack.Size(), len(procs[op.Value.(string)].inputs), op, op.Value.(string))
            for _, in := range procs[op.Value.(string)].inputs {
                a := stack.Pop()
                if in == constants.OP_TYPE_INT && a.typ == TYPE_INT {
                } else if in == constants.OP_TYPE_PTR && a.typ == TYPE_PTR {
                } else if in == constants.OP_TYPE_BOOL && a.typ == TYPE_BOOL {
                } else {
                    foundArguments := a.getTypedString()
                    util.TerminateWithError(op.FilePath, op.Row, "unexpected input for procedure call: " + foundArguments + "\n")
                }
            }
            for in := range procs[op.Value.(string)].outputs {
                if in == constants.OP_TYPE_INT {
                    stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
                } else if in == constants.OP_TYPE_PTR {
                    stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
                } else if in == constants.OP_TYPE_BOOL {
                    stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
                } else {
                    panic("unreachable code")
                }
            }
        case constants.OP_ARGC:
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_ARGV:
            stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
        case constants.OP_SYSCALL1:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "syscall1")
            _ = stack.Pop()
            _ = stack.Pop()
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_SYSCALL2:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "syscall2")
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_SYSCALL3:
            util.CheckNumberOfArguments(stack.Size(), 4, op, "syscall3")
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_SYSCALL6:
            util.CheckNumberOfArguments(stack.Size(), 7, op, "syscall6")
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        default:
        }
    }
    if stack.Size() > 0 {
        errorString := fmt.Sprintf("unhandled data on stack: [ ")
        for stack.Size() > 0 {
            errorString += stack.Pop().getTypedString()
            errorString += " "
        }
        errorString += "]"
        util.TerminateWithError(program[len(program) - 1].FilePath, program[len(program) - 1].Row, errorString)
    }
}
