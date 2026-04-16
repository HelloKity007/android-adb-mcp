package snapshot

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Annotator handles drawing annotations on screenshots
type Annotator struct {
	Scale float64
}

// NewAnnotator creates a new screenshot annotator
func NewAnnotator(scale float64) *Annotator {
	if scale <= 0 {
		scale = 0.7
	}
	return &Annotator{Scale: scale}
}

// getRandomColor generates a random color
func getRandomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

// drawRectangle draws a rectangle outline on the image
func drawRectangle(img *image.RGBA, rect image.Rectangle, col color.RGBA, width int) {
	for i := 0; i < width; i++ {
		// Top and bottom edges
		for x := rect.Min.X; x < rect.Max.X; x++ {
			if rect.Min.Y+i >= 0 && rect.Min.Y+i < img.Bounds().Max.Y {
				img.Set(x, rect.Min.Y+i, col)
			}
			if rect.Max.Y-1-i >= 0 && rect.Max.Y-1-i < img.Bounds().Max.Y {
				img.Set(x, rect.Max.Y-1-i, col)
			}
		}
		// Left and right edges
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			if rect.Min.X+i >= 0 && rect.Min.X+i < img.Bounds().Max.X {
				img.Set(rect.Min.X+i, y, col)
			}
			if rect.Max.X-1-i >= 0 && rect.Max.X-1-i < img.Bounds().Max.X {
				img.Set(rect.Max.X-1-i, y, col)
			}
		}
	}
}

// drawFilledRectangle draws a filled rectangle on the image
func drawFilledRectangle(img *image.RGBA, rect image.Rectangle, col color.RGBA) {
	draw.Draw(img, rect, &image.Uniform{col}, image.Point{}, draw.Src)
}

// drawText draws text on the image at the specified position
func drawText(img *image.RGBA, x, y int, text string, col color.RGBA) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(text)
}

// AnnotateScreenshot draws bounding boxes and labels on a screenshot
func (a *Annotator) AnnotateScreenshot(screenshot image.Image, elements []ElementNode) (*image.RGBA, error) {
	bounds := screenshot.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, screenshot, image.Point{}, draw.Src)

	for i, elem := range elements {
		elemColor := getRandomColor()

		// Scale bounding box
		x1 := int(float64(elem.BoundingBox.X1) * a.Scale)
		y1 := int(float64(elem.BoundingBox.Y1) * a.Scale)
		x2 := int(float64(elem.BoundingBox.X2) * a.Scale)
		y2 := int(float64(elem.BoundingBox.Y2) * a.Scale)

		// Draw bounding box
		rect := image.Rect(x1, y1, x2, y2)
		drawRectangle(rgba, rect, elemColor, 2)

		// Draw label
		label := fmt.Sprintf("%d", i)

		labelWidth := len(label) * 7 // basicfont width
		labelHeight := 13            // basicfont height

		// Label position above bounding box, clamped to image bounds
		labelX1 := x2 - labelWidth
		if labelX1 < 0 {
			labelX1 = 0
		}
		labelY1 := y1 - labelHeight - 4
		if labelY1 < 0 {
			labelY1 = 0
		}
		labelX2 := labelX1 + labelWidth
		labelY2 := labelY1 + labelHeight + 4

		// Draw label background
		labelRect := image.Rect(labelX1, labelY1, labelX2, labelY2)
		drawFilledRectangle(rgba, labelRect, elemColor)

		// Draw label text (white)
		whiteColor := color.RGBA{255, 255, 255, 255}
		drawText(rgba, labelX1+2, labelY1+labelHeight, label, whiteColor)
	}

	return rgba, nil
}

// LoadImage loads an image from file
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

// SavePNG saves an image as PNG
func SavePNG(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// ImageToPNG encodes an image to PNG bytes
func ImageToPNG(img image.Image) ([]byte, error) {
	var buf []byte
	writer := &bytesWriter{data: &buf}
	err := png.Encode(writer, img)
	return buf, err
}

type bytesWriter struct {
	data *[]byte
}

func (w *bytesWriter) Write(p []byte) (n int, err error) {
	*w.data = append(*w.data, p...)
	return len(p), nil
}
