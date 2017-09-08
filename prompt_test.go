package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) (string, string) {
	oldErr := os.Stderr // keep backup of the real stdout
	rE, wE, _ := os.Pipe()
	os.Stderr = wE

	oldOut := os.Stdout // keep backup of the real stdout
	rO, wO, _ := os.Pipe()
	os.Stdout = wO

	f()

	outCerr := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rE)
		outCerr <- buf.String()
	}()

	outCout := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rO)
		outCout <- buf.String()
	}()

	// back to normal state
	wE.Close()
	wO.Close()
	os.Stdout = oldOut // restoring the real stdout
	os.Stderr = oldErr // restoring the real stderr
	outE := <-outCerr
	outO := <-outCout
	return outO, outE
}

func TestMoveCursorUp(t *testing.T) {
	want := "\033[3A"
	_, have := captureOutput(func() {
		moveCursorUp(3)
	})
	assert.Equal(t, want, have)
}

func TestPrintChoice(t *testing.T) {
	want := "Message\n\x1b[36m → opt1\x1b[0m\n   opt2\n   opt3\n"
	_, have := captureOutput(func() {
		printChoice("Message", 0, []string{"opt1", "opt2", "opt3"})
	})
	assert.Equal(t, want, have)
}

func TestSelectUp(t *testing.T) {
	var choice int
	newChoice := selectUp("Message", choice, []string{"opt1", "opt2", "opt3"})
	assert.Equal(t, choice+1, newChoice)
}

func TestSelectUpShouldRotate(t *testing.T) {
	options := []string{"opt1", "opt2", "opt3"}
	newChoice := selectUp("Message", len(options)-1, options)
	assert.Equal(t, 0, newChoice)
}

func TestSelectDown(t *testing.T) {
	choice := 2
	newChoice := selectDown("Message", choice, []string{"opt1", "opt2", "opt3"})
	assert.Equal(t, choice-1, newChoice)
}

func TestSelectDownShouldRotate(t *testing.T) {
	options := []string{"opt1", "opt2", "opt3"}
	newChoice := selectDown("Message", 0, options)
	assert.Equal(t, len(options)-1, newChoice)
}
