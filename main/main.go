package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// func main() {
// 	coords := []onmap.Coord{
// 		{42.1, 19.1},             // Bar
// 		{55.755833, 37.617222},   // Moscow
// 		{41.9097306, 12.2558141}, // Rome
// 		{-31.952222, 115.858889}, // Perth
// 		{42.441286, 19.262892},   // Podgorica
// 		{38.615925, -27.226598},  // Azores
// 		{45.4628329, 9.1076924},  // Milano
// 		{43.7800607, 11.170928},  // Florence
// 		{37.7775, -122.416389},   // San Francisco
// 	}

// 	m := onmap.Pins(coords, onmap.StandardCrop)
// 	f, err := os.Create("out_test.png")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// 	if err := png.Encode(f, m); err != nil {
// 		log.Fatal(err)
// 	}
// }

const (
	Width  = 1000
	Height = 800
)

// Convert lat/lon to pixel on an equirectangular map
func latLonToPixel(lat, lon float64, imgWidth, imgHeight int) (int, int) {
	x := int((lon + 180.0) * float64(imgWidth) / 360.0)
	y := int((90.0 - lat) * float64(imgHeight) / 180.0)
	return x, y
}

func main() {
	// === CONFIG ===
	startLat := 10.7769
	startLon := 106.7009
	scaleFactor := 2.0 // 2x zoom

	// === Load Image ===
	filePath := "out_test.png" // Full equirectangular map
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error loading image:", err)
	}
	defer imgFile.Close()

	imgSrc, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal("Error decoding image:", err)
	}
	imgBounds := imgSrc.Bounds()
	imgW, imgH := imgBounds.Dx(), imgBounds.Dy()

	// === GUI ===
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow := myApp.NewWindow("Interactive Map Zoom")

	// === Canvas Image ===
	fyneImg := canvas.NewImageFromFile(filePath)
	fyneImg.FillMode = canvas.ImageFillOriginal

	// Calculate scaled size
	scaledW := float32(float64(imgW) * scaleFactor)
	scaledH := float32(float64(imgH) * scaleFactor)
	fyneImg.Resize(fyne.NewSize(scaledW, scaledH))

	// Use container without layout to force actual image size
	imgContainer := container.NewWithoutLayout(fyneImg)
	imgContainer.Resize(fyne.NewSize(scaledW, scaledH))

	// === Make scrollable container ===
	scroll := container.NewScroll(imgContainer)

	// === Scroll to target location on start ===
	px, py := latLonToPixel(startLat, startLon, imgW, imgH)
	centerX := float32(float64(px)*scaleFactor - Width/2)
	centerY := float32(float64(py)*scaleFactor - Height/2)
	scroll.ScrollToOffset(fyne.NewPos(centerX, centerY))

	// === Show window ===
	myWindow.SetContent(scroll)
	myWindow.Resize(fyne.NewSize(Width, Height))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()

	fmt.Println("Exiting.")
}

type DraggableImage struct {
	widget.BaseWidget
	image   *canvas.Image
	offset  fyne.Position
	dragPos *fyne.Position
}

func NewDraggableImage(img *canvas.Image) *DraggableImage {
	d := &DraggableImage{
		image: img,
	}
	d.ExtendBaseWidget(d)
	return d
}

func (d *DraggableImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(d.image)
}

func (d *DraggableImage) Dragged(ev *fyne.DragEvent) {
	d.offset = d.offset.Add(ev.Dragged)
	d.image.Move(d.offset)
}

func (d *DraggableImage) DragEnd() {}

func (d *DraggableImage) MouseDown(ev *desktop.MouseEvent) {
	pos := ev.Position
	d.dragPos = &pos
}

// func main() {
// 	a := app.New()
// 	w := a.NewWindow("Simple Drag Example")

// 	img := canvas.NewImageFromFile("mercator.jpg") // Replace with your file
// 	img.FillMode = canvas.ImageFillOriginal

// 	draggable := NewDraggableImage(img)
// 	img.Move(fyne.NewPos(10, 10))

// 	c := container.NewWithoutLayout(draggable)
// 	c.Resize(fyne.NewSize(800, 600))

// 	// Optional: background so you can see it move
// 	bg := canvas.NewRectangle(color.NRGBA{R: 10, G: 10, B: 30, A: 255})
// 	bg.Resize(fyne.NewSize(1600, 1200)) // Bigger than the window

// 	layer := container.NewWithoutLayout(bg, draggable)

// 	w.SetContent(layer)
// 	w.Resize(fyne.NewSize(800, 600))
// 	w.ShowAndRun()
// }
