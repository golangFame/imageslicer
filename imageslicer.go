package imageslicer

import (
	"image"
	"math"
)

type Grid struct {
	Rows    int
	Columns int
}

func Slice(img image.Image, grid [2]int) (tiles []image.Image) {

	tiles = make([]image.Image, 0, grid[0]*grid[1])

	shape := img.Bounds()

	heightF := float64(shape.Max.Y / grid[0])
	widthF := float64(shape.Max.X / grid[1])

	height := int(math.Ceil(heightF))
	width := int(math.Ceil(widthF))

	for y := shape.Min.Y; y+height <= shape.Max.Y; y += height {

		for x := shape.Min.X; x+width <= shape.Max.X; x += width {

			tile := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(x, y, x+width, y+height))

			tiles = append(tiles, tile)
		}
	}

	return
}
