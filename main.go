package main

import (
	"fmt"
	"github.com/go-gui-org/go-gui/gui"
	"github.com/go-gui-org/go-gui/gui/backend"
)

const (
	MaxTextSize = 50
	MinTextSize = 5
)

type State struct {
	textSize float32
}

func main() {
	gui.Debug(true)
	gui.SetTheme(gui.ThemeDark.WithBorders(false))
	w := gui.NewWindow(gui.WindowCfg{
		State: &State{gui.CurrentTheme().N1.Size}, Title: "Subtitle_Visualizer", Width: 300, Height: 50, OnInit: func(w *gui.Window) {
			w.UpdateView(mainView)
		}, BgColor: gui.ColorTransparent, Decorations: gui.DecorationNone,
	})
	w.RegisterCommand(gui.Command{
		ID:       "Q",
		Label:    "Close",
		Shortcut: gui.Shortcut{Key: gui.KeyQ, Modifiers: gui.ModNone},
		Global:   true,
		Execute: func(_ *gui.Event, w *gui.Window) {
			fmt.Println("Hello from EventHandler")
			w.Close()
		},
	})
	w.RegisterCommand(gui.Command{
		ID:       "Esc",
		Label:    "Close",
		Shortcut: gui.Shortcut{Key: gui.KeyEscape, Modifiers: gui.ModNone},
		Global:   true,
		Execute: func(_ *gui.Event, w *gui.Window) {
			w.Close()
		},
	})
	w.RegisterCommand(gui.Command{
		ID:       "+",
		Label:    "IncreaseTextSize",
		Shortcut: gui.Shortcut{Key: gui.KeyUp, Modifiers: gui.ModNone},
		Global:   true,
		Execute: func(_ *gui.Event, w *gui.Window) {
			gui.State[State](w).textSize++
			//fmt.Printf("Increased Text Size to %v", gui.State[State](w).textSize)
		},
		CanExecute: func(w *gui.Window) bool {
			return gui.State[State](w).textSize < MaxTextSize
		},
	})
	w.RegisterCommand(gui.Command{
		ID:       "-",
		Label:    "DecreaseTextSize",
		Shortcut: gui.Shortcut{Key: gui.KeyDown, Modifiers: gui.ModNone},
		Global:   true,
		Execute: func(_ *gui.Event, w *gui.Window) {
			gui.State[State](w).textSize--
			//fmt.Printf("Decreased Text Size to %v", gui.State[State](w).textSize)
		},
		CanExecute: func(w *gui.Window) bool {
			return gui.State[State](w).textSize > MinTextSize
		},
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
	//fmt.Printf("TextSize is %v\n", gui.State[State](w))
	return gui.Column(gui.ContainerCfg{
		Sizing: gui.FillFill,
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
			gui.Label("Here go the subtitles", gui.TextStyle{Size: gui.State[State](w).textSize, Color: gui.CurrentTheme().N1.Color}),
		},
	})
}
