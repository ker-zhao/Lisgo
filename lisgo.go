package main

import (
	"fmt"
	"lisgo/interp"
	"lisgo/parser"
)

func eval(code string) interp.Atom {
	return interp.InterP(parser.Parse(code), interp.GlobalEnv)
}

func evalPrint(code string) string {
	return interp.Stringify(eval(code))
}

func main() {
	//fmt.Println("Hello, 世界")
	td := `
 (quote (1 . 2 . 4 . 5))
`

	parsed := parser.Parse(td)
	parsedS := interp.Stringify(parsed)
	fmt.Println(parsedS)

	r := evalPrint(td)
	fmt.Println("----------------------")
	fmt.Println(r)

}
