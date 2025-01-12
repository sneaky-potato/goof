package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/sneaky-potato/g4th/compiler"
	"github.com/sneaky-potato/g4th/lexer"
)

func callCmd(cmd string, args ...string) {
    command := exec.Command(cmd, args...)
    var outb, errb bytes.Buffer
    command.Stdout = &outb
    command.Stderr = &errb
    fmt.Printf("%s ", cmd)
    fmt.Println(args)

    if err := command.Run(); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s", outb.String())
}

func usage(program string) {
	fmt.Printf("Usage: %s <OPTION> [ARGS]\n", program)
	fmt.Println("OPTIONS:")
	fmt.Println("    sim <file>        Simulate program")
	fmt.Println("    com <file>        Compile program")
	fmt.Println("        SUBOPTIONS:")
	fmt.Println("            -r        run the program after successful compilation")
	fmt.Println("    help              Print this help to stdout")
}

func main() {
    simCmd := flag.NewFlagSet("sim", flag.ExitOnError)
    comCmd := flag.NewFlagSet("com", flag.ExitOnError)
    helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
    runOnCom := comCmd.Bool("r", false, "run")
    
    if len(os.Args) < 2 {
        fmt.Println("expected subcommand")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "sim":
        simCmd.Parse(os.Args[2:])
        filePath := simCmd.Args()[0]
        _ = lexer.LoadProgramFromFile(filePath)
    case "com":
        comCmd.Parse(os.Args[2:])
        filePath := comCmd.Args()[0]
        program := lexer.LoadProgramFromFile(filePath)

        compiler.CompileToAsm("output.asm", program)

        callCmd("nasm", "-felf64", "output.asm")
        callCmd("ld", "-o", "output", "output.o")


        if *runOnCom {
            callCmd("./output")
        }
    default:
        helpCmd.Parse(os.Args[2:])
        usage(os.Args[0])
    }
}
