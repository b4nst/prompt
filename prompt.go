package main

import (
	"fmt"
	"os"

	"github.com/pkg/term"
)

func moveCursorUp(bias int) {
	fmt.Fprintf(os.Stderr, "\033[%dA", bias)
}

func printChoice(message string, curr int, options []string) {
	fmt.Fprintln(os.Stderr, message)
	for idx, opt := range options {
		var str string
		if idx == curr {
			str = "\033[36m" + " → " + opt + "\033[0m"
		} else {
			str = "   " + opt
		}
		fmt.Fprintln(os.Stderr, str)
	}
}

func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}

func selectUp(message string, choice int, options []string) (newchoice int) {
	if choice < (len(options) - 1) {
		newchoice = choice + 1
	} else {
		newchoice = 0
	}
	moveCursorUp(len(options) + 1)
	printChoice(message, newchoice, options)
	return newchoice
}

func selectDown(message string, choice int, options []string) (newchoice int) {
	if choice > 0 {
		newchoice = choice - 1
	} else {
		newchoice = len(options) - 1
	}
	moveCursorUp(len(options) + 1)
	printChoice(message, newchoice, options)
	return newchoice
}

func printUsage() {
	fmt.Println("Usage :\n\tprompt Message option1 option2 [...as much options you want]")
	fmt.Println()
	fmt.Println("\tExample")
	fmt.Println("\tprompt \"Choose one:\" Farnsworth Rick \"Nerdelbaum Frink, Jr\"")
}

func main() {
	if len(os.Args) < 4 {
		printUsage()
		os.Exit(1)
	}
	args := os.Args[1:]
	options := args[1:]
	message := args[0]
	var choice int

	printChoice(message, choice, options)
	ascii, keycode, err := getChar()
	for ascii != 13 && ascii != 3 {
		if err != nil {
			continue
		}

		if keycode != 0 {
			switch keycode {
			case 40, 39:
				choice = selectUp(message, choice, options)
			case 38, 37:
				choice = selectDown(message, choice, options)
			}
		} else {
			switch ascii {
			case 9, 32:
				choice = selectUp(message, choice, options)
			case 127, 0:
				choice = selectDown(message, choice, options)
			}
		}

		ascii, keycode, err = getChar()
	}

	if ascii == 13 {
		fmt.Println(options[choice])
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
