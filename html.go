package main

import (
	"fmt"
	"html"
	"strings"

	"github.com/russross/blackfriday"
)

// Renderer renders an element to HTML
type Renderer interface {
	Render() string
}

// Text represents plain text of html
type Text string

// Render for text, just escapes HTML entities
func (t Text) Render() string {
	return html.EscapeString(string(t))
}

// RawText represents Raw text of html
type RawText string

// Render for text, just escapes HTML entities
func (t RawText) Render() string {
	return string(t)
}

// MDText represents markdown formatted text
type MDText string

// Render for text, expands markdown formats
func (t MDText) Render() string {
	return string(blackfriday.MarkdownCommon([]byte(t)))
}

// Prop represents an HTML property
type Prop struct {
	Name  string
	Value string
}

// Render for property
func (p Prop) Render() string {
	txt := strings.Replace(p.Value, "\\", "\\\\", -1)
	txt = strings.Replace(txt, `"`, `\"`, -1)
	txt = strings.Replace(txt, "\n", `\n`, -1)
	txt = strings.Replace(txt, "\r", `\r`, -1)
	txt = strings.Replace(txt, "\t", `\t`, -1)
	return fmt.Sprintf(`%s="%s"`, p.Name, txt)
}

// Element represents an HTML element
type Element struct {
	Tag     string
	Props   []Prop
	Content []Renderer
}

// AddChild adds a child element to content
func (e *Element) AddChild(c Renderer) *Element {
	if e.Content == nil {
		e.Content = []Renderer{c}
		return e
	}
	e.Content = append(e.Content, c)
	return e
}

// GetContent returns rendered content
func (e *Element) GetContent() string {
	ret := make([]string, 0, len(e.Content))
	for _, c := range e.Content {
		ret = append(ret, c.Render())
	}

	return strings.Join(ret, "\n")
}

// Render for element
func (e *Element) Render() string {
	ret := `<` + e.Tag

	if len(e.Props) > 0 {
		for _, p := range e.Props {
			ret += " " + p.Render()
		}
	}

	if len(e.Content) < 1 {
		return ret + ` />`
	}

	ret += ">" + e.GetContent()
	return ret + "</" + e.Tag + `>`
}
