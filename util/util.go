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

func WarnWithError(filePath string, row int, err string) {
    errorString := fmt.Sprintf("%s:%d -- %s\n", filePath, row, err)
    fmt.Fprintf(os.Stderr, errorString)
}

func CheckNumberOfArguments(stkCnt int, argCnt int, op model.Operation, opString string) {
    if stkCnt < argCnt {
        errorString := fmt.Sprintf("operation %s requires %d arguments but found %d\n", opString, argCnt, stkCnt)
        TerminateWithError(op.FilePath, op.Row, errorString)
    }
}
