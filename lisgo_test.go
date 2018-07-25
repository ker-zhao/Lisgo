package main

import (
	"lisgo/interp"
	"lisgo/parser"
	"testing"
)

func Test_InterP(t *testing.T) {
	code1 := `(let ([x 2])
   (let ([f (lambda (y) (* x y))])
     (let ([x 4])
       (f 3))))`
	expect1 := `6`
	a := interp.InterP(parser.Parse(code1), interp.GlobalEnv)
	b := interp.InterP(parser.Parse(expect1), interp.GlobalEnv)
	if !interp.AtomEqual(a, b) {
		t.Errorf("Interpreter error. input: %s, returns value: %s", code1,
			interp.Stringify(a))
	} else {
		t.Log("OK")
	}

	code2 := "(* 2 3)"
	wrong2 := "(+ 2 3 2)"
	a2 := interp.InterP(parser.Parse(code2), interp.GlobalEnv)
	b2 := interp.InterP(parser.Parse(wrong2), interp.GlobalEnv)
	if interp.AtomEqual(a2, b2) {
		t.Errorf("Interpreter error. input: %s equals  %s", code2, wrong2)
	} else {
		t.Log("OK")
	}

	cons3 := "(car (cdr (cons 1 (cons 2 3))))"
	expect3 := "(car (cdr (append (list 1) (list 2 3 4))))"
	a3 := interp.InterP(parser.Parse(cons3), interp.GlobalEnv)
	b3 := interp.InterP(parser.Parse(expect3), interp.GlobalEnv)
	if interp.AtomEqual(a3, b3) {
		t.Log("OK")
	} else {
		t.Errorf("Error: cons3 \n")
	}
}

func Test_InterP2(t *testing.T) {
	code1 := "(equal? (list 1 2 3) (quote (1 2 3)))"
	code2 := "#t"
	r1 := interp.InterP(parser.Parse(code1), interp.GlobalEnv)
	r2 := interp.InterP(parser.Parse(code2), interp.GlobalEnv)
	if interp.AtomEqual(r1, r2) {
		t.Log("OK")
	} else {
		t.Errorf("Error: Test_InterP2 \n")
	}
}

func Test_InterP3(t *testing.T) {
	code1 := `
(begin
(define a (list 1 2 3))
(define b a)
(eq? b (list 1 2 3))
)
`
	code2 := "#f"
	r1 := interp.InterP(parser.Parse(code1), interp.GlobalEnv)
	r2 := interp.InterP(parser.Parse(code2), interp.GlobalEnv)
	if interp.AtomEqual(r1, r2) {
		t.Log("OK")
	} else {
		t.Errorf("Error: Test_InterP3 \n")
	}
}