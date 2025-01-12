package compiler

import (
    "fmt"
	"log"
	"os"

	"github.com/sneaky-potato/g4th/constants"
	"github.com/sneaky-potato/g4th/lexer"
)

func CompileToAsm(outputFilePath string, program []lexer.Operation) {
    out, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    defer out.Close()
    if err != nil {
        log.Fatal(err)
    }

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
    if constants.COUNT_OPS != 8 {
        panic("Exhaustive handling in compilation")
    }

    for ip < len(program) {
        out.WriteString(fmt.Sprintf("addr_%d:\n", ip))
        operation := program[ip]

        switch operation.Op {
        case constants.OP_PUSH:
            out.WriteString(fmt.Sprintf("    ;; -- push %d --\n", operation.Value))
            out.WriteString(fmt.Sprintf("    push %d\n", operation.Value))
            ip += 1
        case constants.OP_PLUS:
            out.WriteString("    ;; -- plus --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    add rax, rbx\n")
            out.WriteString("    push rax\n")
            ip += 1
        case constants.OP_MINUS:
            out.WriteString("    ;; -- minus --\n")
            out.WriteString("    pop rax\n")
            out.WriteString("    pop rbx\n")
            out.WriteString("    sub rbx, rax\n")
            out.WriteString("    push rbx\n")
            ip += 1
        case constants.OP_DUMP:
            out.WriteString("    ;; -- dump --\n")
            out.WriteString("    pop rdi\n")
            out.WriteString("    call dump\n")
            ip += 1
        default:
            ip += 1
        }
    }
    out.WriteString(fmt.Sprintf("addr_%d:\n", len(program)))
    out.WriteString("    mov rax, 60\n")
    out.WriteString("    mov rdi, 0\n")
    out.WriteString("    syscall\n")
}
