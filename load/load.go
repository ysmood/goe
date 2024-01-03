// Package load loads environment variables from a .env file.
// It uses [goe.Load] to load the .env file, override set to false, expand set to true.
package load

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/ysmood/goe"
)

type info struct{}

var prefix = fmt.Sprintf("[%s]", reflect.TypeOf(info{}).PkgPath())

func init() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, prefix, r)
			os.Exit(1)
		}
	}()

	loadGoeFile()

	err := goe.Load(false, true, goe.DOTENV)
	if err != nil {
		panic(err)
	}

	path, err := goe.LookupFile(goe.DOTENV)
	if err != nil {
		panic(err)
	}

	fmt.Println(prefix+" Loaded environment variables from:", path)
}

func loadGoeFile() {
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
