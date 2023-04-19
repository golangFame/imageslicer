package main

import (
	"fmt"
	"github.com/goferHiro/image-slicer"
	"log"
)

func main() {
	imgUrl := "https://static.wikia.nocookie.net/big-hero-6-fanon/images/0/0f/Hiro.jpg/revision/latest?cb=20180511180437"

	img := imageslicer.GetImageFromUrl(imgUrl)

	if img == nil {
		log.Fatalln("invalid image url or image format not supported!")
	}

	grid := imageslicer.Grid{2, 2} //rows,columns

	tiles := imageslicer.Slice(img, grid)

	expectedTiles := int(grid[0] * grid[1])

	if len(tiles) != expectedTiles {
		log.Fatalf("expected %d rcvd %d\n", expectedTiles, len(tiles))
	}

	//lets join the tiles back

	joinedImg, err := imageslicer.Join(tiles, grid)

	if err != nil {
		log.Fatalf("joining tiles failed due to %s", err)
	}

	shapeJ := joinedImg.Bounds()
	shapeI := img.Bounds()

	fmt.Println(shapeJ, shapeI) //shape might change due to pixel loss

}
