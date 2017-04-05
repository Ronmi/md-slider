package main

import "testing"

func TestHTMLProp(t *testing.T) {
	testData := []struct {
		p      Prop
		expect string
	}{
		{Prop{"prop1", "val1"}, `prop1="val1"`},
		{Prop{"prop1", `"val1"`}, `prop1="\"val1\""`},
		{Prop{"style", "color:#ffffff;border:1px solid black;"}, `style="color:#ffffff;border:1px solid black;"`},
	}

	for _, data := range testData {
		actual := data.p.Render()
		if actual != data.expect {
			t.Errorf("Rendering property, expected %s, got %s", data.expect, actual)
		}
	}
}

func TestHTMLElement(t *testing.T) {
	testData := []struct {
		e      Element
		expect string
	}{
		{
			e: Element{
				Tag: "a",
				Props: []Prop{
					Prop{"href", "https://google.com"},
				},
				Content: []Renderer{
					Text("Google"),
				},
			},
			expect: `<a href="https://google.com">Google</a>`,
		},
		{
			e: Element{
				Tag: "code",
				Content: []Renderer{
					Text("Google"),
					Text("Bing"),
					Text("Yahoo"),
				},
			},
			expect: `<code>Google
Bing
Yahoo</code>`,
		},
		{
			e: Element{
				Tag: "a",
				Props: []Prop{
					Prop{"href", "https://google.com"},
				},
				Content: []Renderer{
					Text("All hail to "),
					&Element{
						Tag:     "bold",
						Content: []Renderer{Text("Google")},
					},
				},
			},
			expect: `<a href="https://google.com">All hail to 
<bold>Google</bold></a>`,
		},
	}

	for _, data := range testData {
		actual := data.e.Render()
		if actual != data.expect {
			t.Errorf("Rendering element, expected %s, got %s", data.expect, actual)
		}
	}
}
