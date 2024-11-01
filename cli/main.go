package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/ysmood/goe"
	_ "github.com/ysmood/goe/load"
	"github.com/ysmood/goe/pkg/utils"
)

var (
	errGetShell = errors.New("failed to get shell")
	errFailRun  = errors.New("failed to run shell")
)

func main() {
	shell, err := Shell()
	if err != nil {
		panic(fmt.Errorf("%w: %w", errGetShell, err))
	}

	cmd := exec.Command(shell)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		var exitErr exec.ExitError
		if errors.Is(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		panic(fmt.Errorf("%w: %w", errFailRun, err))
	}

	utils.Println(goe.Prefix, "Unloaded environment variables")
}
