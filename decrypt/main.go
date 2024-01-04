package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/ysmood/goe"
)

func main() {
	if file, err := goe.LookupFile(goe.DOTENV + goe.GOE_FILE_EXT); err == nil {
		whisperExec("-i", file, "-o", strings.TrimSuffix(file, goe.GOE_FILE_EXT))
	}
}

func whisperExec(args ...string) {
	args = append([]string{"run", goe.WHISPER}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
