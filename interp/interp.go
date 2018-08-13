package interp

func require(exp Atom, cond bool, msg string) {
	if !cond {
		panic(Stringify(exp) + ": " + msg)
	}
}

func Expand(x Atom) Atom {
	sym := (*Symbol)(PairGet(x, 0).Data)
	//exps := (*(*Pair)(x.Data)).Cdr

	require(x, !IsPair(x) || ListLength(x) != 0, MsgWrongLength)
	if !IsPair(x) {
		return x
	} else if sym == Sym(KeyQuote) {
		require(x, ListLength(x) == 2, MsgWrongLength)
		return x
	} else if sym == Sym(KeyIf) {
		require(x, ListLength(x) == 4, MsgWrongLength)
		return Map(Expand, x)
	} else if sym == Sym(KeySet) {
		require(x, ListLength(x) == 3, MsgWrongLength)
		v := PairGet(x, 1)
		require(x, v.IsType(TSymbol), "set! argument must be a symbol")
		return List(NewSymbol(KeySet), v, Expand(PairGet(x, 2)))
	} else if sym == Sym(KeyDefine) {
		require(x, ListLength(x) >= 3, MsgWrongLength)
		def, v, body := PairGet(x, 0), PairGet(x, 1), Cdr(Cdr(x))
		if IsList(v) && ListLength(v) != 0 {
			funcName, args := PairGet(v, 0), Cdr(v)
			return Expand(List(def, funcName, Append(List(NewSymbol(KeyLambda), args), body)))
		} else {
			require(x, ListLength(x) == 3, MsgWrongLength)
			require(x, v.IsType(TSymbol), "define argument must be a symbol")
			exp := PairGet(x, 2)
			return List(def, v, exp)
		}
	} else if sym == Sym(KeyBegin) {
		if ListLength(x) == 1{
			return Void
		} else {
			return Map(Expand, x)
		}
	} else if sym == Sym(KeyLambda) {
		require(x, ListLength(x) >= 3, MsgWrongLength)
		vars, body := PairGet(x, 1), Cdr(Cdr(x))
		if IsList(vars) {
			allSymbol := true
			Foreach(vars, func(_ int, atom Atom) {
				allSymbol = allSymbol && atom.IsType(TSymbol)
			})
			require(x, allSymbol, "lambda arguments list must be symbols")
		} else {
			require(x, vars.IsType(TSymbol), "lambda argument must be a symbol")
		}
		var exp Atom
		if ListLength(body) == 1 {
			exp = PairGet(body, 0)
		} else {
			exp = Cons(NewSymbol(KeyBegin), body)
		}
		return List(NewSymbol(KeyLambda), vars, Expand(exp))
	} else if (*Symbol)(PairGet(x, 0).Data) == Sym("quasiquote") {
		return expandQuote(PairGet(x, 1))
	} else {
		return Map(Expand, x)
	}
}

func expandQuote(exp Atom) Atom {
	if !IsPair(exp) {
		l := NewLinkedList(NewAtom(TypeSymbol, Sym("quote")), exp)
		return l.ToPair()
	} else {
		sym := PairGet(exp, 0)
		if IsPair(sym) && (*Symbol)(PairGet(sym, 0).Data) == Sym("unquote-splicing") {
			require(sym, ListLength(sym) == 2, "wrong length")
			arg := PairGet(PairGet(exp, 0), 1)
			l := NewLinkedList(NewAtom(TypeSymbol, Sym("append")), arg, expandQuote(Cdr(exp)))
			return l.ToPair()
		} else if (*Symbol)(sym.Data) == Sym("unquote") {
			require(exp, ListLength(exp) == 2, "wrong length")
			return PairGet(exp, 1)
		} else {
			l := NewLinkedList(NewAtom(TypeSymbol, Sym("cons")), expandQuote(sym), expandQuote(Cdr(exp)))
			return l.ToPair()
		}
	}
}

func InterP(exp Atom, env *Env) Atom {
	if exp.IsType(TSymbol) {
		x := (*Symbol)(exp.Data)
		v := env.find(x)[x]
		return v
	} else if !exp.IsType(TPair) { // int float bool string void
		return exp
	} else {
		sym := (*Symbol)((*(*Pair)(exp.Data)).Car.Data)
		exps := (*(*Pair)(exp.Data)).Cdr
		if sym == Sym("quote") {
			return PairGet(exps, 0)
		} else if sym == Sym("if") {
			test, conseq, alt := PairGet(exps, 0), PairGet(exps, 1), PairGet(exps, 2)
			if *(*Boolean)(InterP(test, env).Data) {
				return InterP(conseq, env)
			} else {
				return InterP(alt, env)
			}
		} else if sym == Sym("set!") {
			ref, args := (*Symbol)(PairGet(exps, 0).Data), PairGet(exps, 1)
			env.find(ref)[ref] = InterP(args, env)
			return Void
		} else if sym == Sym("define") {
			ref, args := (*Symbol)(PairGet(exps, 0).Data), PairGet(exps, 1)
			env.ext(ref, InterP(args, env))
			return Void
		} else if sym == Sym("lambda") {
			params, args := PairGet(exps, 0), PairGet(exps, 1)
			return NewAtom(TypeClosure, NewClosure(params, args, env))
		} else if sym == Sym("let") {
			binds, body := PairToSlice(PairGet(exps, 0)), PairGet(exps, 1)
			params := NewLinkedList()
			args := NewLinkedList()
			for _, v := range binds {
				params.Insert(PairGet(v, 0))
				args.Insert(InterP(PairGet(v, 1), env))
			}
			return InterP(body, NewEnv(params.ToPair(), args.ToPair(), env))
		} else if sym == Sym("begin") {
			list := PairToSlice(exps)
			for _, v := range list[:len(list)-1] {
				InterP(v, env)
			}
			return InterP(list[len(list)-1], env)
		} else {
			list := PairToSlice(exp)
			values := make([]Atom, 0)
			for _, v := range list {
				values = append(values, InterP(v, env))
			}
			if values[0].IsType(TClosure) {
				function := (*Closure)(values[0].Data)
				return function.call(values[1:]...)
			} else if values[0].IsType(TBuildIn) {
				function := (*BuildInProcedure)(values[0].Data)
				return function.call(values[1:]...)
			} else {
				panic("application: not a procedure. >>> " + Stringify(PairGet(exp, 0)))
			}
		}
	}
}
