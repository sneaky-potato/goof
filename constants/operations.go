package constants

const (
    OP_PUSH_INT = iota
    OP_PUSH_STR
    OP_PUSH_PTR
    OP_PLUS
    OP_MINUS
    OP_MUL
    OP_MOD
    OP_EQ
    OP_NE
    OP_GT
    OP_LT
    OP_IF
    OP_ELIF
    OP_ELSE
    OP_END
    OP_WHILE
    OP_DO
    OP_MACRO
    OP_MEMORY
    OP_CONST
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
    OP_SKIP_PROC
    OP_PREP_PROC
    OP_RET
    OP_CALL
    OP_PROC_SEP
    OP_CAST_INT
    OP_CAST_PTR
    OP_CAST_BOOL
    OP_TYPE_INT
    OP_TYPE_PTR
    OP_TYPE_BOOL
    OP_SYSCALL0
    OP_SYSCALL1
    OP_SYSCALL2
    OP_SYSCALL3
    OP_SYSCALL4
    OP_SYSCALL6
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
    "elif": OP_ELIF,
    "else": OP_ELSE,
    "end": OP_END,
    "while": OP_WHILE,
    "do": OP_DO,
    "macro": OP_MACRO,
    "memory": OP_MEMORY,
    "const": OP_CONST,
    "proc": OP_SKIP_PROC,
    "ret": OP_RET,
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
    "int": OP_TYPE_INT,
    "ptr": OP_TYPE_PTR,
    "bool": OP_TYPE_BOOL,
    "--": OP_PROC_SEP,
    "syscall0": OP_SYSCALL0,
    "syscall1": OP_SYSCALL1,
    "syscall2": OP_SYSCALL2,
    "syscall3": OP_SYSCALL3,
	"syscall4": OP_SYSCALL4,
    "syscall6": OP_SYSCALL6,
    "argv": OP_ARGV,
    "argc": OP_ARGC,
    "dump": OP_DUMP,
    "include": OP_INCLUDE,
    "here": OP_HERE,
}
