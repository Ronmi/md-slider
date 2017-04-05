package main

import "fmt"

func mkcodepen(html, css, js string) *Element {
	return (&Element{
		Tag: "form",
		Props: []Prop{
			Prop{Name: "method", Value: "POST"},
			Prop{Name: "target", Value: "codepen"},
			Prop{Name: "action", Value: "http://codepen.io/pen/define"},
		},
	}).AddChild(&Element{
		Tag: "input",
		Props: []Prop{
			Prop{Name: "type", Value: "hidden"},
			Prop{Name: "name", Value: "data"},
			Prop{Name: "value", Value: fmt.Sprintf(
				`{&#34;html&#34;:&#34;%s&#34;,&#34;css&#34;:&#34;%s&#34;,&#34;js&#34;:&#34;%s&#34;}`,
				html, css, js,
			)},
		},
	}).AddChild(&Element{
		Tag: "button",
		Props: []Prop{
			Prop{Name: "type", Value: "submit"},
		},
		Content: []Renderer{RawText("try it!")},
	})
}
