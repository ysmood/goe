package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/ysmood/goe"
	"github.com/ysmood/goe/pkg/utils"
)

var (
	errGetShell = errors.New("failed to get shell")
	errFailRun  = errors.New("failed to run shell")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), Usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	envFile := flag.Arg(0)

	if envFile == "" {
		envFile = goe.DOTENV
	}

	err := goe.AutoLoad(envFile)
	if err != nil {
		panic(err)
	}

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
