package repl

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"lisgo/interp"
)

var quotes = map[string]*interp.Symbol{
	`'`:  interp.Sym(interp.KeyQuote),
	"`":  interp.Sym(interp.KeyQuasiQuote),
	`,`:  interp.Sym(interp.KeyUnquote),
	`,@`: interp.Sym(interp.KeyUnquoteSplicing),
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
	line   []byte
	reg    *regexp.Regexp
}

// \s*(,@|[('`,)]|"(?:[\\].|[^\\"])*"|;.*|[^\s('"`,;)]*)(.*)
const tokenizer = `\s*(,@|[('` + "`" + `,)]|"(?:[\\].|[^\\"])*"|;.*|[^\s('"` + "`" + `,;)]*)(.*)`

func newInput(reader reader) *Input {
	return &Input{reader, make([]byte, 0), regexp.MustCompile(tokenizer)}
}

func (s *Input) GetToken(depth int, prompt string) (token string, eof bool) {
	for {
		if len(s.line) == 0 {
			if depth > 0 && prompt != "" {
				for i := 0; i <= depth; i++ {
					fmt.Print("  ")
				}
			}
			var err error
			s.line, _, err = s.reader.ReadLine()
			if err == io.EOF {
				return "", true
			} else if err != nil {
				fmt.Errorf("error: s.reader.ReadLine failed. %s", err)
			}
		}
		groups := s.reg.FindSubmatch(s.line)
		tokenBytes := groups[1]
		s.line = groups[2]
		if len(tokenBytes) > 0 && tokenBytes[0] != ';' {
			//if _, ok := quotes[string(tokenBytes)]; ok {
			//	s, eof := s.GetToken(prompt)
			//	return string(tokenBytes) + s, eof
			//}
			return string(tokenBytes), false
		}
	}
}

func (s *Input) Parse(prompt string) (exp interp.Atom, eof bool) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	token, eof := s.GetToken(0, prompt)
	if eof {
		return interp.Void, eof
	} else {
		return s.parseToken(token, prompt, 1), false
	}

}

func (s *Input) parseToken(token string, prompt string, depth int) interp.Atom {
	if token == "(" {
		l := interp.NewLinkedList()
		for {
			token, eof := s.GetToken(depth, prompt)
			if eof {
				panic("Error: parseToken Unexpected End-Of-File.")
			}
			if token == ")" {
				return l.ToPair()
			} else {
				l.Insert(s.parseToken(token, prompt, depth+1))
			}
		}
	} else if token == ")" {
		panic("Error: Unexpected ) in here.")
	} else if v, ok := quotes[token]; ok {
		atom, eof := s.Parse("")
		if eof {
			panic("read: expected an element for quoting ' (found end-of-file)")
		}
		l := interp.NewLinkedList(interp.NewAtom(interp.TypeSymbol, v), atom)
		return l.ToPair()
	} else {
		return atom(token)
	}
}

//func Parse(rd io.Reader, prompt string) {
//	reader := bufio.NewReader(rd)
//	input := newInput(reader)
//	for {
//		atom, eof := input.Parse(prompt)
//		if eof {
//			return
//		}
//		val := interp.InterP(interp.Expand(atom), interp.GlobalEnv)
//		if !val.IsType(interp.TVoid) {
//			fmt.Println(interp.Stringify(val))
//		}
//	}
//}

//func Parse(program string) interp.Atom {
//	var index int
//	return interp.Expand(readFromTokens(tokenize(program), &index))
//}
//
//func ParseUnexpand(program string) interp.Atom {
//	var index int
//	return readFromTokens(tokenize(program), &index)
//}
