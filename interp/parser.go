package interp

import "strconv"

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
		return string(*(*String)(atom.Data))
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
	} else { // Procedure or build-in
		return "#<procedure>"
	}
}