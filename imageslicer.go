package imageslicer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strings"
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

func Join(tiles []image.Image, grid [2]uint) (img image.Image, err error) {

	expectedNoOfTiles := int(grid[0] * grid[1])

	if len(tiles) != expectedNoOfTiles || expectedNoOfTiles == 0 {
		err = fmt.Errorf("expected %d != %d", expectedNoOfTiles, len(tiles))
		return
	}

	i := 0

	shape := tiles[0].Bounds()

	height := shape.Max.Y * int(grid[0])
	width := shape.Max.X * int(grid[1])

	shapeOrig := image.Rect(shape.Min.X, shape.Min.Y, width, height)

	srcImage := image.NewRGBA(shapeOrig)

	for y := 0; y < int(grid[0]); y++ {
		for x := 0; x < int(grid[1]); x++ {

			tile := tiles[i]

			draw.Draw(srcImage, tile.Bounds(), tile, image.Point{
				x * shape.Max.X, y * shape.Max.Y,
			}, draw.Over)

			i += 1
			//shape.Min.X += shape.Max.X
		}
	}
	img = srcImage
	return
}

func GetBytes(i image.Image) (b []byte) {
	var outWriter bytes.Buffer

	err := jpeg.Encode(&outWriter, i, nil)
	if err != nil {
		fmt.Println(err)
	}
	b = outWriter.Bytes()

	return
}

func GetImageFromUrl(imgUrl string) (img image.Image) {
	res, err := http.Get(imgUrl)
	if err != nil {
		log.Println("err", err)
		return
	}
	defer res.Body.Close()
	img, _, err = image.Decode(res.Body)
	if err != nil {
		log.Println("err", err)
		return
	}
	return
}

func GetImageFromBase64(base64Img string) (img image.Image, err error) {

	mimeTypeIndex := strings.Index(base64Img, ";base64,")

	imageType := ""

	if mimeTypeIndex != -1 {
		mimeType := base64Img[:mimeTypeIndex]
		base64Img = strings.TrimPrefix(base64Img, mimeType+";base64,")
		imageType = strings.TrimPrefix(mimeType, "data:image/")

	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Img))

	switch imageType {

	case "jpeg", "jpg":
		img, err = jpeg.Decode(reader)
	case "png":
		img, err = png.Decode(reader)

	default:
		img, _, err = image.Decode(reader)
	}

	if err != nil {
		err = fmt.Errorf("unable to decode the img due to %s", err)
		return
	}
	return
}
