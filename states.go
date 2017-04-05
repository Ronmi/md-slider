package main

import (
	"log"
	"regexp"
	"strings"
)

// 可能的狀態
const (
	SContent = iota
	SCodeBlock

	// 這幾個不會成為 current state
	SUL
	SOL
	SH1
	SH2
	SH3
)

const cuset = " \t\r\n" // 去頭尾用的

// regexps
var (
	rHeading *regexp.Regexp
	rUL      *regexp.Regexp
	rOL      *regexp.Regexp
)

func init() {
	f := func(title, regex string) *regexp.Regexp {
		ret, err := regexp.Compile(regex)
		if err != nil {
			log.Fatalf("Cannot compile %s regexp: %s", title, err)
		}

		return ret
	}
	rHeading = f("heading", `^#{1,3} [^#]*`)
	rUL = f("list", `^[*+-] [^*+-]*`)
	rOL = f("ordered list", `^[1-9][0-9]*\. `)
}

// line 代表一行解析後的資料
type line struct {
	indentSpaces int
	indentTabs   int
	indentLevel  int
	typ          int // 同上方的狀態
	content      string
}

// 取得 indent level，先比 tab 再比 space
func (l *line) getIndentLevel(spaces, tabs int) int {
	switch {
	case l.indentTabs > tabs:
		return l.indentLevel - 1
	case l.indentTabs < tabs:
		return l.indentLevel + 1
	case l.indentSpaces > spaces:
		return l.indentLevel - 1
	case l.indentSpaces < spaces:
		return l.indentLevel + 1
	}

	return l.indentLevel
}

type stateMachine struct {
	state    int
	lastLine *line
}

func newMachine() *stateMachine {
	return &stateMachine{
		state:    SContent,
		lastLine: &line{},
	}
}

func (m *stateMachine) parseIndent(leading string) (spaces, tabs int) {
	for _, c := range []rune(leading) {
		switch c {
		case ' ':
			spaces++
		case '\t':
			tabs++
		}
	}

	return
}

func (m *stateMachine) parse(txt string) *line {
	tmp := strings.TrimLeft(txt, " \t")
	leading := string(txt[:len(txt)-len(tmp)])
	real := strings.TrimRight(tmp, cuset)

	switch m.state {
	case SContent:
		m.lastLine = m.normalParse(leading, real)
		return m.lastLine
	case SCodeBlock:
		m.lastLine = m.codeBlockParse(leading, real)
		return m.lastLine
	}

	return &line{}
}

func (m *stateMachine) codeBlockParse(leading, txt string) *line {
	if txt != "```" {
		return &line{
			typ:     SCodeBlock,
			content: leading + txt,
		}
	}

	m.state = SContent
	return &line{
		typ:     SCodeBlock,
		content: leading,
	}
}

func (m *stateMachine) normalParse(leading, txt string) *line {
	switch {
	case strings.HasPrefix(txt, "```"):
		m.state = SCodeBlock
		// code block
		return &line{
			typ:     SCodeBlock,
			content: strings.TrimLeft(txt, "`"),
		}

	case rHeading.MatchString(txt):
		// heading
		arr := strings.SplitN(txt, " ", 2)
		return &line{
			typ:     SH1 + len(arr[0]) - 1,
			content: arr[1],
		}

	case rUL.MatchString(txt):
		m.state = SContent
		s, t := m.parseIndent(leading)
		// list
		return &line{
			indentSpaces: s,
			indentTabs:   t,
			indentLevel:  m.lastLine.getIndentLevel(s, t),
			typ:          SUL,
			content:      strings.TrimLeft(txt, "*-+ "),
		}
	case rOL.MatchString(txt):
		m.state = SContent
		s, t := m.parseIndent(leading)
		// list
		return &line{
			indentSpaces: s,
			indentTabs:   t,
			indentLevel:  m.lastLine.getIndentLevel(s, t),
			typ:          SOL,
			content:      strings.TrimLeft(txt, "0123456789. "),
		}
	}

	return &line{
		typ:     SContent,
		content: txt,
	}
}
