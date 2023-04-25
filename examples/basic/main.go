package main

import (
	"fmt"
	"github.com/goferHiro/image-slicer"
	"log"
	"os"
)

func main() {
	// define the URL of the image to be sliced
	imgUrl := "https://static.wikia.nocookie.net/big-hero-6-fanon/images/0/0f/Hiro.jpg/revision/latest?cb=20180511180437"

	// download the image from the given URL and store it in a variable 'img'
	img := imageslicer.GetImageFromUrl(imgUrl)

	// if the image is not downloaded successfully, log the error and exit the program
	if img == nil {
		log.Fatalln("invalid image url or image format not supported!")
	}

	// define the grid structure to slice the image
	grid := imageslicer.Grid{2, 2} //rows,columns

	// slice the image using the defined grid
	tiles := imageslicer.Slice(img, grid)

	// create a directory to store the sliced tiles

	tilesDir := "sliced_tiles"

	err := os.Mkdir(tilesDir, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("failed to create directory: %v", err)

		}
	}

	// save each tile to the directory
	for i, tile := range tiles {
		tileName := fmt.Sprintf(fmt.Sprintf("%s/tile_%d", tilesDir, i))
		err := imageslicer.SaveTile(tile, tileName)
		if err != nil {
			log.Fatalf("failed to save tile %d: %v", i, err)
		}
	}

	fmt.Printf("Sliced tiles saved to directory %s successfully.\n", tilesDir)
}
