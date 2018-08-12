package repl

import (
	"strconv"
	"regexp"

	"lisgo/interp"
)

var quotes = map[string]*interp.Symbol{
	`'`:  interp.Sym("quote"),
	"`":  interp.Sym("quasiquote"),
	`,`:  interp.Sym("unquote"),
	`,@`: interp.Sym("unquote-splicing"),
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


type reader interface {
	ReadLine() (line []byte, isPrefix bool, err error)
}

type Input struct {
	reader reader
	line []byte
	reg *regexp.Regexp
}

// \s*(,@|[('`,)]|"(?:[\\].|[^\\"])*"|;.*|[^\s('"`,;)]*)(.*)
const tokenizer = `\s*(,@|[('` + "`" + `,)]|"(?:[\\].|[^\\"])*"|;.*|[^\s('"` + "`" + `,;)]*)(.*)`

func newInput(reader reader) *Input {
	return &Input{reader, make([]byte, 0), regexp.MustCompile(tokenizer)}
}

func (s *Input) GetToken() (token string, eof bool) {
	for {
		s.read()
		if len(s.line) == 0 {
			return "", true
		} else {
			groups := s.reg.FindSubmatch(s.line)
			tokenBytes := groups[1]
			s.line = groups[2]
			if len(tokenBytes) > 0 && tokenBytes[0] != ';' {
				return string(tokenBytes), false
			}
		}
	}
}

func (s *Input) read() {
	if len(s.line) == 0 {
		var err error
		s.line, _, err = s.reader.ReadLine()
		checkError(err, "Error: s.reader.ReadLine failed. %s")
	}
}

func (s Input) GetExp() interp.Atom {
	token, eof := s.GetToken()
	if eof {
		panic("Error: GetExp Unexpected End-Of-File.")
	} else {
		return s.parseToken(token)
	}

}

func (s Input) parseToken(token string) interp.Atom {
	if token == "(" {
		l := interp.NewLinkedList()
		for {
			token, eof := s.GetToken()
			if eof {
				panic("Error: parseToken Unexpected End-Of-File.")
			}
			if token == ")" {
				return l.ToPair()
			} else {
				l.Insert(s.parseToken(token))
			}
		}
	} else if token == ")" {
		panic("Error: Unexpected ) in here.")
	} else if v, ok := quotes[token]; ok {
		l := interp.NewLinkedList(interp.NewAtom(interp.TypeSymbol, v), s.GetExp())
		return l.ToPair()
	} else {
		return atom(token)
	}
}

//func Parse(program string) interp.Atom {
//	var index int
//	return interp.Expand(readFromTokens(tokenize(program), &index))
//}
//
//func ParseUnexpand(program string) interp.Atom {
//	var index int
//	return readFromTokens(tokenize(program), &index)
//}
