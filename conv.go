package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

func mkpage(n int) *Element {
	return &Element{
		Tag: "div",
		Props: []Prop{
			Prop{Name: "class", Value: "page"},
			Prop{Name: "id", Value: "page" + strconv.Itoa(n)},
		},
	}
}

// Page 代表一張投影片
type Page struct {
	Num     int
	Title   string
	Content string
	Notes   []string // presenter notes
}

// ToElement 把投影片轉成 HTML element tree
func (p *Page) ToElement(footer string) *Element {
	ret := &Element{Tag: "div"}
	ret.AddClass("page")

	head := (&Element{Tag: "header"}).AddClass("pageTitle").AddChild(
		(&Element{Tag: "h1"}).AddChild(Text(p.Title)),
	)
	body := (&Element{Tag: "section"}).AddClass("pageContent").AddChild(MDText(p.Content))
	foot := (&Element{Tag: "footer"}).AddClass("pageFoot").AddChild(MDText(footer))

	if p.Content == "" {
		head.AddClass("lonely")
	}

	ret.AddChild(head)
	if p.Content != "" {
		ret.AddChild(body)
	}

	if len(p.Notes) > 0 {
		n := &Element{Tag: "script"}
		n.AddChild(RawText(`notes=('undefined'==typeof notes)?{}:notes;notes["p` + strconv.Itoa(p.Num) + `"] = "` + url.QueryEscape(strings.Join(p.Notes, "\n")) + `";`))
		ret.AddChild(n)
	}

	return ret.AddChild(foot)
}

// Slides 代表整份投影片
type Slides struct {
	Title    string
	Subtitle string // optional

	// information about the author/presenter
	Name     string
	Email    string
	Notes    []string // optional
	PreNotes []string // optional

	// page customization
	Footer string // optional

	pages []*Page
}

func (s *Slides) makeFirstPage() *Element {
	ret := (&Element{Tag: "div"}).AddClass("pageContent")

	ret.AddChild((&Element{Tag: "div"}).AddClass("slideTitle").AddChild(
		(&Element{Tag: "h1"}).AddChild(Text(s.Title))),
	)
	if s.Subtitle != "" {
		ret.AddChild((&Element{Tag: "div"}).AddClass("slideSubtitle").AddChild(Text(s.Subtitle)))
	}
	ret.AddChild((&Element{Tag: "div"}).AddClass("authorName").AddChild(Text(s.Name)))
	if len(s.PreNotes) > 0 {
		for _, n := range s.PreNotes {
			ret.AddChild((&Element{Tag: "div"}).AddClass("authorNotes").AddChild(MDText(n)))
		}
	}

	return (&Element{Tag: "div"}).AddChild(ret).AddClass("page").AddClass("firstPage")
}

func (s *Slides) makeLastPage() *Element {
	ret := (&Element{Tag: "div"}).AddClass("pageContent")

	ret.AddChild((&Element{Tag: "h2"}).AddChild(Text("Thank you")))
	ret.AddChild((&Element{Tag: "div"}).AddClass("authorName").AddChild(MDText("**" + s.Name + "**")))
	ret.AddChild((&Element{Tag: "div"}).AddClass("authorMail").AddChild(MDText(s.Email)))
	if s.Notes != nil && len(s.Notes) > 0 {
		for _, n := range s.Notes {
			ret.AddChild((&Element{Tag: "div"}).AddClass("authorNotes").AddChild(MDText(n)))
		}
	}

	return (&Element{Tag: "div"}).AddChild(ret).AddClass("page").AddClass("lastPage")
}

// Len 取得頁數(不含第一及最後一張)
func (s *Slides) Len() int {
	return len(s.pages)
}

// ToElements 把整份投影片轉成一個一個的 HTML element
func (s *Slides) ToElements() []*Element {
	var e *Element
	ret := make([]*Element, 0, len(s.pages)+2)

	cnt := 1
	e = s.makeFirstPage()
	e.AppendProp("id", "page"+strconv.Itoa(cnt))
	cnt++
	ret = append(ret, e)
	for _, p := range s.pages {
		p.Num = cnt
		e = p.ToElement(s.Footer)
		e.AppendProp("id", "page"+strconv.Itoa(cnt))
		ret = append(ret, e)
		cnt++
	}
	e = s.makeLastPage()
	e.AppendProp("id", "page"+strconv.Itoa(cnt))
	cnt++
	ret = append(ret, e)

	return ret
}

