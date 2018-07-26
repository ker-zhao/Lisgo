package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"lisgo/interp"
)

var quotes = map[string]*interp.Symbol{
	`'`:  interp.Sym("quote"),
	"`":  interp.Sym("quasiquote"),
	`,`:  interp.Sym("unquote"),
	`,@`: interp.Sym("unquote-splicing"),
}

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
			word = fmt.Sprintf("%s%c", word, char)
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
	keyChars := []string{"(", ")", "'", "`", ",@"}
	for _, v := range keyChars {
		chars = strings.Replace(chars, v, " "+v+" ", -1)
	}
	chars = strings.Replace(chars, "[", " ( ", -1)
	chars = strings.Replace(chars, "]", " ) ", -1)
	unquote, _ := regexp.Compile(",([^@])")
	chars = unquote.ReplaceAllString(chars, " , $1")
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
	} else if v, ok := quotes[token]; ok {
		l := interp.NewLinkedList(interp.NewAtom(interp.TypeSymbol, v), readFromTokens(tokens, index))
		return l.ToPair()
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
	return interp.Expand(readFromTokens(tokenize(program), &index))
}

func ParseUnexpand(program string) interp.Atom {
	var index int
	return readFromTokens(tokenize(program), &index)
}

func checkError(err error, info string) {
	if err != nil {
		fmt.Errorf(info, err.Error())
	}
}
