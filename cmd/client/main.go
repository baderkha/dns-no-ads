package main

import (
	"baderkha-no-dns/pkg/dns"
	"baderkha-no-dns/pkg/osproc"
	"fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/julien040/go-ternary"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("Dns AdBlocker"))
		window.Option(app.Size(unit.Dp(400), unit.Dp(600)))

		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops
	actionButtonState := &widget.Clickable{}

	serverIsUp := false
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			fmt.Println("destroying bye ...")
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			isBtnClicked := actionButtonState.Clicked(gtx)
			if isBtnClicked {
				if serverIsUp {
					go dns.Stop()
					serverIsUp = false
				} else {
					go dns.Start()

					serverIsUp = true
				}
			}

			if osproc.IsRoot() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(35),
								Left:   unit.Dp(35),
							}
							return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								status := material.H4(theme, ternary.If(serverIsUp, "Server Started ...", "Server Stopped"))
								status.Alignment = text.Middle
								// Define an large label with an appropriate text:
								return status.Layout(gtx)
							})

						},
					),
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(35),
								Left:   unit.Dp(35),
							}
							return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								status := material.H4(theme, fmt.Sprintf("Ad Hosts Blocked :[%d]", dns.AdsBlocked))
								status.Alignment = text.Middle
								// Define an large label with an appropriate text:
								return status.Layout(gtx)
							})

						},
					),
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(35),
								Left:   unit.Dp(35),
							}
							return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(theme, actionButtonState, ternary.If(serverIsUp, "Stop", "Start"))

								btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
								// Define an large label with an appropriate text:

								return btn.Layout(gtx)
							})

						},
					),
				)
			} else {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceEvenly,
				}.Layout(gtx,
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(35),
								Left:   unit.Dp(35),
							}
							return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								status := material.Body1(theme, "This application Requires admin/root Priveleges. Please Restart")
								status.Alignment = text.Middle
								// Define an large label with an appropriate text:
								return status.Layout(gtx)
							})
						},
					),
				)
			}

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
