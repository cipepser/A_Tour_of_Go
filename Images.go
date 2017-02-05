package main

import (
	"image"
	"image/color"
	"golang.org/x/tour/pic"
)

type Image struct{}

func (img *Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (img *Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, 200, 200)
}

func (img *Image) At(x, y int) color.Color {
	return color.RGBA{uint8(x % 256), uint8(y %  256), uint8((x * y) % 256 ), 255}
}


func main() {
	m := Image{}
	pic.ShowImage(m)
}