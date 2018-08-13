package interp

import "fmt"
import "unsafe"

type ObjTypeKind int

const (
	TSymbol ObjTypeKind = iota
	TBoolean
	TInt
	TFloat
	TString
	TClosure
	TBuildIn
	TPair
	TVoid
)

type Symbol string
type Boolean bool
type Int int
type Float float64
type String string

func Sym(s string) *Symbol {
	if v, ok := symbolTable[s]; ok {
		return v
	} else {
		sym := new(Symbol)
		*sym = Symbol(s)
		symbolTable[s] = sym
		return sym
	}
}

func NewSymbol(s string) Atom {
	return NewAtom(TypeSymbol, Sym(s))
}

type ObjType struct {
	ObjTypeKind
}

func NewObjType(t ObjTypeKind) *ObjType {
	return &ObjType{t}
}

type EmptyInterface struct {
	Type unsafe.Pointer
	Word unsafe.Pointer
}

type Atom struct {
	ObjType *ObjType
	Data    unsafe.Pointer
	I       interface{}
}

func NewAtom(objT *ObjType, data interface{}) Atom {
	empty := *((*EmptyInterface)(unsafe.Pointer(&data)))
	return Atom{
		objT,
		empty.Word,
		data,
	}
}

func NewAtomPtr(objT *ObjType, data interface{}) *Atom {
	empty := *((*EmptyInterface)(unsafe.Pointer(&data)))
	return &Atom{
		objT,
		empty.Word,
		data,
	}
}

func (s Atom) IsType(t ObjTypeKind) bool {
	return s.ObjType.ObjTypeKind == t
}

func IsPair(atom Atom) bool {
	return atom.IsType(TPair) && (atom.Data != nil)
}

func AtomEqual(x Atom, y Atom) bool {
	if x.IsType(y.ObjType.ObjTypeKind) {
		if x.IsType(TPair) {
			if IsList(x) && IsList(y) {
				if ListLength(x) != ListLength(y) {
					return false
				}
				for i, pair := 0, (*Pair)(x.Data); pair != nil; i, pair = i+1, (*Pair)(pair.Cdr.Data) {
					if !AtomEqual(pair.Car, PairGet(y, i)) {
						return false
					}
				}
				return true
			} else if IsList(x) == IsList(y) {
				for i, pair, pairY := 0, (*Pair)(x.Data), (*Pair)(y.Data); pair != nil; i, pair, pairY = i+1, (*Pair)(pair.Cdr.Data), (*Pair)(pairY.Cdr.Data) {
					if !AtomEqual(pair.Car, pairY.Car) || !AtomEqual(pair.Cdr, pairY.Cdr) {
						return false
					}
					if !pair.Cdr.IsType(TPair) {
						return true
					}
				}
				return true
			} else {
				return false
			}
		} else if x.IsType(TBoolean) {
			return *(*Boolean)(x.Data) == *(*Boolean)(y.Data)
		} else if x.IsType(TInt) {
			return *(*Int)(x.Data) == *(*Int)(y.Data)
		} else if x.IsType(TFloat) {
			return *(*Float)(x.Data) == *(*Float)(y.Data)
		} else if x.IsType(TString) {
			return *(*String)(x.Data) == *(*String)(y.Data)
		} else if x.IsType(TSymbol) {
			return *(*Symbol)(x.Data) == *(*Symbol)(y.Data)
		} else if x.IsType(TClosure) {
			return *(*Closure)(x.Data) == *(*Closure)(y.Data)
		} else if x.IsType(TBuildIn) {
			return *(*BuildInProcedure)(x.Data).Opt == *(*BuildInProcedure)(y.Data).Opt
		} else {
			fmt.Errorf("AtomEqual, input error, unknown type. \n ")
			fmt.Println(x, y)
			return false
		}
	} else {
		return false
	}
}

type Pair struct {
	Car Atom
	Cdr Atom
}

func NewPairAtom(carPart Atom, cdrPart Atom) Atom {
	return NewAtom(TypePair, Pair{
		carPart,
		cdrPart,
	})
}

func NewPairAtomPtr(carPart Atom, cdrPart Atom) *Atom {
	return NewAtomPtr(TypePair, Pair{
		carPart,
		cdrPart,
	})
}

