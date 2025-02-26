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
    OP_LOAD64
    OP_STORE64

    OP_CAST_INT
    OP_CAST_PTR
    OP_CAST_BOOL

    OP_SYSCALL1
    OP_SYSCALL3
    OP_ARGV
    OP_ARGC
    OP_DUMP
    OP_INCLUDE
    OP_HERE
    COUNT_OPS
)

var BUILTIN_WORDS = map[string]int{
    "+": OP_PLUS,
    "-": OP_MINUS,
    "*": OP_MUL,
    "divmod": OP_MOD,
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
    "|": OP_OR,
    "&": OP_AND,
    "<": OP_LT,
    ">": OP_GT,
    "mem": OP_MEM,
    ",": OP_LOAD,
    ".": OP_STORE,
    ",64": OP_LOAD64,
    ".64": OP_STORE64,
    "(int)": OP_CAST_INT,
    "(ptr)": OP_CAST_PTR,
    "(bool)": OP_CAST_BOOL,
    "syscall1": OP_SYSCALL1,
    "syscall3": OP_SYSCALL3,
    "argv": OP_ARGV,
    "argc": OP_ARGC,
    "dump": OP_DUMP,
    "include": OP_INCLUDE,
    "here": OP_HERE,
}

const MEM_CAPACITY = 640_000
const STR_CAPACITY = 640_000
