package main

import (
	"lisgo/repl"
	"os"

)

//func eval(code string) interp.Atom {
//	return interp.InterP(parser.Parse(code), interp.GlobalEnv)
//}

//func evalPrint(code string) string {
//	return interp.Stringify(eval(code))
//}

func main() {
	//测试代码 `(,@x 3 4 ,`x)可以测试多重的quasiquote

	//td := "(begin (define x '(1 2)) (define y '(3))  `(,@x ,y 4))"
	//
	//parsed := parser.ParseUnexpand(td)
	//
	//parsedRaw := interp.Stringify(parsed)
	//fmt.Println(parsedRaw)
	//
	//parsedS := interp.Stringify(interp.Expand(parsed))
	//fmt.Println(parsedS)
	//
	//r := evalPrint(td)
	//fmt.Println("----------------------")
	//fmt.Println(r)
	//


	if len(os.Args) > 1 {
		fileName := os.Args[1]
		repl.ExecFile(fileName)
	} else {
		repl.REPL()
	}

}