func conv(fn, css string) ([]byte, error) {
	barr, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	s := &Slides{}

	supportedMetadata := []struct {
		tag     string
		process func(data string)
	}{
		{"TITLE", func(data string) { s.Title = data }},
		{"SUBTITLE", func(data string) { s.Subtitle = data }},
		{"AUTHOR", func(data string) { s.Name = data }},
		{"EMAIL", func(data string) {
			s.Email = fmt.Sprintf("[%s](mailto:%s)", data, data)
		}},
		{"FOOTER", func(data string) { s.Footer = data }},
		{"FACEBOOK", func(data string) {
			s.Notes = append(s.Notes, fmt.Sprintf("[Facebook](https://www.facebook.com/%s)", data))
		}},
		{"TWITTER", func(data string) {
			s.Notes = append(s.Notes, fmt.Sprintf("[Twitter @%s](https://twitter.com/%s)", data, data))
		}},
		{"URL", func(data string) {
			s.Notes = append(s.Notes, fmt.Sprintf("[%s](%s)", data, data))
		}},
		{"TEXT", func(data string) {
			s.Notes = append(s.Notes, data)
		}},
		{"TITLETEXT", func(data string) {
			s.PreNotes = append(s.PreNotes, data)
		}},
	}

	save := func(page *Page, buf string) *Page {
		if page != nil {
			buf = strings.TrimSpace(buf)
			if buf != "" {
				page.Content = buf
			}
			s.pages = append(s.pages, page)
		}

		return &Page{}
	}

	strs := strings.Split(string(barr), "\n")
	firstTitleLine := 0
	for k, l := range strs {
		if len(l) < 2 {
			continue
		}

		if l[0] != '#' {
			// plain text, skip
			continue
		}

		if l[1] != '+' {
			firstTitleLine = k
			break
		}

		l = l[2:]

		// 處理各種 metadata tag
		for _, meta := range supportedMetadata {
			if !strings.HasPrefix(l, meta.tag+":") {
				continue
			}

			meta.process(strings.TrimSpace(l[len(meta.tag)+1:]))
			break
		}
	}

	var page *Page
	buf := ""
	for _, l := range strs[firstTitleLine:] {
		sz := len(l)
		if sz < 2 {
			buf += l + "\n"
			continue
		}
		if l == "##" {
			buf += l + "\n"
			continue
		}
		if l[0] != '#' {
			// 非 headings
			buf += l + "\n"
			continue
		}

		if strings.HasPrefix(l, "#+NOTE:") {
			// presenter notes
			if page != nil {
				page.Notes = append(page.Notes, strings.TrimSpace(l[7:]))
				continue
			}
		}

		if len(l) >= 3 && l[1] == '#' && l[2] == '#' {
			// 至少會是 h3
			buf += l + "\n"
			continue
		}

		l = l[1:]
		if l[0] == '#' {
			l = l[1:]
		}

		l = strings.TrimSpace(l)
		page = save(page, buf)
		buf = ""
		page.Title = l
	}

	save(page, buf)

	html := &Element{Tag: "html"}
	head := &Element{Tag: "head"}
	html.AddChild(
		head.AddChild(
			(&Element{Tag: "meta"}).AppendProp("charset", "UTF-8"),
		).AddChild(
			(&Element{Tag: "meta"}).AppendProp(
				"name", "viewport",
			).AppendProp(
				"content", "width=device-width, initial-scale=1.0",
			),
		).AddChild(
			(&Element{Tag: "link"}).AppendProp(
				"rel", "stylesheet",
			).AppendProp(
				"href", "/assets/prism.css",
			),
		).AddChild(
			(&Element{Tag: "link"}).AppendProp(
				"rel", "stylesheet",
			).AppendProp(
				"href", "/assets/"+css+".css",
			),
		),
	)
	body := &Element{Tag: "body", Props: []Prop{Prop{Name: "id", Value: "body"}}}
	for _, e := range s.ToElements() {
		body.AddChild(e)
	}
	body.AddChild(
		(&Element{Tag: "script"}).AddChild(RawText("var maxPage=" + strconv.Itoa(s.Len()+2) + ";")),
	).AddChild(
		(&Element{Tag: "script"}).AppendProp("src", "/assets/js/before.js").AddChild(RawText("")),
	).AddChild(
		(&Element{Tag: "script"}).AppendProp("src", "/assets/prism.js").AddChild(RawText("")),
	).AddChild(
		(&Element{Tag: "script"}).AppendProp("src", "/assets/js/after.js").AddChild(RawText("")),
	).AddChild(
		(&Element{Tag: "div"}).AddChild(RawText("")).AppendProp("style", "height:10vh;"),
	)

	if css != "main" {
		body.AddChild((&Element{Tag: "script"}).AppendProp("src", "/assets/js/present.js").AddChild(RawText("")))
	}
	html.AddChild(body)
	return []byte("<!DOCTYPE html>" + html.Render()), nil
}