func Foreach(p Atom, f func(int, Atom)) {
	for i, pair := 0, (*Pair)(p.Data); pair != nil; i, pair = i+1, (*Pair)(pair.Cdr.Data) {
		f(i, pair.Car)
	}
}

func Map(f func(x Atom) Atom, p Atom) Atom {
	l := NewLinkedList()
	for i, pair := 0, (*Pair)(p.Data); pair != nil; i, pair = i+1, (*Pair)(pair.Cdr.Data) {
		l.Insert(f(pair.Car))
	}
	return l.ToPair()
}

func PairToSlice(p Atom) []Atom {
	l := make([]Atom, 0)
	for pair := (*Pair)(p.Data); pair != nil; pair = (*Pair)(pair.Cdr.Data) {
		l = append(l, pair.Car)
	}
	return l
}

func ListLength(p Atom) int {
	i := 0
	for pair := (*Pair)(p.Data); pair != nil; pair = (*Pair)(pair.Cdr.Data) {
		i += 1
	}
	return i
}

func PairGet(p Atom, n int) Atom {
	pair := (*Pair)(p.Data)
	for i := 0; i < n; pair, i = (*Pair)(pair.Cdr.Data), i+1 {
	}
	return pair.Car
}

type LinkedList struct {
	First *Atom // Pair
	Last  *Atom // Pair
}

func NewLinkedList(xs ...Atom) *LinkedList {
	e := EmptyPair
	l := LinkedList{&e, nil}
	for _, v := range xs {
		l.Insert(v)
	}
	return &l
}

func (s *LinkedList) Insert(x Atom) {
	p := NewPairAtomPtr(x, EmptyPair)
	if s.Last == nil {
		s.First, s.Last = p, p
	} else {
		lastPair := (*Pair)(s.Last.Data)
		lastPair.Cdr = *p
		s.Last = p
		if s.First.Data == nil {
			s.First = p
		}
	}
}

func (s *LinkedList) ToPair() Atom {
	return *s.First
}

func IsList(atom Atom) bool {
	if !atom.IsType(TPair) {
		return false
	}
	for pair := (*Pair)(atom.Data); pair != nil; pair = (*Pair)(pair.Cdr.Data) {
		if !pair.Cdr.IsType(TPair) {
			return false
		}
	}
	return true
}

type Env struct {
	vars  map[*Symbol]Atom
	outer *Env
}

func NewEnv(params Atom, args Atom, outer *Env) *Env {
	e := Env{
		make(map[*Symbol]Atom),
		outer,
	}
	require(params, ListLength(params) == ListLength(args), MsgWrongLength+
		fmt.Sprintf(" expect %d, giving %d, got: %s",
			ListLength(params), ListLength(args), Stringify(args)))
	e.zipUpdate(params, args)
	return &e
}

func (s *Env) ext(x *Symbol, v Atom) *Env {
	s.vars[x] = v
	return s
}

func (s *Env) extBuildIn(x string, handler func(...Atom) Atom) *Env {
	s.ext(Sym(x), NewBuildInAtom(Sym(x), handler))
	return s
}

func (s *Env) zipUpdate(params Atom, args Atom) *Env {
	p, a := PairToSlice(params), PairToSlice(args)
	for i, v := range p {
		s.vars[(*Symbol)(v.Data)] = a[i]
	}
	return s
}

func (s *Env) update(dict map[*Symbol]Atom) *Env {

	for i, v := range dict {
		s.vars[i] = v
	}
	return s
}

func (s *Env) find(x *Symbol) map[*Symbol]Atom {
	if _, ok := s.vars[x]; ok {
		return s.vars
	} else {
		if s.outer != nil {
			return s.outer.find(x)
		} else {
			fmt.Printf("ERROR: Lookup error: %s \n", *x)
			return nil
		}
	}
}

type Closure struct {
	params Atom
	body   Atom
	env    *Env
}

func NewClosure(params Atom, exp Atom, env *Env) Closure {
	return Closure{params, exp, env}
}

func (s Closure) call(args ...Atom) Atom {
	l := NewLinkedList(args...)
	return InterP(s.body, NewEnv(s.params, *l.First, s.env))
}
