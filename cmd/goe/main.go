package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/ysmood/goe"
	"github.com/ysmood/goe/pkg/utils"
)

var (
	errGetShell = errors.New("failed to get shell")
	errFailRun  = errors.New("failed to run shell")
)

func main() {
	envFile := goe.DOTENV

	if len(os.Args) == 2 && slices.Contains([]string{"-h", "--help", "-help"}, os.Args[1]) {
		utils.Println(USAGE)
		os.Exit(0)
	}

	if len(os.Args) > 1 {
		envFile = os.Args[1]
	}

	err := goe.AutoLoad(envFile)
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 { //nolint: mnd
		runShell()
	} else {
		runCommand()
	}
}

func runShell() {
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
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		panic(fmt.Errorf("%w: %w", errFailRun, err))
	}

	utils.Println(goe.Prefix, "Unloaded environment variables")
}

func runCommand() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		panic(fmt.Errorf("%w: %w", errFailRun, err))
	}
}
