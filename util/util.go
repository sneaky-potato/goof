package util

import (
    "fmt"
    "os"
)

func TerminateWithError(filePath string, row int, err string) {
    errorString := fmt.Sprintf("%s:%d -- %s\n", filePath, row, err)
    fmt.Fprintf(os.Stderr, errorString)
    os.Exit(1)
}

