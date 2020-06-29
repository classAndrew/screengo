package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

const offset int = 90
const scale float64 = 1.15762 // this value with vary across machines, I don't know how fyne handles scaling of images

type tappableIcon struct {
	widget.Icon
	c int
	p [][]int
}

func (t *tappableIcon) Tapped(f *fyne.PointEvent) {
	t.p[t.c] = []int{int(float64((f.Position.X - offset)) * scale), int(float64(f.Position.Y) * scale)}
	fmt.Printf("%v\n", t.p)
	t.c++
	if t.c == 2 {
		f, _ := os.Open("screenshot.png")
		im1, _, _ := image.Decode(f)
		x, y := math.Abs(float64(t.p[1][0]-t.p[0][0])), math.Abs(float64(t.p[1][1]-t.p[0][1]))
		m := image.NewRGBA(image.Rect(0, 0, int(x), int(y)))
		draw.Draw(m, m.Bounds(), im1, image.Point{t.p[0][0], t.p[0][1]}, draw.Src)
		to, _ := os.Create("new.png")
		png.Encode(to, m)
		f.Close()
		to.Close()
		os.Exit(0)
	}
}

func (t *tappableIcon) TappedSecondary(_ *fyne.PointEvent) {
}

func newTappableIcon(res fyne.Resource) *tappableIcon {
	icon := &tappableIcon{}
	icon.p = [][]int{[]int{0, 0}, []int{0, 0}}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	return icon
}

func main() {
	appb := app.New()
	w := appb.NewWindow("Grab")
	sw := fyne.CurrentApp().NewWindow("Crop")

	sw.Resize(fyne.NewSize(1366, 768))
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Capture!"),
		widget.NewButton("Start", func() {
			w.Hide()
			Screenshot(1)
			s, _ := fyne.LoadResourceFromPath("screenshot.png")
			sw.SetContent(newTappableIcon(s))
			sw.Show()
			w.Close()
		}),
	))
	w.ShowAndRun()

	fmt.Println("Screenshot Loaded")

}
