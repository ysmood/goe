// Package dotenv loads environment variables from a .env file.
package dotenv

import (
	"github.com/joho/godotenv"
	"github.com/ysmood/dotenv/pkg/utils"
)

func init() {
	_ = godotenv.Load(utils.LookupFile(".env"))
}
