package prs

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderSetupError(t *testing.T) {
	oldStdOut := os.Stdout

	r, w, _ := os.Pipe()
	defer func() {
		os.Stdout = oldStdOut
		r.Close()
	}()

	os.Stdout = w

	renderSetupError(fmt.Errorf("test error"))
	w.Close()

	output, _ := io.ReadAll(r)
	text := string(output)
	runes := []rune(text)
	t.Log(runes)
	assert.Contains(t, runes, []rune("⚠")[0])
}
