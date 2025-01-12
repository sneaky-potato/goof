package main

import (
    "fmt"
    "flag"
    "os"
    "os/exec"

    "github.com/sneaky-potato/g4th/lexer"
)

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
        _ = lexer.LoadProgramFromFile(filePath)

        // compile_program(program, "output.asm")

        exec.Command("nasm", "-felf64", "output.asm")
        exec.Command("ld", "-o", "output")

        if *runOnCom {
            exec.Command("./output")
        }
    default:
        helpCmd.Parse(os.Args[2:])
        usage(os.Args[0])
    }
}
