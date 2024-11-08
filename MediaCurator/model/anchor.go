package model

import "fmt"

type Anchor struct {
	Text string
	HRef string
}

func NewAnchor(text, href string) Anchor {
	return Anchor{
		Text: text,
		HRef: href,
	}
}

func (a Anchor) AsString() string {
	return fmt.Sprintf("%s: %s", a.Text, a.HRef)
}
