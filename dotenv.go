// Package dotenv loads environment variables from a .env file.
package dotenv

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(findDotenv(".env"))
	if err != nil {
		panic(err)
	}
}

// recursively search for .env.
func findDotenv(name string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prev := ""

	for dir != prev {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p
		}

		prev = dir
		dir = filepath.Dir(dir)
	}

	return ""
}
