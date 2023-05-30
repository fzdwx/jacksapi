package main

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	ai "github.com/fzdwx/jacksapi"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"math/rand"
	"os"
	"strings"
	"time"
)

func ask() {
	if len(os.Args) == 1 {
		fmt.Println("Please input your ask")
		os.Exit(1)
	}

	var (
		content = strings.Join(os.Args[1:], " ")
	)
	stream := client.ChatStream(
		[]ai.ChatMessage{
			{Role: "system", Content: "Format the response as Markdown."},
			{Role: "user", Content: content},
		}).Stream(false)

	m := &model{chatStream: stream}
	p := tea.NewProgram(m)
	m.p = p
	_, err := p.Run()
	if err != nil {
		panic(err)
	}
}

type model struct {
	chatStream *ai.ChatStream

	viewport viewport.Model
	spinner  tea.Model
	p        *tea.Program

	content *bytes.Buffer

	showSpinner bool // Whether or not to show the spinner
	w           int
	h           int
	err         error
}

type start struct{}
type stop struct{ err error }
type chatContent struct{ r rune }
type setViewportContent struct{ content string }

func (m *model) Init() tea.Cmd {
	m.viewport = viewport.New(0, 0)
	m.content = bytes.NewBufferString("")
	m.spinner = newCyclingChars(10, "Generating", lipgloss.NewRenderer(os.Stderr, termenv.WithColorCache(true)))

	return tea.Batch(func() tea.Msg {
		return start{}
	}, m.spinner.Init())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m.w = msg.Width
		m.h = msg.Height
		cmds = append(cmds, m.flush())
	case start:
		m.showSpinner = true
		return m, m.start()
	case stop:
		m.showSpinner = false
		m.err = msg.err
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case chatContent:
		m.showSpinner = false
		m.content.WriteRune(msg.r)
		return m, m.flush()
	case setViewportContent:
		m.viewport.SetContent(msg.content)
		m.viewport.GotoBottom()
	}

	if m.showSpinner {
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
	}

	var viewportCmd tea.Cmd
	m.viewport, viewportCmd = m.viewport.Update(msg)
	cmds = append(cmds, viewportCmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s", m.err)
	}

	if m.showSpinner {
		return lipgloss.NewStyle().Width(m.w).Height(m.h).Align(lipgloss.Center, lipgloss.Center).Render(m.spinner.View())
	}

	return m.viewport.View()
}

func (m *model) start() tea.Cmd {
	return func() tea.Msg {
		go m.chatStream.DoWithCallback(ai.With(func(r rune, done bool, err error) {
			if err != nil || done {
				m.p.Send(stop{err})
				return
			}

			m.p.Send(chatContent{r})
		}))
		return nil
	}
}

func (m *model) flush() tea.Cmd {
	return func() tea.Msg {
		go func() {
			content, _ := newRender(m.viewport.Width).Render(m.content.String())
			m.p.Send(setViewportContent{content})
		}()
		return nil
	}
}

func newRender(wordWrap int) *glamour.TermRenderer {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(wordWrap),
		glamour.WithEmoji(),
	)
	return renderer
}

