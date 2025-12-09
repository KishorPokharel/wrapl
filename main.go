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
	var pipeOut string

	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.StringVar(&commandTemplate, "command", "", "command with {{}} placeholder")
	flag.StringVar(&historyFile, "history-file", "", "path to history file")
	flag.StringVar(&pipeOut, "pipe-out", "", "pipe output through this command (e.g., 'glow' or 'jq .')")
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
		if pipeOut != "" {
			if err := runWithPipe(cmd, pipeOut, debug); err != nil {
				fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
			}
		} else {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
			}
		}
	}
}

// runWithPipe executes cmd and pipes its output through pipeCmd
func runWithPipe(cmd *exec.Cmd, pipeCmd string, debug bool) error {
	pipe := exec.Command("bash", "-c", pipeCmd)
	if debug {
		fmt.Println("PIPE:", pipeCmd)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	pipe.Stdin = stdout
	pipe.Stdout = os.Stdout
	pipe.Stderr = os.Stderr
	cmd.Stderr = os.Stderr

	if err := pipe.Start(); err != nil {
		return fmt.Errorf("failed to start pipe command: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start main command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		stdout.Close()
		pipe.Wait()
		return fmt.Errorf("main command failed: %w", err)
	}

	if err := pipe.Wait(); err != nil {
		return fmt.Errorf("pipe command failed: %w", err)
	}

	return nil
}
