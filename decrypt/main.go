package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ysmood/goe"
)

func main() {
	if file, err := goe.LookupFile(goe.DOTENV + goe.GOE_FILE_EXT); err == nil {
		meta, err := exec.Command("go", "run", goe.WHISPER, "-i", file, "-m").CombinedOutput() //nolint: gosec
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}

		fmt.Println(string(meta))

		whisperExec("-i", file, "-o", strings.TrimSuffix(file, goe.GOE_FILE_EXT))
	}
}

func whisperExec(args ...string) {
	args = append([]string{"run", goe.WHISPER}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if cmd.Run() != nil {
		os.Exit(1)
	}
}
