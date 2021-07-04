package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// printTester 標準出力のテストをする
type printTester struct {
	t            *testing.T
	r, w, stdout *os.File
}

func (p *printTester) setupPrintTest() {
	r, w, err := os.Pipe()
	if err != nil {
		p.t.Fatal(err)
	}

	p.r, p.w = r, w
	p.stdout = os.Stdout
	os.Stdout = w
}

func (p *printTester) donePrintTest() string {
	os.Stdout = p.stdout
	err := p.w.Close()
	if err != nil {
		p.t.Fatal(err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, p.r)
	if err != nil {
		p.t.Fatal(err)
	}

	return buf.String()
}
