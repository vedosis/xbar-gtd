package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	paragraph = `Merlin the Magnificent stood amidst bubbling cauldrons and crackling enchantments, his workshop a chaos
  of half-finished spells and forgotten ingredients. "Eksbaur!" he called desperately, and his faithful familiar—a 
  shimmering spirit with eyes like hourglasses—materialized beside him. "The resurrection elixir needs stirring in three
  minutes, the transformation tonic is overheating, and the levitation potion requires moonstone dust before midnight,"
  Eksbaur recited calmly, its voice cutting through the magical din. As purple sparks began cascading from the ceiling,
  Eksbaur chimed a warning: "Ambient magic levels critical—you must ground the excess immediately." Merlin waved his
  staff, channeling the wild energy into a containment crystal, and sighed with relief. "What would I do without you?"
  he muttered, rushing to stir the resurrection elixir just as Eksbaur's internal timer reached zero. The familiar
  simply flickered with satisfaction, already tracking the next seventeen tasks that required the wizard's attention.`
)

func TestWrapLines(t *testing.T) {
	inputString := strings.Join(strings.Fields(paragraph), " ")
	output := WrapLines(inputString, 400)
	assert.Less(t, len(output[0]), 400, "Output line length should not exceed 400 characters")
	assert.Greater(t, len(output), 2, "Should break into multiple lines")
	assert.Contains(t, output[len(output)-1], "required the wizard's attention.", "Last line should contain the full paragraph")

	// Dealing with those shorter lines
	output = WrapLines(inputString, 5)
	assert.Equal(t, output[0], "Merli")
	assert.Equal(t, output[1], "n the")
}
