package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Mdaiki0730/hackasm/code"
	"github.com/Mdaiki0730/hackasm/symtable"
)

type Parser struct {
	file             *os.File
	outfile          *os.File
	scanner          bufio.Scanner
	round            int
	count            int
	symtable         symtable.SymTable
	latestAddressNum int
	Command          string
}

const A_COMMAND = 0
const C_COMMAND = 1
const L_COMMAND = 2

var aChecker = regexp.MustCompile(`@`)
var cChecker = regexp.MustCompile(`=|;`)
var lChecker = regexp.MustCompile(`\(.+?\)`)

func NewParser(f, of *os.File) Parser {
	return Parser{
		file:             f,
		outfile:          of,
		round:            1,
		latestAddressNum: 15,
		symtable:         symtable.NewSymTable(),
		scanner:          *bufio.NewScanner(f),
	}
}

func (p *Parser) HasMoreCommands() bool {
	if !p.scanner.Scan() {
		if p.round >= 2 {
			return false
		} else if p.round == 1 {
			p.scanner = *bufio.NewScanner(p.file)
			p.file.Seek(0, 0)
			p.round += 1
			p.count = 0
			return true
		}
		fmt.Println("unexpected statement")
		os.Exit(1)
	}
	return true
}

func (p *Parser) Advance() {
	// ignore comment out
	line := p.scanner.Text()
	index := strings.Index(line, "//")
	trimmedCommentString := line
	if index != -1 {
		trimmedCommentString = line[:index]
	}
	p.Command = strings.TrimSpace(trimmedCommentString)
	if p.Command == "" || p.Command == "\n" {
		return
	}

	if p.round == 1 {
		if p.CommandType() != L_COMMAND {
			p.count += 1
			return
		}
		symbol := p.symbol()
		p.symtable.AddEntry(symbol, p.count)
	} else {
		switch p.CommandType() {
		case A_COMMAND:
			a := p.symbol()
			value, err := strconv.Atoi(a)
			if err != nil {
				// a = symbol
				if !p.symtable.Contains(a) {
					p.symtable.AddEntry(a, p.latestAddressNum+1)
					p.latestAddressNum += 1
				}
				value = p.symtable.GetAddress(a)
			}

			// validate value
			bin := fmt.Sprintf("%015b", value)
			if len(bin) > 15 {
				fmt.Println("too large value")
				os.Exit(1)
			}
			p.write("0" + bin + "\n")
		case C_COMMAND:
			p.write("111" + code.Comp(p.comp()) + code.Dest(p.dest()) + code.Jump(p.jump()) + "\n")
		case L_COMMAND:
		default:
			fmt.Println("unexpected command")
			os.Exit(1)
		}
	}
}

func (p *Parser) CommandType() int {
	if aChecker.MatchString(p.Command) {
		return A_COMMAND
	} else if cChecker.MatchString(p.Command) {
		return C_COMMAND
	} else if lChecker.MatchString(p.Command) {
		return L_COMMAND
	}
	fmt.Println("unexpected token")
	os.Exit(1)
	return 0
}

func (p *Parser) symbol() string {
	charsToRemove := "@()"
	return strings.Trim(p.Command, charsToRemove)
}

func (p *Parser) dest() string {
	if strings.Contains(p.Command, ";") {
		return "null"
	}
	mnemonics := strings.Split(p.Command, "=")
	return mnemonics[0]
}

func (p *Parser) comp() string {
	if strings.Contains(p.Command, "=") {
		mnemonics := strings.Split(p.Command, "=")
		return mnemonics[1]
	}
	mnemonics := strings.Split(p.Command, ";")
	return mnemonics[0]
}

func (p *Parser) jump() string {
	if strings.Contains(p.Command, "=") {
		return "null"
	}
	mnemonics := strings.Split(p.Command, ";")
	return mnemonics[1]
}

func (p *Parser) write(bin string) error {
	_, err := p.outfile.Write([]byte(bin))
	if err != nil {
		fmt.Println("failed to write file")
		os.Exit(1)
	}
	return nil
}
