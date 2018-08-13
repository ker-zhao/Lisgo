package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"lisgo/interp"
)

const Prompt = "> "

func REPL() {
	for {
		repl()
	}
}

func repl() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	Exec(os.Stdin, Prompt)
}

func Exec(rd io.Reader, prompt string) {
	reader := bufio.NewReader(rd)
	input := newInput(reader)
	for {
		atom, eof := input.Parse(prompt)
		if eof {
			return
		}
		val := interp.InterP(interp.Expand(atom), interp.GlobalEnv)
		if !val.IsType(interp.TVoid) {
			fmt.Println(interp.Stringify(val))
		}
	}
}

func checkError(err error, info string) {
	if err != nil {
		fmt.Printf(info, err.Error()+"\n")
	}
}

func checkErrorPanic(err error, info string) {
	if err != nil {
		panic(fmt.Sprintf(info, err.Error()))
	}
}
