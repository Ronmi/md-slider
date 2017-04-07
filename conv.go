package main

import (
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
	barr, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	root := &Element{Tag: "div"}
	curPage := 0
	var page *Element
	buf := ""

	save := func(curPage int, page *Element, buf string) (int, *Element) {
		if page != nil {
			buf = strings.TrimSpace(buf)
			if buf != "" {
				page.AddChild(MDText(buf))
			} else {
				page.Props[0].Value += " lonely"
			}
			root.AddChild(page)
		}

		curPage++
		return curPage, mkpage(curPage)
	}

	for _, l := range strings.Split(string(barr), "\n") {
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
		curPage, page = save(curPage, page, buf)
		buf = ""
		page.AddChild(&Element{
			Tag:     "h1",
			Content: []Renderer{Text(l)},
		})
	}

	save(curPage, page, buf)

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
		Content: []Renderer{RawText("var maxPage=" + strconv.Itoa(curPage) + ";")},
	}).AddChild(&Element{
		Tag:     "script",
		Props:   []Prop{Prop{Name: "src", Value: "/assets/before.js"}},
		Content: []Renderer{RawText("")},
	}).AddChild(&Element{
		Tag:     "script",
		Props:   []Prop{Prop{Name: "src", Value: "/assets/prism.js"}},
		Content: []Renderer{RawText("")},
	}).AddChild(&Element{
		Tag:     "script",
		Props:   []Prop{Prop{Name: "src", Value: "/assets/after.js"}},
		Content: []Renderer{RawText("")},
	}))
	return []byte(html.Render()), nil
}
