package constants

const (
    OP_PUSH_INT = iota
    OP_PUSH_STR
    OP_DUMP
    OP_PLUS
    OP_MINUS
    OP_EQUAL
    OP_IF
    OP_ELSE
    OP_END
    OP_DUP
    OP_2DUP
    OP_DROP
    OP_OVER
    OP_SHR
    OP_SHL
    OP_OR
    OP_AND
    OP_SWAP
    OP_ROT
    OP_GT
    OP_LT
    OP_WHILE
    OP_DO
    OP_MEM
    OP_LOAD
    OP_STORE
    OP_SYSCALL3
    COUNT_OPS
)

var BUILTIN_WORDS = map[string]int{
    "+": OP_PLUS,
    "-": OP_MINUS,
    "=": OP_EQUAL,
    "if": OP_IF,
    "else": OP_ELSE,
    "end": OP_END,
    "dump": OP_DUMP,
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
    "while": OP_WHILE,
    "do": OP_DO,
    "mem": OP_MEM,
    ",": OP_LOAD,
    ".": OP_STORE,
    "syscall3": OP_SYSCALL3,
}

const MEM_CAPACITY = 640_000
const STR_CAPACITY = 640_000
