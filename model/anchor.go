package model

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
