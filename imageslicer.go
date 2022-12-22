package imageslicer

import (
	"image"
	"math"
)

func Slice(img image.Image, grid [2]int) (tiles []image.Image) {

	tiles = make([]image.Image, 0, grid[0]*grid[1])

	shape := img.Bounds()

	widthF := float64(shape.Max.X / grid[1])
	heightF := float64(shape.Max.Y / grid[0])

	width := int(math.Ceil(widthF))
	height := int(math.Ceil(heightF))

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
