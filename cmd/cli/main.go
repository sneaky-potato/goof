package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
    "strings"

	"github.com/sneaky-potato/goof/compiler"
	"github.com/sneaky-potato/goof/lexer"
	"github.com/sneaky-potato/goof/types"
)

func callCmd(cmd string, args ...string) {
    command := exec.Command(cmd, args...)
    var outb, errb bytes.Buffer
    command.Stdout = &outb
    command.Stderr = &errb
    fmt.Printf("%s ", cmd)
    fmt.Println(args)

    if err := command.Run(); err != nil {
        fmt.Printf("%s", errb.String())
        log.Fatal(err)
    }

    fmt.Printf("%s", outb.String())
}

func usage(program string) {
	fmt.Printf("Usage: %s [ARGS] <file>\n", program)
	fmt.Println("ARGUMENTS:")
	fmt.Println("    -r         run the program after successful compilation")
	fmt.Println("    -unsafe    run in UNSAFE mode - skip static type checking before compilation")
}

func main() {
    runOnCom := flag.Bool("r", false, "run")
    skipTypeChecking := flag.Bool("unsafe", false, "skip static type checking [UNSAFE]")
    
    if len(os.Args) < 2 {
        fmt.Println("expected <file>")
        usage(os.Args[0])
        os.Exit(1)
    }

    flag.Parse()

    if len(flag.Args()) < 1 {
        fmt.Println("expected <file>")
        usage(os.Args[0])
        os.Exit(1)
    }

    filePath := flag.Args()[0]
    program := lexer.LoadProgramFromFile(filePath)
    fileName := filepath.Base(filePath)
    fileNameWithoutExtension := strings.TrimSuffix(fileName, filepath.Ext(fileName))
    if !(*skipTypeChecking) {
        types.TypeCheckingProgram(program)
    }

    compiler.CompileToAsm(fileNameWithoutExtension + ".asm", program)

    callCmd("nasm", "-felf64", fileNameWithoutExtension + ".asm")
    callCmd("ld", "-o", fileNameWithoutExtension, fileNameWithoutExtension + ".o")

    if *runOnCom {
        callCmd("./" + fileNameWithoutExtension, flag.Args()[1:]...)
    }
}
