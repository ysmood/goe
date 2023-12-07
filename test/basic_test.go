package test_test

import (
	"os"
	"testing"

	_ "github.com/ysmood/dotenv"
)

func TestBasic(t *testing.T) {
	if os.Getenv("STR") != "hello" {
		t.Fail()
	}
}
