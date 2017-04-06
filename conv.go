package main

import (
	"html"
	"io/ioutil"
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

func conv(fn string) ([]byte, error) {
	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	root := &Element{
		Tag:     "div",
		Content: []Renderer{},
	}

	m := newMachine()

	var (
		page      *Element
		paragraph *Element
		code      *Element
		ol        *Element
		ul        *Element
	)

	save := func(p, child *Element) *Element {
		if child != nil && len(child.Content) > 0 {
			p.AddChild(child)
		}

		return nil
	}
	savecode := func(p *Element) {
		if code == nil {
			return
		}

		// 最後一行如果是空白就砍了
		l := len(code.Content)
		last := code.Content[l-1].Render()
		if strings.TrimSpace(last) == "" {
			code.Content = code.Content[0 : l-1]
		}

		lang := code.Props[0].Value
		pre := &Element{
			Tag:     "pre",
			Props:   []Prop{Prop{Name: "class", Value: "line-numbers"}},
			Content: []Renderer{code},
		}
		p.AddChild(pre)

		var button *Element
		// 各種 code blocks 特殊處理
		switch lang {
		case "language-html":
			button = mkcodepen(code.GetContent(), "", "")
		case "language-css":
			button = mkcodepen(html.EscapeString("<html><head></head><body><h1>heading</h1><p>some text</p></body></html>"), code.GetContent(), "")
		case "language-js":
			button = mkcodepen("", "", code.GetContent())
		}

		if button != nil {
			p.AddChild(button)
		}

		code = nil
	}

	saveall := func(p *Element) {
		paragraph = save(p, paragraph)
		savecode(p)
		ol = save(p, ol)
		ul = save(p, ul)
	}

	curPage := 1

	for _, txt := range strings.Split(string(buf), "\n") {
		l := m.parse(txt)

		switch l.typ {
		case SH1, SH2:
			if page != nil {
				// 特殊處理只有標題的 page
				if len(page.Content) == 1 {
					page.Props[0].Value += " lonely"
				}
				saveall(page)
				root.AddChild(page)
			}
			page = mkpage(curPage)
			curPage++
			page.AddChild(&Element{
				Tag:     "h1",
				Content: []Renderer{Text(l.content)},
			})
		case SH3:
			saveall(page)
			page.AddChild(&Element{
				Tag:     "h2",
				Content: []Renderer{Text(l.content)},
			})
		case SUL:
			paragraph = save(page, paragraph)
			savecode(page)
			ol = save(page, ol)
			if ul == nil {
				ul = &Element{
					Tag: "ul",
					Props: []Prop{
						Prop{Name: "indentLevel", Value: strconv.Itoa(l.indentLevel)},
					},
				}
			}
			ul.AddChild(&Element{
				Tag:     "li",
				Content: []Renderer{MDText(l.content)},
			})
		case SOL:
			paragraph = save(page, paragraph)
			savecode(page)
			ul = save(page, ul)
			if ol == nil {
				ol = &Element{
					Tag: "ol",
					Props: []Prop{
						Prop{Name: "indentLevel", Value: strconv.Itoa(l.indentLevel)},
					},
					Content: []Renderer{},
				}
			}
			ol.AddChild(&Element{
				Tag:     "li",
				Content: []Renderer{MDText(l.content)},
			})
		case SCodeBlock:
			paragraph = save(page, paragraph)
			ol = save(page, ol)
			ul = save(page, ul)

			if code == nil {
				code = &Element{
					Tag: "code",
					Props: []Prop{
						Prop{Name: "class", Value: "language-" + l.content},
					},
					Content: []Renderer{},
				}
				break
			}

			code.AddChild(Text(l.content))
		case SContent:
			savecode(page)
			ol = save(page, ol)
			ul = save(page, ul)

			if l.content == "" {
				save(page, paragraph)
				paragraph = &Element{Tag: "p"}
				break
			}

			if paragraph == nil {
				page.AddChild(MDText(l.content))
				break
			}

			paragraph.AddChild(MDText(l.content))
		}
	}

	if page != nil {
		saveall(page)
		root.AddChild(page)
	}

	html := &Element{Tag: "html"}
	head := &Element{Tag: "head"}
	html.AddChild(
		head.AddChild(&Element{
			Tag:   "meta",
			Props: []Prop{Prop{Name: "charset", Value: "utf8"}},
		}).AddChild(&Element{
			Tag: "link",
			Props: []Prop{
				Prop{Name: "rel", Value: "stylesheet"},
				Prop{Name: "href", Value: "/assets/prism.css"},
			},
		}).AddChild(&Element{
			Tag: "link",
			Props: []Prop{
				Prop{Name: "rel", Value: "stylesheet"},
				Prop{Name: "href", Value: "/assets/main.css"},
			},
		}),
	)
	body := &Element{Tag: "body", Props: []Prop{Prop{Name: "id", Value: "body"}}}
	html.AddChild(body.AddChild(root).AddChild(&Element{
		Tag:     "script",
		Props:   []Prop{Prop{Name: "src", Value: "/assets/prism.js"}},
		Content: []Renderer{RawText("")},
	}).AddChild(&Element{
		Tag: "script",
		Content: []Renderer{RawText("var maxPage=" + strconv.Itoa(curPage-1) + `
`)},
	}).AddChild(&Element{
		Tag:     "script",
		Props:   []Prop{Prop{Name: "src", Value: "/assets/event.js"}},
		Content: []Renderer{RawText("")},
	}))
	return []byte(html.Render()), nil
}
