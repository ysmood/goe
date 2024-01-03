// Encrypt the .env file to .env.goe for the github user ids in GOE_ENV_VIEWERS,
// so that the github users can decrypt the file with their private keys.
package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/ysmood/goe"
	"github.com/ysmood/goe/pkg/envparse"
)

func main() {
	envFile, err := goe.LookupFile(goe.DOTENV)
	if err != nil {
		panic(err)
	}

	viewers := getViewers(envFile)

	if len(viewers) == 0 {
		panic("No viewers found in env var " + goe.GOE_ENV_VIEWERS + " defined in the .env file")
	}

	encrypt(envFile, envFile+goe.GOE_FILE_EXT, viewers)
}

func getViewers(envFile string) []string {
	content, err := os.ReadFile(envFile)
	if err != nil {
		panic(err)
	}

	ps, err := envparse.Parse(bytes.NewReader(content))
	if err != nil {
		panic(err)
	}

	viewers := []string{}
	conf := ps.Get(goe.GOE_ENV_VIEWERS)

	if conf != "" {
		viewers = strings.Split(conf, ",")
	}

	for i := range viewers {
		viewers[i] = strings.TrimSpace(viewers[i])
	}

	return viewers
}

func encrypt(envFile, goeFile string, viewers []string) {
	args := []string{"-i", envFile, "-o", goeFile}
	for _, v := range viewers {
		args = append(args, "-e", v)
	}

	whisperExec(args...)
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
