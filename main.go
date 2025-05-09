package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s '<command with {{}} placeholder>'\n", os.Args[0])
		os.Exit(1)
	}

	template := os.Args[1]

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	historyFile := filepath.Join(homeDir, ".wrapl_history")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     historyFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Println("Type 'exit' to quit.")

	for {
		line, err := rl.Readline()
		if err != nil {
			break // EOF or interrupt
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			break
		}

		cmdStr := strings.Replace(template, "{{}}", line, -1)
		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Command error: %v\n", err)
		}
	}
}
