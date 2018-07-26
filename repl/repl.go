package repl

import (
	"bufio"
	"fmt"
	"lisgo/interp"
	"lisgo/parser"
	"os"
)

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
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		checkError(err, "Read input error: %s \n")
		//fmt.Println(text)
		val := interp.InterP(parser.Parse(text), interp.GlobalEnv)
		if !val.IsType(interp.TVoid) {
			fmt.Println(interp.Stringify(val))
		}
	}
}

func checkError(err error, info string) {
	if err != nil {
		fmt.Errorf(info, err.Error())
	}
}
