package repl

import (
	"bufio"
	"fmt"
	"os"

	"lisgo/interp"
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
		input := newInput(reader)
		//text, err := reader.ReadString('\n')
		//bytes, _, err := reader.ReadLine()
		//text := string(bytes)
		//checkError(err, "Read input error: %s \n")
		val := interp.InterP(input.GetExp(), interp.GlobalEnv)
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
