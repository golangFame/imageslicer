package imageslicer

import (
	"image"
)

type Grid struct {
	Rows    int
	Columns int
}

func Slice(img image.Image, grid [2]int) (tiles []image.Image) {

	tiles = make([]image.Image, 0, grid[0]*grid[1])

	shape := img.Bounds()

	height := shape.Max.Y / grid[0]
	width := shape.Max.X / grid[1]

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
