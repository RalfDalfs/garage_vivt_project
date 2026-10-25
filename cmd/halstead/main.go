// halstead применяет описанную в README лексическую методику подсчёта.
package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"math"
	"os"
)

type Result struct {
	File      string             `json:"file"`
	Operators map[string]int     `json:"operators"`
	Operands  map[string]int     `json:"operands"`
	Metrics   map[string]float64 `json:"metrics"`
}

func analyze(name string, src []byte) (Result, error) {
	r := Result{File: name, Operators: map[string]int{}, Operands: map[string]int{}}
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, name, src, parser.AllErrors); err != nil {
		return r, err
	}
	var s scanner.Scanner
	s.Init(fset.AddFile(name, -1, len(src)), src, nil, 0)
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		// Вставленные сканером точки с запятой не являются явным кодом.
		if tok == token.SEMICOLON && lit == "\n" {
			continue
		}
		if tok.IsLiteral() {
			r.Operands[lit]++
		} else {
			r.Operators[tok.String()]++
		}
	}
	var N1, N2 float64
	for _, v := range r.Operators {
		N1 += float64(v)
	}
	for _, v := range r.Operands {
		N2 += float64(v)
	}
	n1, n2 := float64(len(r.Operators)), float64(len(r.Operands))
	if n1 == 0 || n2 == 0 {
		return r, fmt.Errorf("пустой словарь")
	}
	n, N := n1+n2, N1+N2
	V := N * math.Log2(n)
	D := n1 / 2 * N2 / n2
	E := V * D
	r.Metrics = map[string]float64{
		"n1": n1, "n2": n2, "N1": N1, "N2": N2, "n": n, "N": N,
		"N_hat":   n1*math.Log2(n1) + n2*math.Log2(n2),
		"n2_star": 4, "V_star": 6 * math.Log2(6), "V": V, "D": D,
		"L": 1 / D, "E": E, "B_V": V / 3000, "B_E": math.Pow(E, 2.0/3) / 3000,
		"T_seconds": E / 18, "T_hours": E / 18 / 3600,
	}
	return r, nil
}

func run(paths []string) error {
	if len(paths) == 0 {
		return fmt.Errorf("укажите пути к .go файлам")
	}
	var results []Result
	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		r, err := analyze(path, src)
		if err != nil {
			return err
		}
		results = append(results, r)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
