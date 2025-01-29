package compiler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

    "github.com/sneaky-potato/g4th/model"
	"github.com/sneaky-potato/g4th/constants"
	"github.com/sneaky-potato/g4th/util"
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

func TypeCheckingProgram(program []model.Operation) {

    var stack = new(util.Stack[typedOperand])
    var ip int = 0
    for ip < len(program) {
        op := program[ip]
        switch op.Op {
        case constants.OP_PUSH_INT:
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_PUSH_STR:
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
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
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for +")
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
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for -")
            }
        case constants.OP_MUL:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "*")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ != TYPE_INT {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for *")
            }
            if b.typ != TYPE_INT {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for *")
            }
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_MOD:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "divmod")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ != TYPE_INT {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for divmod")
            }
            if b.typ != TYPE_INT {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for divmod")
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
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for =")
            }
        case constants.OP_NE:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "!=")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_PTR) {
                stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
            } else {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for !=")
            }
        case constants.OP_IF:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "if")
            var a typedOperand
            a = stack.Pop()
            if a.typ != TYPE_BOOL {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for if")
            }
        case constants.OP_ELSE:
        case constants.OP_END:
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
            b := stack.Peek(1)
            stack.Push(b)
        case constants.OP_SHR:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "shr")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == TYPE_INT && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for shr")
            }
        case constants.OP_SHL:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "shl")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == TYPE_INT && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for shl")
            }
        case constants.OP_OR:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "|")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == TYPE_INT && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for |")
            }
        case constants.OP_AND:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "&")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == TYPE_INT && b.typ == TYPE_INT {
                stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
            } else {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for &")
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
            stack.Push(b)
            stack.Push(a)
            stack.Push(c)
        case constants.OP_LT:
            util.CheckNumberOfArguments(stack.Size(), 2, op, "<")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_PTR) {
                stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
            } else {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for <")
            }
        case constants.OP_GT:
            util.CheckNumberOfArguments(stack.Size(), 2, op, ">")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ == b.typ && (a.typ == TYPE_INT || a.typ == TYPE_PTR) {
                stack.Push(typedOperand{ TYPE_BOOL, op.FilePath, op.Row })
            } else {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for >")
            }
        case constants.OP_WHILE:
        case constants.OP_DO:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "do")
            var a typedOperand
            a = stack.Pop()
            if a.typ != TYPE_BOOL {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for do")
            }
        case constants.OP_MEM:
            stack.Push(typedOperand{ TYPE_PTR, op.FilePath, op.Row })
        case constants.OP_LOAD:
            util.CheckNumberOfArguments(stack.Size(), 1, op, ",")
            a := stack.Pop()
            if a.typ != TYPE_PTR {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for ,")
            }
            stack.Push(typedOperand{ TYPE_INT, op.FilePath, op.Row })
        case constants.OP_STORE:
            util.CheckNumberOfArguments(stack.Size(), 2, op, ".")
            var a, b typedOperand
            a = stack.Pop()
            b = stack.Pop()
            if a.typ != TYPE_INT || b.typ != TYPE_PTR {
                util.WarnWithError(op.FilePath, op.Row, "invalid arguments for .")
            }
        case constants.OP_SYSCALL3:
            util.CheckNumberOfArguments(stack.Size(), 1, op, "syscall3")
            _ = stack.Pop()
            _ = stack.Pop()
            _ = stack.Pop()
        }
        ip += 1
    }

}

