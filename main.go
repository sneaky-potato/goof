package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

    "github.com/sneaky-potato/goof/lexer"
    "github.com/sneaky-potato/goof/types"
	"github.com/sneaky-potato/goof/compiler"
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

    filePath := flag.Args()[0]
    program := lexer.LoadProgramFromFile(filePath)

    if *skipTypeChecking == false {
        types.TypeCheckingProgram(program)
    }

    compiler.CompileToAsm("output.asm", program)

    callCmd("nasm", "-felf64", "output.asm")
    callCmd("ld", "-o", "output", "output.o")

    if *runOnCom {
        callCmd("./output", flag.Args()[1:]...)
    }
}
