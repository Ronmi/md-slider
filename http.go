//go:generate go-bindata assets/

package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func loadAsset(fn string) ([]byte, error) {
	p := "assets/" + fn
	if _, err := os.Stat("." + p); err != nil {
		return Asset(p)
	}

	return ioutil.ReadFile("." + p)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	dir := "." + strings.TrimRight(r.URL.Path, "/") // strip leading and tailing slash
	arr := strings.Split(dir, "/")
	l := len(arr)

	// 先看看是不是 assets
	if l >= 2 && arr[l-2] == "assets" {

		ret, err := loadAsset(arr[l-1])
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

		ret, err := loadAsset("list.html")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		dir = dir[1:]
		if dir == "" {
			dir = "root"
		}
		str := strings.Replace(string(ret), "{{list}}", ul.Render(), -1)
		str = strings.Replace(str, "{{path}}", dir, -1)
		w.Write([]byte(str))
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
