package constants

const (
    OP_PUSH_INT = iota
    OP_PUSH_STR

    OP_PLUS
    OP_MINUS
    OP_MUL
    OP_MOD

    OP_EQ
    OP_NE
    OP_GT
    OP_LT

    OP_IF
    OP_ELSE
    OP_END
    OP_WHILE
    OP_DO
    OP_MACRO

    OP_DUP
    OP_2DUP
    OP_SWAP
    OP_DROP
    OP_OVER
    OP_ROT

    OP_SHR
    OP_SHL
    OP_OR
    OP_AND

    OP_MEM
    OP_LOAD
    OP_STORE

    OP_SYSCALL3
    OP_DUMP
    OP_INCLUDE
    COUNT_OPS
)

var BUILTIN_WORDS = map[string]int{
    "+": OP_PLUS,
    "-": OP_MINUS,
    "*": OP_MUL,
    "mod": OP_MOD,
    "=": OP_EQ,
    "!=": OP_NE,
    "if": OP_IF,
    "else": OP_ELSE,
    "end": OP_END,
    "while": OP_WHILE,
    "do": OP_DO,
    "macro": OP_MACRO,
    "dup": OP_DUP,
    "2dup": OP_2DUP,
    "swap": OP_SWAP,
    "rot": OP_ROT,
    "drop": OP_DROP,
    "over": OP_OVER,
    "shl": OP_SHL,
    "shr": OP_SHR,
    "or": OP_OR,
    "and": OP_AND,
    "<": OP_LT,
    ">": OP_GT,
    "mem": OP_MEM,
    ",": OP_LOAD,
    ".": OP_STORE,
    "syscall3": OP_SYSCALL3,
    "dump": OP_DUMP,
    "include": OP_INCLUDE,
}

const MEM_CAPACITY = 640_000
const STR_CAPACITY = 640_000
