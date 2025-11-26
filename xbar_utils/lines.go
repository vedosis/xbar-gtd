package xbar_utils

import (
	"xbar/utils"
)

type XBarLine struct {
	message   string
	command   *XBarLineCommand
	href      string
	color     string
	fontName  string
	fontSize  int
	maxLength int    // Truncates string with '...'
	image     string // Base64 encoded
	disabled  bool
	children  []*XBarLine
}

type XBarLineCommand struct {
	command string
	args    []string
}

type XBarLineOption func(*XBarLine)

func NewXBarLine(message string, opts ...XBarLineOption) *XBarLine {
	defaultObj := XBarLine{
		message: message,
	}
	for _, o := range opts {
		o(&defaultObj)
	}
	return &defaultObj
}

func (l *XBarLine) Clone(message string, opts ...XBarLineOption) *XBarLine {
	values := *l

	if message != "" {
		values.message = message
	}

	for _, o := range opts {
		o(&values)
	}
	return &values
}

func WithMessage(message string) XBarLineOption {
	return func(line *XBarLine) {
		line.message = message
	}
}

func WithHref(href string) XBarLineOption {
	return func(line *XBarLine) {
		line.href = href
	}
}

func WithColor(color string) XBarLineOption {
	return func(line *XBarLine) {
		line.color = utils.Color(color)
	}
}

func WithCommand(args ...string) XBarLineOption {
	return func(line *XBarLine) {
		command := &XBarLineCommand{
			command: args[0],
			args:    args[1:],
		}
		line.command = command
	}
}

func WithFontName(name string) XBarLineOption {
	return func(line *XBarLine) {
		line.fontName = name
	}
}
func WithFontSize(size int) XBarLineOption {
	return func(line *XBarLine) {
		line.fontSize = size
	}
}
func WithMaxLength(length int) XBarLineOption {
	return func(line *XBarLine) {
		line.maxLength = length
	}
}
func WithImage(image string) XBarLineOption {
	return func(line *XBarLine) {
		line.image = image
	}
}
func WithDisabled() XBarLineOption {
	return func(line *XBarLine) {
		line.disabled = true
	}
}

func WithChildren(children ...*XBarLine) XBarLineOption {
	return func(line *XBarLine) {
		line.children = children
	}
}

func (l *XBarLine) ChildrenAsInterfaces() []interface{} {
	result := make([]interface{}, len(l.children))
	for i, child := range l.children {
		result[i] = child
	}
	return result
}
