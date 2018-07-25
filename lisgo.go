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
(list "\"123" "1234" "string" 123 (list 1 2 3 "4"))
`

	parsed := parser.Parse(td)
	parsedS := interp.Stringify(parsed)
	fmt.Println(parsedS)

	r := evalPrint(td)
	fmt.Println("----------------------")
	fmt.Println(r)

}
