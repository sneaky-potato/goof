package util

import (
	"fmt"
	"os"

    "github.com/sneaky-potato/g4th/model"
)

func TerminateWithError(filePath string, row int, err string) {
    errorString := fmt.Sprintf("%s:%d -- %s\n", filePath, row, err)
    fmt.Fprintf(os.Stderr, errorString)
    os.Exit(1)
}

func checkNumberOfArguments(stkCnt int, argCnt int, op model.Operation) {
    if stkCnt < argCnt {
        TerminateWithError(op.FilePath, op.Row, "1 argument is required for dup")
    }
}
