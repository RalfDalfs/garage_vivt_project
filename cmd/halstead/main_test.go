package main

import (
	"reflect"
	"testing"
)

func TestKnownTokens(t *testing.T) {
	r, err := analyze("tiny.go", []byte("package p\nfunc main() { x := 1; x++ } // comment\n"))
	if err != nil {
		t.Fatal(err)
	}
	ops := map[string]int{"package": 1, "func": 1, "(": 1, ")": 1, "{": 1, "}": 1, ":=": 1, ";": 1, "++": 1}
	args := map[string]int{"p": 1, "main": 1, "x": 2, "1": 1}
	if !reflect.DeepEqual(r.Operators, ops) || !reflect.DeepEqual(r.Operands, args) {
		t.Fatalf("unexpected dictionaries: %+v", r)
	}
	if r.Metrics["N1"] != 9 || r.Metrics["N2"] != 5 {
		t.Fatal(r.Metrics)
	}
}

func TestCommentAndWhitespaceInvariant(t *testing.T) {
	a, _ := analyze("a.go", []byte("package p\nvar x = 1\n"))
	b, _ := analyze("b.go", []byte("// заголовок\npackage p\n\nvar x = 1 // число\n"))
	if !reflect.DeepEqual(a.Metrics, b.Metrics) {
		t.Fatal("format affects metrics")
	}
}

func TestRejectInvalidSource(t *testing.T) {
	if _, err := analyze("bad.go", []byte("package p\nfunc {")); err == nil {
		t.Fatal("invalid source accepted")
	}
}
