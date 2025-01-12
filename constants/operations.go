package constants

const (
    OP_PUSH = iota
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

const MEM_CAPACITY = 640_000
