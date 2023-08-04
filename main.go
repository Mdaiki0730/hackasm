package main

import (
	"fmt"
	"os"

	"github.com/Mdaiki0730/hackasm/parser"
)

func main() {
	// arg validation
	if len(os.Args) != 2 {
		fmt.Println("please enter only target file name")
		os.Exit(1)
	}

	// file setting
	f, err := os.Open(os.Args[1])
	defer f.Close()
	if err != nil {
		fmt.Println("no such file")
		os.Exit(1)
	}

	of, err := os.Create("out.hack")
	defer of.Close()
	if err != nil {
		fmt.Println("failed to file create")
		os.Exit(1)
	}

	// main process
	p := parser.NewParser(f, of)
	for p.HasMoreCommands() {
		p.Advance()
	}
}
