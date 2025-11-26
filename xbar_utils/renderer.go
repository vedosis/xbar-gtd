package xbar_utils

import (
	"fmt"
	"strings"
)

type XBarRenderer struct {
	header       *XBarHeader
	font         string
	sectionDepth uint8
}

func NewXBarRenderer() *XBarRenderer {
	return &XBarRenderer{
		header: NewXBarHeader("No Title"),
	}
}

func (r *XBarRenderer) Output(o ...interface{}) {
	if !r.header.hasRendered {
		if r.header.icon != "" {
			fmt.Print(fmt.Sprintf("%s  ", r.header.icon))
		}
		fmt.Print(r.header.title)
		if r.header.image != "" {
			fmt.Print(fmt.Sprintf(" | icon=data:image/png;base64,%s", r.header.image))
		}
		if r.font != "" {
			fmt.Print(fmt.Sprintf(" | font='%s'", r.font))
		}
		fmt.Println()
		r.header.hasRendered = true
		fmt.Println("---")
	}

	prefix := ""
	if r.sectionDepth > 0 {
		prefix += strings.Repeat("--", int(r.sectionDepth))
		prefix += " "
	}

	for _, line := range o {
		switch v := line.(type) {
		case string:
			fmt.Println(prefix + v)
		case *XBarLine:
			fmt.Println(prefix + r.RenderLine(v))
			if len(v.children) > 0 {
				r.sectionDepth += 1
				r.Output(v.ChildrenAsInterfaces()...)
				r.sectionDepth -= 1
			}
		default:
			fmt.Println(fmt.Sprintf("unknown object, %s", v))
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
	if line.fontName != "" {
		rVal += fmt.Sprintf(" | font=%s", line.fontName)
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

func (r *XBarRenderer) SetFont(font string) {
	r.font = font
}
