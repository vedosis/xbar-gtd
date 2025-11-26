package xbar_utils

type XBarHeader struct {
	title       string
	icon        string
	image       string
	hasRendered bool
}

func NewXBarHeader(title string) *XBarHeader {
	return NewXBarHeaderWithIcon(title, "")
}

func NewXBarHeaderWithIcon(title string, icon string) *XBarHeader {
	return &XBarHeader{
		title:       title,
		icon:        icon,
		hasRendered: false,
	}
}

func (r *XBarRenderer) SetTitle(title string) *XBarRenderer {
	if r.header == nil {
		r.header = NewXBarHeader("No Title")
	}
	r.header.title = title
	return r
}

func (r *XBarRenderer) SetIcon(icon string) *XBarRenderer {
	if r.header == nil {
		r.header = NewXBarHeader("No Icon")
	}
	r.header.icon = icon
	return r
}

func (r *XBarRenderer) SetHeaderImage(image string) *XBarRenderer {
	r.header.image = image
	return r
}
