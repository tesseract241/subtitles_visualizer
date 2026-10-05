package main

import (
	"bufio"
	"fmt"
	"github.com/go-gui-org/go-gui/gui"
	"github.com/go-gui-org/go-gui/gui/backend"
	"os"
	"strings"
	"unicode"
)

const (
	InitialTextSize = 50
	MaxTextSize     = 150
	MinTextSize     = 5
)

type State struct {
	textList TextList
	text string
}

type TextList struct {
	b [2]strings.Builder
	f int
	c int
}


func (l *TextList) Append(s string) {
	l.b[l.c].WriteString(s)
}

func (l *TextList) NewLine(s string) {
	l.f = l.c
	l.c ^= 1
	l.b[int(l.c)].Reset()
	l.b[int(l.c)].WriteString(s)
}

func (l *TextList) ToText() string {
	return fmt.Sprintf("%s\n%s", l.b[l.f].String(), l.b[l.f^1].String())
}

func _scanner_to_channel(in *bufio.Scanner) chan string {
	c := make(chan string, 20)
	go func(){
		for in.Scan() {
			c <- in.Text()
		}
		if err := in.Err(); err!=nil {
			fmt.Fprintln(os.Stderr, "Error ", err)
		}
		close(c)
	}()
	return c
}

func parse(t string, l *TextList, out *bufio.Writer) {
	fmt.Fprint(out, t)
	out.Flush()
	index := strings.IndexFunc(t, unicode.IsUpper)
	if index > 0 {
		l.Append(t[:index])
	}
	if index != -1 {
		l.NewLine(t[index:])
	}
}

func main() {
	gui.Debug(true)
	gui.DebugCategories(gui.DebugAll | gui.DebugUnscopedIDs)
	state := &State{text : "Initial Text"}
	t, e := gui.ThemeDark.WithBorders(false).AdjustFontSize(InitialTextSize, MinTextSize, MaxTextSize)
	if e != nil {
		fmt.Fprintln(os.Stderr, "Error changing font size: ", e)
			return 
	}
	gui.SetTheme(t)
	in  := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	textList := new(TextList)
	c := _scanner_to_channel(in)
	w := gui.NewWindow(gui.WindowCfg{
		State: state, Title: "Subtitle_Visualizer", Width: 800, Height: 100,
		BgColor: gui.RGBA(0, 0, 0, 10), Decorations: gui.DecorationNone, Transparent: true,
	})
	w.RegisterCommand(gui.Command{
		ID:       "Q",
		Label:    "Close",
		Shortcut: gui.Shortcut{Key: gui.KeyQ, Modifiers: gui.ModNone},
		Global:   true,
		Execute: func(_ *gui.Event, w *gui.Window) {
			os.Stdin.Write([]byte{0})
			os.Stdin.Close()
			w.Close()
		},
	})
	w.RegisterCommand(gui.Command{
		ID:       "Esc",
		Label:    "Close",
		Shortcut: gui.Shortcut{Key: gui.KeyEscape, Modifiers: gui.ModNone},
		Global:   true,
		Execute: func(_ *gui.Event, w *gui.Window) {
			os.Stdin.Close()
			w.Close()
		},
	})
	w.RegisterCommand(gui.Command{
		ID:       "+",
		Label:    "IncreaseTextSize",
		Shortcut: gui.Shortcut{Key: gui.KeyUp, Modifiers: gui.ModNone},
		Global:   true,
		Execute: func(_ *gui.Event, w *gui.Window) {
			t, e := w.Theme().AdjustFontSize(1, MinTextSize, MaxTextSize)
			if e != nil {
				fmt.Fprintln(os.Stderr, "Error changing font size: ", e)
				return 
			}
			w.SetTheme(t)
			//fmt.Printf("Increased Text Size to %v", state.textSize)
		},
	})
	w.RegisterCommand(gui.Command{
		ID:       "-",
		Label:    "DecreaseTextSize",
		Shortcut: gui.Shortcut{Key: gui.KeyDown, Modifiers: gui.ModNone},
		Global:   true,
		Execute: func(_ *gui.Event, w *gui.Window) {
			t, e := w.Theme().AdjustFontSize(-1, MinTextSize, MaxTextSize)
			if e != nil {
				fmt.Fprintln(os.Stderr, "Error changing font size: ", e)
				return 
			}
			w.SetTheme(t)
			//fmt.Printf("Decreased Text Size to %v", state.textSize)
		},
	})
	w.SetView(mainView)
	gui.Stream(w, c, func(w *gui.Window, t string) {
		parse(t, textList, out)
		state.text = textList.ToText()
	})

	//w.OnEvent = func(e *gui.Event, w *gui.Window) {
	//switch e.Type {
	//	case gui.EventKeyDown:
	//		fmt.Printf("Key pressed: %v\n", e.KeyCode)
	//	}
	//}
	backend.Run(w)
}

func mainView(w *gui.Window) gui.View {
	state := gui.State[State](w)
	return gui.Column(gui.ContainerCfg{
		ID    : "text_column",
		Sizing: gui.FillFit,
		HAlign: gui.HAlignCenter,
		VAlign: gui.VAlignMiddle,
		OnMouseDown: func(ctx gui.EventCtx) {
			switch ctx.Event.MouseButton {
			case gui.MouseLeft:
				ctx.Window.StartWindowDrag()
			case gui.MouseRight:
				ctx.Window.StartWindowResize(gui.EdgeBottomRight)
			default:
			}
			ctx.Consume()
		},
		Content: []gui.View{
			gui.Text(gui.TextCfg{ID : "Sub_Text", Text : state.text, Mode : gui.TextModeWrap}),
		},
	})
}
