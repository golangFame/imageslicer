package imageslicer

import (
	"image"
)

type Grid struct {
	Rows    int
	Columns int
}

func Slice(img image.Image, grid [2]uint) (tiles []image.Image) {

	tiles = make([]image.Image, 0, grid[0]*grid[1])

	if cap(tiles) == 0 {
		return
	}

	shape := img.Bounds()

	height := shape.Max.Y / int(grid[0])
	width := shape.Max.X / int(grid[1])

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
