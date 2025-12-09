package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	var debug bool
	var commandTemplate string
	var historyFile string

	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.StringVar(&commandTemplate, "command", "", "command with {{}} placeholder")
	flag.StringVar(&historyFile, "history-file", "", "path to history file")
	flag.Parse()

	if commandTemplate == "" {

		fmt.Fprintf(os.Stderr, "Error: empty command\n")
		fmt.Fprintf(os.Stderr, "\nUsage: \n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	if historyFile == "" {
		historyFile = filepath.Join(homeDir, ".wrapl_history")
	}

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

		cmdStr := strings.Replace(commandTemplate, "{{}}", line, 1)
		if debug {
			fmt.Println("CMD:", cmdStr)
		}
		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
		}
	}
}
