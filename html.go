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
	r := blackfriday.HtmlRenderer(
		0|
			blackfriday.HTML_USE_XHTML|
			blackfriday.HTML_USE_SMARTYPANTS|
			blackfriday.HTML_SMARTYPANTS_FRACTIONS|
			blackfriday.HTML_SMARTYPANTS_DASHES|
			blackfriday.HTML_SMARTYPANTS_LATEX_DASHES,
		"", "",
	)
	return string(blackfriday.Markdown(
		[]byte(t),
		r,
		0|
			blackfriday.EXTENSION_NO_INTRA_EMPHASIS|
			blackfriday.EXTENSION_TABLES|
			blackfriday.EXTENSION_FENCED_CODE|
			blackfriday.EXTENSION_AUTOLINK|
			blackfriday.EXTENSION_STRIKETHROUGH|
			blackfriday.EXTENSION_SPACE_HEADERS|
			blackfriday.EXTENSION_HEADER_IDS|
			blackfriday.EXTENSION_BACKSLASH_LINE_BREAK|
			blackfriday.EXTENSION_DEFINITION_LISTS|
			blackfriday.EXTENSION_FOOTNOTES,
	))
}

// Prop represents an HTML property
type Prop struct {
	Name  string
	Value string
}

// Render for property
func (p Prop) Render() string {
	txt := strings.ReplaceAll(p.Value, "\\", "\\\\")
	txt = strings.ReplaceAll(txt, `"`, `\"`)
	txt = strings.ReplaceAll(txt, "\n", `\n`)
	txt = strings.ReplaceAll(txt, "\r", `\r`)
	txt = strings.ReplaceAll(txt, "\t", `\t`)
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

// AppendProp adds a property to element, insteads of overwritting, add to tail if exists
func (e *Element) AppendProp(n, v string) *Element {
	idx := -1
	for k, p := range e.Props {
		if p.Name == n {
			idx = k
			break
		}
	}
	if idx == -1 {
		e.Props = append(e.Props, Prop{Name: n, Value: v})
		return e
	}

	e.Props[idx].Value += " " + v

	return e
}

// AddClass adds a css class to this element
func (e *Element) AddClass(c string) *Element {
	return e.AppendProp("class", c)
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
