package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor(t *testing.T) {
	oldMode := IsDarkMode
	defer func() { IsDarkMode = oldMode }()
	IsDarkMode = true
	assert.Equal(t, "darkred", Color("red"))
	IsDarkMode = false
	assert.Equal(t, "crimson", Color("red"))

	assert.Equal(t, "totally-legit-color", Color("totally-legit-color"))
}
