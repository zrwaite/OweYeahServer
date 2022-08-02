package tests

import (
	"testing"

	"github.com/zrwaite/OweMate/terminal"
)

func PrintErr(t *testing.T, err string) {
	t.Error(terminal.Red + err + terminal.Reset)
}
func PrintPass(t *testing.T, msg string) {
	t.Log(terminal.Green + "PASS: " + msg + terminal.Reset)
}
