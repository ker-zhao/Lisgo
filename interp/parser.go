package interp

import (
	"fmt"
	"strconv"
	"strings"
)

var quotesReflect = map[*Symbol]string{
	Sym(KeyQuote):           `'`,
	Sym(KeyQuasiQuote):      "`",
	Sym(KeyUnquote):         `,`,
	Sym(KeyUnquoteSplicing): `,@`,
}

func Stringify(atom Atom) string {
	return StringifyInner(atom, false)
}

func StringifyInner(atom Atom, inQuote bool) string {
	if IsList(atom) && ListLength(atom) >= 2 && PairGet(atom, 0).IsType(TSymbol) {
		symbol := (*Symbol)(PairGet(atom, 0).Data)
		if v, ok := quotesReflect[symbol]; ok {
			return v + StringifyInner(PairGet(atom, 1), inQuote)
		}
	}
	if atom.IsType(TBoolean) {
		if *(*Boolean)(atom.Data) {
			return "#t"
		} else {
			return "#f"
		}
	} else if atom.IsType(TSymbol) {
		r := string(*(*Symbol)(atom.Data))
		if !inQuote {
			r = "'" + r
		}
		return r
	} else if atom.IsType(TString) {
		r := strings.Replace(string(*(*String)(atom.Data)), `"`, `\"`, -1)
		return `"` + r + `"`
	} else if atom.IsType(TPair) {
		r := "("
		if !inQuote {
			inQuote = true
			r = "'" + r
		}
		i := 0
		for pair := (*Pair)(atom.Data); pair != nil; i, pair = i+1, (*Pair)(pair.Cdr.Data) {
			r += StringifyInner(pair.Car, inQuote) + " "
			if !pair.Cdr.IsType(TPair) {
				i += 1
				r += "." + " " + StringifyInner(pair.Cdr, inQuote) + " "
				break
			}
		}
		if i > 0 {
			r = r[:len(r)-1]
		}
		return r + ")"
	} else if atom.IsType(TInt) {
		return strconv.Itoa(int(*(*Int)(atom.Data)))
	} else if atom.IsType(TFloat) { // Float
		f := float64(*(*Float)(atom.Data))
		r := strconv.FormatFloat(f, 'E', -1, 64)
		return r
	} else if atom.IsType(TClosure) || atom.IsType(TBuildIn) {
		tmp := string(*atom.ObjInfo.Name)
		return fmt.Sprintf("#<procedure:%s>", tmp)
	} else { // Void, maybe never should be see this?
		return "#<void>"
	}
}
