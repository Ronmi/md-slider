package main

import "testing"

func TestRegex(t *testing.T) {
	testData := []struct {
		data    string
		heading bool
		ul      bool
		ol      bool
	}{
		{"abc", false, false, false},
		{"#abc", false, false, false},
		{"##abc", false, false, false},
		{"###abc", false, false, false},
		{"####abc", false, false, false},
		{"# abc", true, false, false},
		{"## abc", true, false, false},
		{"### abc", true, false, false},
		{"#### abc", false, false, false},
		{"*abc", false, false, false},
		{"*abc*", false, false, false},
		{"* abc", false, true, false},
		{"* abc*", false, true, false},
		{"+abc", false, false, false},
		{"+abc+", false, false, false},
		{"+ abc", false, true, false},
		{"+ abc+", false, true, false},
		{"-abc", false, false, false},
		{"-abc-", false, false, false},
		{"- abc", false, true, false},
		{"- abc-", false, true, false},
		{"0.abc", false, false, false},
		{"0. abc", false, false, false},
		{"1.abc", false, false, false},
		{"1. abc", false, false, true},
		{"19.abc", false, false, false},
		{"19. abc", false, false, true},
	}
	conv := func(b bool) string {
		if b {
			return "pass"
		}
		return "fail"
	}

	for _, v := range testData {
		if res := rHeading.MatchString(v.data); res != v.heading {
			t.Errorf(
				"Matching %s against %s (heading): expected to %s, but %sed",
				v.data,
				rHeading,
				conv(v.heading),
				conv(res),
			)
		}
		if res := rUL.MatchString(v.data); res != v.ul {
			t.Errorf(
				"Matching %s against %s (ul): expected to %s, but %sed",
				v.data,
				rUL,
				conv(v.ul),
				conv(res),
			)
		}
		if res := rOL.MatchString(v.data); res != v.ol {
			t.Errorf(
				"Matching %s against %s (ol): expected to %s, but %sed",
				v.data,
				rOL,
				conv(v.ol),
				conv(res),
			)
		}
	}
}
