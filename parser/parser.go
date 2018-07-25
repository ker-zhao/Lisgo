package parser

import (
	"fmt"
	"strconv"
	"strings"

	"lisgo/interp"
)

func splitAtom(str string) []string {
	l := make([]string, 0)
	inStr := false
	pre := ' '
	word := ""
	for _, char := range str {
		if (char == ' ' || char == '\n' || char == '\t') && !inStr {
			if word != "" {
				l = append(l, word)
				word = ""
			}
		} else if char == '"' && pre != '\\' {
			inStr = !inStr
		} else {
			word = fmt.Sprintf("%s%c", word, char)
		}
		pre = char
	}
	if word != "" {
		l = append(l, word)
	}
	return l
}

func tokenize(chars string) []string {
	chars = strings.Replace(chars, "(", " ( ", -1)
	chars = strings.Replace(chars, ")", " ) ", -1)
	chars = strings.Replace(chars, "[", " ( ", -1)
	chars = strings.Replace(chars, "]", " ) ", -1)
	return splitAtom(chars)
}

func readFromTokens(tokens []string, index *int) interp.Atom {
	token := tokens[*index]
	*index += 1

	if token == "(" {
		l := interp.NewLinkedList()
		for tokens[*index] != ")" {
			l.Insert(readFromTokens(tokens, index))
		}
		*index += 1
		return *l.First
	} else {
		return atom(token)
	}
}

func atom(str string) interp.Atom {
	if str == "#t" {
		return interp.NewAtom(interp.TypeBoolean, interp.Boolean(true))
	} else if str == "#f" {
		return interp.NewAtom(interp.TypeBoolean, interp.Boolean(false))
	} else if str[0] == '"' {
		s, err := strconv.Unquote(str)
		checkError(err, "Unquote string atom filed, detail: %s.")
		return interp.NewAtom(interp.TypeString, interp.String(s))
	} else if Int, err2 := strconv.Atoi(str); err2 == nil {
		return interp.NewAtom(interp.TypeInt, interp.Int(Int))
	} else if Float, err3 := strconv.ParseFloat(str, 64); err3 == nil {
		return interp.NewAtom(interp.TypeFloat, interp.Float(Float))
	} else {
		return interp.NewAtom(interp.TypeSymbol, interp.Sym(str))
	}
}

func Parse(program string) interp.Atom {
	var index int
	return readFromTokens(tokenize(program), &index)
}



func checkError(err error, info string) {
	if err != nil {
		fmt.Errorf(info, err.Error())
	}
}
