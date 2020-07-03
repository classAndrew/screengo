package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

const offset int = 90
const scale float64 = 1.15762 // this value with vary across machines, I don't know how fyne handles scaling of images

type tappableIcon struct {
	widget.Icon
	win *fyne.Window
	c   int
	p   [][]int
}

func (t *tappableIcon) Tapped(f *fyne.PointEvent) {
	t.p[t.c] = []int{int(float64((f.Position.X - offset)) * scale), int(float64(f.Position.Y) * scale)}
	fmt.Printf("%v\n", t.p)
	t.c++
	if t.c == 2 {
		go uploadProcess(t)
		(*(t.win)).Hide()
	}
}

func (t *tappableIcon) TappedSecondary(_ *fyne.PointEvent) {
}

func newTappableIcon(res fyne.Resource, win *fyne.Window) *tappableIcon {
	icon := &tappableIcon{}
	icon.win = win
	icon.p = [][]int{[]int{0, 0}, []int{0, 0}}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	return icon
}

func uploadProcess(t *tappableIcon) {
	f, _ := os.Open("screenshot.png")
	im1, _, _ := image.Decode(f)
	x, y := math.Abs(float64(t.p[1][0]-t.p[0][0])), math.Abs(float64(t.p[1][1]-t.p[0][1]))
	m := image.NewRGBA(image.Rect(0, 0, int(x), int(y)))
	draw.Draw(m, m.Bounds(), im1, image.Point{t.p[0][0], t.p[0][1]}, draw.Src)
	to, _ := os.Create("new.png")
	png.Encode(to, m)
	f.Close()
	to.Close()
	res := UploadImgur(EncodeB64("new.png"))
	del := FindDelHash(res)
	id := FindID(res)
	logf, _ := os.OpenFile("logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	logf.WriteString(id + "-" + del + "\n")
	// fmt.Printf("%s\n%s\n%s\n", res, del, id)
	w := fyne.CurrentApp().NewWindow("Output")
	entry := widget.NewEntry()
	entry.Text = "https://i.imgur.com/" + id + ".png"
	w.SetContent(widget.NewVBox(
		entry,
	))
	w.Show()
	(*(t.win)).Close()
}

var appb fyne.App

func main() {
	appb = app.New()
	w := appb.NewWindow("Grab")
	sw := fyne.CurrentApp().NewWindow("Crop")
	sw.Resize(fyne.NewSize(1366, 768))
	texten := widget.NewEntry()
	texten.Text = "1"
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Capture!"),
		widget.NewHBox(
			texten,
			widget.NewButton("Start", func() {
				w.Hide()
				i, _ := strconv.ParseInt(texten.Text, 10, 64)
				Screenshot(int(i))
				s, _ := fyne.LoadResourceFromPath("screenshot.png")
				sw.SetContent(newTappableIcon(s, &sw))
				sw.Show()
				w.Close()
			})),
	))
	w.ShowAndRun()

	fmt.Println("Screenshot Loaded")

}
