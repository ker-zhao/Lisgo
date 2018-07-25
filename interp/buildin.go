package interp

type BuildInProcedure struct {
	Opt *Symbol
	_   *Symbol
	handleFunc func(...Atom) Atom
}

func NewBuildInProcedure(Opt *Symbol, handle func(...Atom) Atom) BuildInProcedure {
	return BuildInProcedure{Opt, Opt, handle}
}

func NewBuildInAtom(Opt *Symbol, handle func(...Atom) Atom) Atom {
	return NewAtom(TypeBuildIn, NewBuildInProcedure(Opt, handle))
}

func (s BuildInProcedure) call(args ...Atom) Atom {
	return s.handleFunc(args...)
}

func Cons(args ...Atom) Atom {
	carPart, cdrPart := args[0], args[1]
	return NewPairAtom(carPart, cdrPart)
}

func Car(args ...Atom) Atom {
	pair := *(*Pair)(args[0].Data)
	return pair.Car
}

func Cdr(args ...Atom) Atom {
	pair := *(*Pair)(args[0].Data)
	return pair.Cdr
}

func Append(args ...Atom) Atom {
	l := NewLinkedList()
	for _, v := range args[:len(args) - 1] {
		if v.IsType(TPair) {
			Foreach(v, func(i int, atom Atom) {
				l.Insert(atom)
			})
		} else {
			//fmt.Errorf("expected: list? given: %s\n", parser.Stringify(v))
		}
	}
	pair := (*Pair)(l.Last.Data)
	(*pair).Cdr = args[len(args) - 1]
	return l.ToPair()
}

func List(args ...Atom) Atom {
	l := NewLinkedList(args...)
	return l.ToPair()
}

func Eq(args ...Atom) Atom {
	x, y := args[0], args[1]
	if x.IsType(TPair) {
		return NewAtom(TypeBoolean, x.Data == y.Data)
	} else {
		return Equal(args...)
	}

}

func Equal(args ...Atom) Atom {
	return NewAtom(TypeBoolean, AtomEqual(args[0], args[1]))
}

func Length(args ...Atom) Atom {
	return NewAtom(TypeInt, len(PairToSlice(args[0])))
}


func BasicOptMaker(opt *Symbol) func(...Atom) Atom{
	return func(atom ...Atom) Atom {
		return BasicOpt(opt, atom...)
	}
}

func BasicOpt(opt *Symbol, args ...Atom) Atom {
	for {
		var result Atom
		switch opt {
		case Sym("+"), Sym("-"), Sym("*"), Sym("/"):
			x, y := args[0], args[1]
			if x.IsType(TInt) {
				var r Int
				switch opt {
				case Sym("+"):
					r = *(*Int)(x.Data) + *(*Int)(y.Data)
				case Sym("-"):
					r = *(*Int)(x.Data) - *(*Int)(y.Data)
				case Sym("*"):
					r = *(*Int)(x.Data) * *(*Int)(y.Data)
				case Sym("/"):
					r = *(*Int)(x.Data) / *(*Int)(y.Data)
				}
				result = NewAtom(TypeInt, r)
			} else {
				var r Float
				switch opt {
				case Sym("+"):
					r = *(*Float)(x.Data) + *(*Float)(y.Data)
				case Sym("-"):
					r = *(*Float)(x.Data) - *(*Float)(y.Data)
				case Sym("*"):
					r = *(*Float)(x.Data) * *(*Float)(y.Data)
				case Sym("/"):
					r = *(*Float)(x.Data) / *(*Float)(y.Data)
				}
				result = NewAtom(TypeFloat, r)
			}
		case Sym(">"), Sym("<"), Sym(">="), Sym("<="), Sym("="):
			x, y := args[0], args[1]
			var r Boolean
			if x.IsType(TInt) {
				switch opt {
				case Sym(">"):
					r = *(*Int)(x.Data) > *(*Int)(y.Data)
				case Sym("<"):
					r = *(*Int)(x.Data) < *(*Int)(y.Data)
				case Sym(">="):
					r = *(*Int)(x.Data) >= *(*Int)(y.Data)
				case Sym("<="):
					r = *(*Int)(x.Data) <= *(*Int)(y.Data)
				case Sym("="):
					r = *(*Int)(x.Data) == *(*Int)(y.Data)
				}
				result = NewAtom(TypeBoolean, r)
			} else {
				switch opt {
				case Sym(">"):
					r = *(*Float)(x.Data) > *(*Float)(y.Data)
				case Sym("<"):
					r = *(*Float)(x.Data) < *(*Float)(y.Data)
				case Sym(">="):
					r = *(*Float)(x.Data) >= *(*Float)(y.Data)
				case Sym("<="):
					r = *(*Float)(x.Data) <= *(*Float)(y.Data)
				case Sym("="):
					r = *(*Float)(x.Data) == *(*Float)(y.Data)
				}
				result = NewAtom(TypeBoolean, r)
			}
		case Sym("and"), Sym("or"), Sym("xor"):
			x, y := args[0], args[1]
			var r Boolean
			switch opt {
			case Sym("and"):
				r = *(*Boolean)(x.Data) && *(*Boolean)(y.Data)
			case Sym("or"):
				r = *(*Boolean)(x.Data) || *(*Boolean)(y.Data)
			case Sym("xor"):
				r = *(*Boolean)(x.Data) != *(*Boolean)(y.Data)
			}
			result = NewAtom(TypeBoolean, r)
		case Sym("not"):
			r := !(*(*Boolean)(args[0].Data))
			result = NewAtom(TypeBoolean, r)
		}
		
		if len(args) <= 2 {
			return result
		} else {
			newArgs := []Atom{result}
			args = append(newArgs, args[2:]...)
		}
	}
}




func StandardEnv(env *Env) *Env {
	opts := []string{"+", "-", "*", "/", ">", "<", ">=", "<=", "="}
	for _, v := range opts {
		env.extBuildIn(v, BasicOptMaker(Sym(v)))
	}
	env.extBuildIn("cons", Cons)
	env.extBuildIn("car", Car)
	env.extBuildIn("cdr", Cdr)
	env.extBuildIn("append", Append)
	env.extBuildIn("list", List)
	env.extBuildIn("eq?", Eq)
	env.extBuildIn("equal?", Equal)
	env.extBuildIn("length", Length)

	return env
}

