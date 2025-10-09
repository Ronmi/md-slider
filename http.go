package main

import (
	"embed"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
)

//go:embed assets
var assetsFS embed.FS

type httpHandler struct {
	devMode bool
}

func (h httpHandler) loadAsset(fn string) ([]byte, error) {
	p := "assets/" + fn
	if !h.devMode {
		return assetsFS.ReadFile(p)
	}

	if _, err := os.Stat("." + p); err != nil {
		return assetsFS.ReadFile(p)
	}

	return os.ReadFile("." + p)
}

func (h httpHandler) setMIME(w http.ResponseWriter, fn string) {
	arr := strings.Split(fn, ".")
	if l := len(arr); l > 0 {
		if typ := mime.TypeByExtension("." + arr[len(arr)-1]); typ != "" {
			w.Header().Set("Content-Type", typ)
		}
	}
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir := "." + strings.TrimRight(r.URL.Path, "/") // strip leading and tailing slash

	// 先看看是不是 assets
	if idx := strings.Index(dir, "/assets/"); idx > 0 {

		ret, err := h.loadAsset(dir[idx+8:])
		if err != nil {
			w.WriteHeader(404)
			return
		}

		h.setMIME(w, dir)
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
			w.WriteHeader(403)
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

		ret, err := h.loadAsset("list.html")
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

	name := "main"
	if !strings.HasSuffix(info.Name(), ".md") {
		f, err := os.Open(dir)
		if err != nil {
			w.WriteHeader(403)
			return
		}
		defer f.Close()

		io.Copy(w, f)
		return
	}

	if r.URL.Query().Get("present") != "" {
		name = "present"
	}

	ret, err := conv(dir, name)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	w.Write(ret)
}
