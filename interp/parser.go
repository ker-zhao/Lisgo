package interp

import (
	"strconv"
	"strings"
	"fmt"
)

func Stringify(atom Atom) string {
	if atom.IsType(TBoolean) {
		if *(*Boolean)(atom.Data) {
			return "#t"
		} else {
			return "#f"
		}
	} else if atom.IsType(TSymbol) {
		return string(*(*Symbol)(atom.Data))
	} else if atom.IsType(TString) {
		r := strings.Replace(string(*(*String)(atom.Data)), `"`, `\"`, -1)
		return `"` + r + `"`
	} else if atom.IsType(TPair) {
		r := "("
		i := 0
		for pair := (*Pair)(atom.Data); pair != nil; i, pair = i+1, (*Pair)(pair.Cdr.Data) {
			r += Stringify(pair.Car) + " "
			if !pair.Cdr.IsType(TPair) {
				i += 1
				r += "." + " " + Stringify(pair.Cdr) + " "
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