func CompileToAsm(outputFilePath string, program []model.Operation) {
    out, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    defer out.Close()
    if err != nil {
        log.Fatal(err)
    }

    var strs []string = []string{}

    out.WriteString("segment .text\n")
    out.WriteString("dump:\n")
    out.WriteString("    push    rbp\n")
    out.WriteString("    mov     rbp, rsp\n")
    out.WriteString("    sub     rsp, 64\n")
    out.WriteString("    mov     QWORD [rbp-56], rdi\n")
    out.WriteString("    mov     QWORD [rbp-8], 1\n")
    out.WriteString("    mov     eax, 32\n")
    out.WriteString("    sub     rax, QWORD [rbp-8]\n")
    out.WriteString("    mov     BYTE [rbp-48+rax], 10\n")
    out.WriteString(".L2:\n")
    out.WriteString("    mov     rcx, QWORD [rbp-56]\n")
    out.WriteString("    mov     rdx, -3689348814741910323\n")
    out.WriteString("    mov     rax, rcx\n")
    out.WriteString("    mul     rdx\n")
    out.WriteString("    shr     rdx, 3\n")
    out.WriteString("    mov     rax, rdx\n")
    out.WriteString("    sal     rax, 2\n")
    out.WriteString("    add     rax, rdx\n")
    out.WriteString("    add     rax, rax\n")
    out.WriteString("    sub     rcx, rax\n")
    out.WriteString("    mov     rdx, rcx\n")
    out.WriteString("    mov     eax, edx\n")
    out.WriteString("    lea     edx, [rax+48]\n")
    out.WriteString("    mov     eax, 31\n")
    out.WriteString("    sub     rax, QWORD [rbp-8]\n")
    out.WriteString("    mov     BYTE [rbp-48+rax], dl\n")
    out.WriteString("    add     QWORD [rbp-8], 1\n")
    out.WriteString("    mov     rax, QWORD [rbp-56]\n")
    out.WriteString("    mov     rdx, -3689348814741910323\n")
    out.WriteString("    mul     rdx\n")
    out.WriteString("    mov     rax, rdx\n")
    out.WriteString("    shr     rax, 3\n")
    out.WriteString("    mov     QWORD [rbp-56], rax\n")
    out.WriteString("    cmp     QWORD [rbp-56], 0\n")
    out.WriteString("    jne     .L2\n")
    out.WriteString("    mov     eax, 32\n")
    out.WriteString("    sub     rax, QWORD [rbp-8]\n")
    out.WriteString("    lea     rdx, [rbp-48]\n")
    out.WriteString("    lea     rcx, [rdx+rax]\n")
    out.WriteString("    mov     rax, QWORD [rbp-8]\n")
    out.WriteString("    mov     rdx, rax\n")
    out.WriteString("    mov     rsi, rcx\n")
    out.WriteString("    mov     edi, 1\n")
    out.WriteString("    mov     rax, 1\n")
    out.WriteString("    syscall\n")
    out.WriteString("    nop\n")
    out.WriteString("    leave\n")
    out.WriteString("    ret\n")

    out.WriteString("global _start\n")
    out.WriteString("_start:\n")

    ip := 0
    if constants.COUNT_OPS != 36 {
        panic("Exhaustive handling in compilation")
    }

    for ip < len(program) {
        out.WriteString(fmt.Sprintf("addr_%d:\n", ip))
        operation := program[ip]

        switch operation.Op {
        case constants.OP_PUSH_INT:
            out.WriteString(fmt.Sprintf("    ;; -- push int %d --\n", operation.Value))
            out.WriteString(fmt.Sprintf("    push %d\n", operation.Value))
        case constants.OP_PUSH_STR:
            out.WriteString(fmt.Sprintf("    ;; -- push str %s --\n", operation.Value))
            val, _ := operation.Value.(string)
            val, _ = strconv.Unquote(`"` + val + `"`)
            out.WriteString(fmt.Sprintf("    mov rax, %d\n", len(val)))
            out.WriteString("    push rax\n")
            out.WriteString(fmt.Sprintf("    push str_%d\n", len(strs)))
            strs = append(strs, val)
        case constants.OP_PLUS:
            out.WriteString("    ;; -- plus --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    add rax, rbx\n")
            out.WriteString("    push rax\n")
        case constants.OP_MINUS:
            out.WriteString("    ;; -- minus --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    sub rbx, rax\n")
            out.WriteString("    push rbx\n")
        case constants.OP_MUL:
            out.WriteString("    ;; -- mul --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    mul rbx\n")
            out.WriteString("    push rax\n")
        case constants.OP_MOD:
            out.WriteString("    ;; -- mod --\n")
            out.WriteString("    xor rdx, rdx\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    div rbx\n")
            out.WriteString("    push rax\n")
            out.WriteString("    push rdx\n")
        case constants.OP_DUMP:
            out.WriteString("    ;; -- dump --\n")
            out.WriteString("    pop rdi\n")
            out.WriteString("    call dump\n")
        case constants.OP_EQ:
            out.WriteString("    ;; -- equal --\n")
            out.WriteString("    mov rcx, 0\n")
            out.WriteString("    mov rdx, 1\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    cmp rbx, rax\n")
            out.WriteString("    cmove rcx, rdx\n")
            out.WriteString("    push rcx\n")
        case constants.OP_NE:
            out.WriteString("    ;; -- not equal --\n")
            out.WriteString("    mov rcx, 0\n")
            out.WriteString("    mov rdx, 1\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    cmp rbx, rax\n")
            out.WriteString("    cmove rdx, rcx\n")
            out.WriteString("    push rdx\n")
        case constants.OP_IF:
            out.WriteString("    ;; -- if --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    test rax, rax\n")
            if operation.Jump < 0 {
                panic("`if` instruction does not have reference to end of its block, please use end after if")
            }
            out.WriteString(fmt.Sprintf("    jz addr_%d\n", operation.Jump))
        case constants.OP_ELSE:
            out.WriteString("    ;; -- else --\n")
            if operation.Jump < 0 {
                panic("`else` instruction does not have reference to end of its block, please use end after else")
            }
            out.WriteString(fmt.Sprintf("    jmp addr_%d\n", operation.Jump))
        case constants.OP_END:
            out.WriteString("    ;; -- end --\n")
            if operation.Jump < 0 {
                panic("`end` instruction does not have reference to next instruction")
            }
            if ip + 1 != operation.Jump {
                out.WriteString(fmt.Sprintf("    jmp addr_%d\n", operation.Jump))
            }
        case constants.OP_DUP:
            out.WriteString("    ;; -- dup --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    push rax\n")
            out.WriteString("    push rax\n")
        case constants.OP_2DUP:
            out.WriteString("    ;; -- 2dup --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    push rbx\n")
            out.WriteString("    push rax\n")
            out.WriteString("    push rbx\n")
            out.WriteString("    push rax\n")
        case constants.OP_DROP:
            out.WriteString("    ;; -- drop --\n")
            out.WriteString("    pop rax\n")
        case constants.OP_OVER:
            out.WriteString("    ;; -- over --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    push rbx\n")
            out.WriteString("    push rax\n")
            out.WriteString("    push rbx\n")
        case constants.OP_SHR:
            out.WriteString("    ;; -- shr --\n")
            out.WriteString("    pop rcx\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    shr rbx, cl\n")
            out.WriteString("    push rbx\n")
        case constants.OP_SHL:
            out.WriteString("    ;; -- shl --\n")
            out.WriteString("    pop rcx\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    shl rbx, cl\n")
            out.WriteString("    push rbx\n")
        case constants.OP_OR:
            out.WriteString("    ;; -- or --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    or rbx, rax\n")
            out.WriteString("    push rbx\n")
        case constants.OP_AND:
            out.WriteString("    ;; -- and --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    and rbx, rax\n")
            out.WriteString("    push rbx\n")
        case constants.OP_SWAP:
            out.WriteString("    ;; -- swap --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    push rax\n")
            out.WriteString("    push rbx\n")
        case constants.OP_ROT:
            out.WriteString("    ;; -- rot --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    pop rcx\n")
            out.WriteString("    push rax\n")
            out.WriteString("    push rcx\n")
            out.WriteString("    push rbx\n")
        case constants.OP_LT:
            out.WriteString("    ;; -- lt --\n")
            out.WriteString("    mov rcx, 0\n")
            out.WriteString("    mov rdx, 1\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    cmp rax, rbx\n")
            out.WriteString("    cmovl rcx, rdx\n")
            out.WriteString("    push rcx\n")
        case constants.OP_GT:
            out.WriteString("    ;; -- gt --\n")
            out.WriteString("    mov rcx, 0\n")
            out.WriteString("    mov rdx, 1\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    cmp rax, rbx\n")
            out.WriteString("    cmovg rcx, rdx\n")
            out.WriteString("    push rcx\n")
        case constants.OP_WHILE:
            out.WriteString("    ;; -- while --\n")
        case constants.OP_DO:
            out.WriteString("    ;; -- do --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    test rax, rax\n")
            if operation.Jump < 0 {
                panic("`do` instruction does not have reference to end of its block, please use end after else")
            }
            out.WriteString(fmt.Sprintf("    jz addr_%d\n", operation.Jump))
        case constants.OP_MEM:
            out.WriteString("    ;; -- mem --\n")
            out.WriteString("    push mem\n")
        case constants.OP_LOAD:
            out.WriteString("    ;; -- load --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    xor rbx, rbx\n")
            out.WriteString("    mov bl, [rax]\n")
            out.WriteString("    push rbx\n")
        case constants.OP_STORE:
            out.WriteString("    ;; -- store --\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    mov [rax], bl\n")
        case constants.OP_LOAD64:
            out.WriteString("    ;; -- load --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    xor rbx, rbx\n")
            out.WriteString("    mov rbx, [rax]\n")
            out.WriteString("    push rbx\n")
        case constants.OP_STORE64:
            out.WriteString("    ;; -- store --\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    mov [rax], rbx\n")
        case constants.OP_ARGC:
            panic("not implemented")
        case constants.OP_ARGV:
            panic("not implemented")
        case constants.OP_SYSCALL3:
            out.WriteString("    ;; -- syscall --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rdi\n")
            out.WriteString("    pop rsi\n")
            out.WriteString("    pop rdx\n")
            out.WriteString("    syscall\n")
        }

        ip += 1
    }
    out.WriteString(fmt.Sprintf("addr_%d:\n", len(program)))
    out.WriteString("    mov rax, 60\n")
    out.WriteString("    mov rdi, 0\n")
    out.WriteString("    syscall\n")
    out.WriteString("segment .data\n")
    for idx, s := range(strs) {
        out.WriteString(fmt.Sprintf("str_%d: ", idx))
        out.WriteString("db ")
        bytes := []byte(s)
        var stringHex []string = []string{}
        for _, b := range(bytes) {
            stringHex = append(stringHex, fmt.Sprintf("0x%x", b))
        }
        out.WriteString(strings.Join(stringHex, ","))
        out.WriteString("\n")
    }
    out.WriteString("segment .bss\n")
    out.WriteString(fmt.Sprintf("mem resb %d\n", constants.MEM_CAPACITY))
}
