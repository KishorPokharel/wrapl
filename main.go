package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s '<command with {{}} placeholder>'\n", os.Args[0])
		os.Exit(1)
	}

	template := os.Args[1]
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("REPL started. Type 'exit' to quit.")

	for {
		input, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		if input == "exit" {
			break
		}

		cmdStr := strings.Replace(template, "{{}}", input, -1)
		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
