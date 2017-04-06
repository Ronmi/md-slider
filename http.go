//go:generate go-bindata assets/

package main

import (
	"net/http"
	"os"
	"strings"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	dir := "." + strings.TrimRight(r.URL.Path, "/") // strip leading and tailing slash
	arr := strings.Split(dir, "/")
	l := len(arr)

	// 先看看是不是 assets
	if l >= 2 && arr[l-2] == "assets" {
		ret, err := Asset("assets/" + arr[l-1])
		if err != nil {
			w.WriteHeader(404)
			return
		}

		w.Write(ret)
		return
	}

	info, err := os.Stat(dir)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if info.IsDir() {
		f, err := os.Open(dir)
		if err != nil {
			w.WriteHeader(401)
			return
		}
		infos, _ := f.Readdir(-1)

		ul := &Element{Tag: "ul"}
		for _, i := range infos {
			if strings.HasPrefix(i.Name(), ".") {
				continue
			}
			if !i.IsDir() && !strings.HasSuffix(i.Name(), ".md") {
				continue
			}
			ul.AddChild(
				(&Element{Tag: "li"}).AddChild(&Element{
					Tag:     "a",
					Props:   []Prop{{Name: "href", Value: "/" + dir + "/" + i.Name()}},
					Content: []Renderer{Text(i.Name())},
				}),
			)
		}

		ret := (&Element{Tag: "html"}).AddChild(
			(&Element{Tag: "head"}).AddChild(&Element{
				Tag:   "meta",
				Props: []Prop{{Name: "charset", Value: "utf8"}},
			}),
		).AddChild((&Element{Tag: "body"}).AddChild(ul))

		w.Write([]byte(ret.Render()))
		return
	}

	if !strings.HasSuffix(info.Name(), ".md") {
		w.WriteHeader(400)
		return
	}

	ret, err := conv(dir)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	w.Write(ret)
}