/*
inspired by https://github.com/charmbracelet/mods/blob/main/cyclingchars.go

# MIT License

# Copyright (c) 2023 Charmbracelet, Inc

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
const (
	charCyclingFPS  = time.Second / 22
	maxCyclingChars = 120
)

var (
	charRunes = []rune("0123456789abcdefABCDEF~!@#$£€%^&*()+=_")

	ellipsisSpinner = spinner.Spinner{
		Frames: []string{"", ".", "..", "..."},
		FPS:    time.Second / 3, //nolint:gomnd
	}
)

type charState int

const (
	charInitialState charState = iota
	charCyclingState
	charEndOfLifeState
)

// cyclingChar is a single animated character.
type cyclingChar struct {
	finalValue   rune // if < 0 cycle forever
	currentValue rune
	initialDelay time.Duration
	lifetime     time.Duration
}

func (c cyclingChar) randomRune() rune {
	return (charRunes)[rand.Intn(len(charRunes))] //nolint:gosec
}

func (c cyclingChar) state(start time.Time) charState {
	now := time.Now()
	if now.Before(start.Add(c.initialDelay)) {
		return charInitialState
	}
	if c.finalValue > 0 && now.After(start.Add(c.initialDelay)) {
		return charEndOfLifeState
	}
	return charCyclingState
}

type stepCharsMsg struct{}

func stepChars() tea.Cmd {
	return tea.Tick(charCyclingFPS, func(_ time.Time) tea.Msg {
		return stepCharsMsg{}
	})
}

// cyclingChars is the model that manages the animation that displays while the
// output is being generated.
type cyclingChars struct {
	start           time.Time
	chars           []cyclingChar
	ramp            []lipgloss.Style
	label           []rune
	ellipsis        spinner.Model
	ellipsisStarted bool
	styles          lipgloss.Style
}

func newCyclingChars(initialCharsSize uint, label string, r *lipgloss.Renderer) cyclingChars {
	n := int(initialCharsSize)
	if n > maxCyclingChars {
		n = maxCyclingChars
	}

	gap := " "
	if n == 0 {
		gap = ""
	}

	c := cyclingChars{
		start:    time.Now(),
		label:    []rune(gap + label),
		ellipsis: spinner.New(spinner.WithSpinner(ellipsisSpinner)),
		styles:   r.NewStyle().Foreground(lipgloss.Color("#FF87D7")),
	}

	// If we're in truecolor mode (and there are enough cycling characters)
	// color the cycling characters with a gradient ramp.
	const minRampSize = 3
	if n >= minRampSize && r.ColorProfile() == termenv.TrueColor {
		c.ramp = make([]lipgloss.Style, n)
		ramp := makeGradientRamp(n)
		for i, color := range ramp {
			c.ramp[i] = r.NewStyle().Foreground(color)
		}
	}

	makeDelay := func(a int32, b time.Duration) time.Duration {
		return time.Duration(rand.Int31n(a)) * (time.Millisecond * b) //nolint:gosec
	}

	makeInitialDelay := func() time.Duration {
		return makeDelay(8, 60) //nolint:gomnd
	}

	c.chars = make([]cyclingChar, n+len(c.label))

	// Initial characters that cycle forever.
	for i := 0; i < n; i++ {
		c.chars[i] = cyclingChar{
			finalValue:   -1, // cycle forever
			initialDelay: makeInitialDelay(),
		}
	}

	// Label text that only cycles for a little while.
	for i, r := range c.label {
		c.chars[i+n] = cyclingChar{
			finalValue:   r,
			initialDelay: makeInitialDelay(),
			lifetime:     makeDelay(5, 180), //nolint:gomnd
		}
	}

	return c
}

// Init initializes the animation.
func (c cyclingChars) Init() tea.Cmd {
	return stepChars()
}

// Update handles messages.
func (c cyclingChars) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case stepCharsMsg:
		for i, char := range c.chars {
			switch char.state(c.start) {
			case charInitialState:
				c.chars[i].currentValue = '.'
			case charCyclingState:
				c.chars[i].currentValue = char.randomRune()
			case charEndOfLifeState:
				c.chars[i].currentValue = char.finalValue
			}
		}

		if !c.ellipsisStarted {
			var eol int
			for _, char := range c.chars {
				if char.state(c.start) == charEndOfLifeState {
					eol++
				}
			}
			if eol == len(c.label) {
				// If our entire label has reached end of life, start the
				// ellipsis "spinner" after a short pause.
				c.ellipsisStarted = true
				cmd = tea.Tick(time.Millisecond*220, func(_ time.Time) tea.Msg { //nolint:gomnd
					return c.ellipsis.Tick()
				})
			}
		}

		return c, tea.Batch(stepChars(), cmd)
	case spinner.TickMsg:
		var cmd tea.Cmd
		c.ellipsis, cmd = c.ellipsis.Update(msg)
		return c, cmd
	default:
		return c, nil
	}
}

// View renders the animation.
func (c cyclingChars) View() string {
	var b strings.Builder
	for i, char := range c.chars {
		var (
			s *lipgloss.Style
			r = char.currentValue
		)
		if len(c.ramp) > 0 && i < len(c.ramp) {
			// There's a gradient ramp style defined for this char. Style it
			// accordingly.
			s = &c.ramp[i]
		} else if char.finalValue < 0 {
			// No gradient ramp defined, but this color will cycle forever so
			// let's color it accordingly.
			s = &c.styles
		}
		if s != nil {
			b.WriteString(s.Render(string(r)))
			continue
		}
		b.WriteRune(r)
	}
	return b.String() + c.ellipsis.View()
}

func makeGradientRamp(length int) []lipgloss.Color {
	const startColor = "#F967DC"
	const endColor = "#6B50FF"
	var (
		c        = make([]lipgloss.Color, length)
		start, _ = colorful.Hex(startColor)
		end, _   = colorful.Hex(endColor)
	)
	for i := 0; i < length; i++ {
		step := start.BlendLuv(end, float64(i)/float64(length))
		c[i] = lipgloss.Color(step.Hex())
	}
	return c
}

func makeGradientText(baseStyle lipgloss.Style, str string) string {
	const minSize = 3
	if len(str) < minSize {
		return str
	}
	b := strings.Builder{}
	runes := []rune(str)
	for i, c := range makeGradientRamp(len(str)) {
		b.WriteString(baseStyle.Copy().Foreground(c).Render(string(runes[i])))
	}
	return b.String()
}
