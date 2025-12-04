package xbar_utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleLongText = strings.Join(strings.Fields(`
 The workshop door burst open as Eksbaur's crystalline form flickered urgently. "Master, your cryptocurrency cauldrons are bubbling over—Bitcoin up twelve percent, Ethereum requiring immediate attention, and
   that obscure potion coin you invested in has vanished entirely," the familiar announced, its voice shifting between concerned and smugly judgmental. Merlin squinted at the floating runes displaying his
  GitHub notifications: forty-seven pull requests awaiting review, sixteen issues marked critical, and someone had commented "needs more cowbell" on his ancient spell-casting library. "And the temperature
  sensors?" he asked wearily, knowing the answer. "Your MacBook is approaching molten lava temperatures while running those neural network divinations," Eksbaur chimed, spawning a miniature weather widget
  that showed storm clouds gathering over the data center where his cloud spells resided. The familiar's form pulsed red as RSS feeds flooded in—the Magic Council had published new regulations, Stack Overflow
   had three answers to his levitation optimization question, and his favorite podcast about interdimensional travel had just dropped. "Master," Eksbaur said gently, "you have sixty-two items in your menu
  bar, your battery is at four percent, and you haven't looked away from that code spell in three hours." Merlin sighed and reached for his cold coffee, grateful that at least Eksbaur kept everything
  organized at a glance.
`), " ")

func TestXBarRenderer_Output(t *testing.T) {
	var buffer []string
	renderer := NewXBarRenderer()

	renderer.SetTitle("Test Title")
	renderer.SetIcon("💕")
	renderer.SetHeaderImage("image.png")
	renderer.SetFont("HelveticaOld")

	var printBuffer string
	renderer.cmdPrintLnFn = func(a ...interface{}) (int, error) {
		if len(a) == 0 {
			a = append(a, "")
		}
		outputString := a[0].(string)
		if printBuffer != "" {
			outputString += printBuffer
			printBuffer = ""
		}
		buffer = append(buffer, outputString)
		return len(outputString), nil
	}

	renderer.cmdPrintFn = func(a ...interface{}) (int, error) {
		if len(a) == 0 {
			return 0, nil
		}
		value := a[0].(string)
		printBuffer += value
		return len(value), nil
	}

	baseLine := NewXBarLine("Test XBarLine", WithColor("red"))
	renderer.Output(
		"Test String Line",
		baseLine,
		baseLine.Clone("Test Clone Line", WithFontSize(100)),
		NewXBarLine(sampleLongText,
			WithColor("secondary"),
			WithCommand("echo", "this is a test", "OtherArgs"),
			WithFontName("HelveticaNeue"),
			WithFontSize(12),
			WithHref("https://www.google.com"),
			WithImage("image.png"),
			WithDisabled(),
			WithChildren(
				NewXBarLine("Child Line", WithChildren(
					NewXBarLine("Child Line Child"))),
			),
			WithMaxLength(210),
			WithWrapTextLength(200),
		),
		NewXBarDebugError("Test Error", fmt.Errorf("extra special cool error")),
		NewXBarLine("osascript call", WithCommand(OpenInTerminal("echo done")...)),
		struct {
			Name string
		}{
			Name: "Random Struct",
		},
	)

	if printBuffer != "" {
		renderer.cmdPrintLnFn(printBuffer)
	}

	var lines [][]string
	for _, line := range buffer {
		parts := strings.Split(line, "|")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		lines = append(lines, parts)
	}
	outputRunes := []rune(lines[0][0])
	assert.Equal(t, []rune("💕")[0], outputRunes[0])
	assert.Contains(t, lines[0][0], "Test Title")
	assert.Contains(t, lines[0], "font='HelveticaOld'")

	assert.Equal(t, []string{"---"}, lines[1])
	assert.Equal(t, []string{"Test String Line"}, lines[2])
	assert.Equal(t, []string{"Test XBarLine", "color=darkred", "font='HelveticaOld'"}, lines[3])
	assert.Equal(t, []string{"Test Clone Line", "color=darkred", "size=100", "font='HelveticaOld'"}, lines[4])

	assert.LessOrEqual(t, len(lines[5][0]), 200)
	assert.Contains(t, lines[5][0], "Eksbaur")
	var overrideCounter = map[string]int{}
	for _, line := range lines[5][1:] {
		overrideType := strings.Split(line, "=")[0]
		overrideCounter[overrideType]++
		switch overrideType {
		case "shell":
			assert.Equal(t, "shell=echo param1='this is a test' param2=OtherArgs", line, "shell param escaping")
		case "color":
			assert.Equal(t, "color=#333030", line)
		case "href":
			assert.Equal(t, "href=https://www.google.com", line)
		case "font":
			assert.Equal(t, "font=HelveticaNeue", line)
		case "size":
			assert.Equal(t, "size=12", line)
		case "image":
			assert.Equal(t, "image=image.png", line)
		case "disabled":
			assert.Equal(t, "disabled=true", line)
		case "length":
			assert.Equal(t, "length=210", line)
		default:
			t.Errorf("unexpected override type: %s", overrideType)
		}
	}
	expectedKeys := []string{"shell", "color", "href", "font", "size", "image", "disabled"}
	for _, key := range expectedKeys {
		_, ok := overrideCounter[key]
		assert.True(t, ok, "missing %s override", key)
		assert.Equal(t, 1, overrideCounter[key], "duplicate %s override", key)
	}

	assert.Equal(t, "-- ", lines[13][0][0:3])
	assert.Equal(t, "---- ", lines[14][0][0:5])
	assert.Equal(t, []string{"error details", "color=#333030"}, lines[15][:2])
	assert.Equal(t, "-- Test Error: extra special cool error", lines[16][0])
	assert.Equal(t, "shell=osascript param1=-e", lines[17][1][:25])
	assert.Contains(t, lines[17][1], "\"echo done\"")

	assert.Equal(t, []string{"unknown object, {Random Struct}"}, lines[18])
}
