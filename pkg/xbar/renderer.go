package xbar

import (
	"fmt"
	"os"
	"strings"
	"xbar/pkg/utils"
)

type XBarRenderer struct {
	header       *XBarHeader
	font         string
	sectionDepth uint8
	cmdPrintFn   func(...any) (int, error)
	cmdPrintLnFn func(...any) (int, error)
}

type XBarRendererOption func(*XBarRenderer)

func NewXBarRenderer(opts ...XBarRendererOption) *XBarRenderer {
	renderer := &XBarRenderer{
		header:       NewXBarHeader("No Title"),
		cmdPrintFn:   fmt.Print,
		cmdPrintLnFn: fmt.Println,
	}

	for _, opt := range opts {
		opt(renderer)
	}

	return renderer
}

func WithCmdPrintFn(fn func(...any) (int, error)) XBarRendererOption {
	return func(r *XBarRenderer) {
		r.cmdPrintFn = fn
	}
}

func WithCmdPrintLnFn(fn func(...any) (int, error)) XBarRendererOption {
	return func(r *XBarRenderer) {
		r.cmdPrintLnFn = fn
	}
}

func (r *XBarRenderer) Output(o ...interface{}) {
	os.Stdout.Sync()
	if !r.header.hasRendered {
		if r.header.icon != "" {
			r.cmdPrintFn(fmt.Sprintf("%s  ", r.header.icon))
		}
		r.cmdPrintFn(r.header.title)
		if r.header.image != "" {
			r.cmdPrintFn(fmt.Sprintf(" | icon=data:image/png;base64,%s", r.header.image))
		}
		if r.font != "" {
			r.cmdPrintFn(fmt.Sprintf(" | font='%s'", r.font))
		}
		r.cmdPrintLnFn()
		r.header.hasRendered = true
		r.cmdPrintLnFn("---")
	}

	prefix := ""
	if r.sectionDepth > 0 {
		prefix += strings.Repeat("--", int(r.sectionDepth))
		prefix += " "
	}

	for _, line := range o {
		switch v := line.(type) {
		case string:
			r.cmdPrintLnFn(prefix + v)
		case *XBarLine:
			var lines []string
			if v.wrapTextLength == 0 {
				lines = []string{v.message}
			} else {
				lines = utils.WrapLines(v.message, v.wrapTextLength)
			}
			for _, sentence := range lines {
				r.cmdPrintLnFn(prefix + r.RenderLine(v.Clone(sentence)))
			}
			if len(v.children) > 0 {
				r.sectionDepth += 1
				r.Output(v.ChildrenAsInterfaces()...)
				r.sectionDepth -= 1
			}
		default:
			r.cmdPrintLnFn(fmt.Sprintf("unknown object, %s", v))
		}
	}
}

func (r *XBarRenderer) RenderLine(line *XBarLine) string {
	rVal := line.message
	if line.command != nil {
		rVal += fmt.Sprintf(" | shell=%s", line.command.command)
		for idx, param := range line.command.args {
			if strings.Contains(param, " ") && param[0] != '\'' {
				param = "'" + param + "'"
			}
			rVal += fmt.Sprintf(" param%d=%s", idx+1, param)
		}
	}
	if line.color != "" {
		rVal += fmt.Sprintf(" | color=%s", line.color)
	}
	if line.href != "" {
		rVal += fmt.Sprintf(" | href=%s", line.href)
	}
	if line.fontSize != 0 {
		rVal += fmt.Sprintf(" | size=%d", line.fontSize)
	}
	if line.maxLength != 0 {
		rVal += fmt.Sprintf(" | length=%d", line.maxLength)
	}
	if line.image != "" {
		rVal += fmt.Sprintf(" | image=%s", line.image)
	}
	if line.disabled {
		rVal += fmt.Sprintf(" | disabled=true")
	}
	if line.fontName != "" {
		rVal += fmt.Sprintf(" | font=%s", line.fontName)
	} else if r.font != "" {
		rVal += fmt.Sprintf(" | font='%s'", r.font)
	}
	return rVal
}

func (r *XBarRenderer) SetFont(font string) *XBarRenderer {
	r.font = font
	return r
}

func (r *XBarRenderer) SetTitle(title string) *XBarRenderer {
	r.header.title = title
	return r
}

func (r *XBarRenderer) SetIcon(icon string) *XBarRenderer {
	r.header.icon = icon
	return r
}

func (r *XBarRenderer) SetHeaderImage(image string) *XBarRenderer {
	r.header.image = image
	return r
}
