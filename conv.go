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

func conv(fn string) (string, error) {
	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
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
		head.
			AddChild(&Element{Tag: "meta", Props: []Prop{Prop{Name: "charset", Value: "utf8"}}}).
			AddChild(&Element{
				Tag: "link",
				Props: []Prop{
					Prop{Name: "rel", Value: "stylesheet"},
					Prop{Name: "href", Value: "https://ronmi.tw/prism.all.min.css"},
				},
			}).AddChild(&Element{Tag: "style", Content: []Renderer{RawText(`
body {
  background-color: #cccccc;
}
form {
  margin: 0;
}
button {
  background-color: white;
  z-index: 999999;
  border: 1px solid #9a9a9a;
  border-radius: 5px;
  padding: 3px 10px;
  margin-top: -2em;
  margin-left: 1em;
  position: relative;
}
button:hover {
  background-color: #eeeeee;
}
div.page {
  background-color: white;
  box-shadow: 10px 10px 5px grey;
  width: 80%;
  height: 80%;
  padding: 1em 1.8em;
  margin: 5% auto;
  border: 1px solid white;
  border-radius: 1em;
  font-size: x-large;
}
li,p,span,b,i,strike {
  font-size: x-large;
}
div.lonely {
  vertical-align: middle;
}
div.lonely > h1 {
  font-size: 3em;
}
h1 {
  font-size: 2em;
  text-align: center;
  margin: 0;
}
pre[class*="language-"] {
  width: 80%;
  max-height: 50%;
  overflow: auto;
}`)}}),
	)
	body := &Element{Tag: "body", Props: []Prop{Prop{Name: "id", Value: "body"}}}
	html.AddChild(body.AddChild(root).AddChild(&Element{
		Tag:     "script",
		Props:   []Prop{Prop{Name: "src", Value: "https://ronmi.tw/prism.all.min.js"}},
		Content: []Renderer{RawText("")},
	}).AddChild(&Element{
		Tag: "script",
		Content: []Renderer{RawText("var maxPage=" + strconv.Itoa(curPage-1) + `
var cur = 1;
document.getElementById("body").addEventListener("keypress", function(e){
  var x = e.which || e.keyCode;
  if (x != 37 && x != 39) return;

  if (x == 37) {
    // left
    cur--;
    if (cur < 1) cur = 1;
  } else if (x == 39) {
    // right
    cur++;
    if (cur > maxPage) cur = maxPage;
  }

  e = document.getElementById("page" + cur + "");
  e.scrollIntoView();
  var style = e.currentStyle || window.getComputedStyle(e);
  window.scrollBy(0, -parseInt(style.marginTop)/2);
});
`)},
	}))
	return html.Render(), nil
}
