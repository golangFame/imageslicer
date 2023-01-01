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
	"math"
	"net/http"
	"os"
	"path"
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

	fheight := float64(shape.Max.Y / int(grid[0]))
	fwidth := float64(shape.Max.X / int(grid[1]))

	height := int(math.Ceil(fheight))
	width := int(math.Ceil(fwidth))

	//fmt.Printf("h%dw%d\n", height, width)

	for y := shape.Min.Y; y+height <= shape.Max.Y; y += height {

		for x := shape.Min.X; x+width <= shape.Max.X; x += width {

			tile := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(x, y, x+width, y+height))

			tiles = append(tiles, tile)
		}
	}

	/*	if int(grid[0]*grid[1]) != len(tiles) {
		log.Fatalf("expected %v got %d", grid[0]*grid[1], len(tiles))
	}*/

	return
}

func Join(tiles []image.Image, grid [2]uint) (img image.Image, err error) {

	if err = CheckSlice(tiles, grid); err != nil {
		return
	}

	i := 0

	shape := tiles[0].Bounds()

	height := shape.Max.Y * int(grid[0])
	width := shape.Max.X * int(grid[1])

	shapeOrig := image.Rect(shape.Min.X, shape.Min.Y, width, height)

	srcImage := image.NewRGBA(shapeOrig)

	for y, minY := 0, shapeOrig.Min.Y; y < int(grid[0]); y++ {

		for x, minX := 0, shapeOrig.Min.X; x < int(grid[1]); x++ {

			tile := tiles[i]
			shape := tile.Bounds()

			draw.Draw(srcImage, shape, tile, shape.Min, draw.Src)

			i += 1

			minX += shape.Min.X //not required really

		}
		minY += shape.Min.Y

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

func GetImageFromUrl(imgUrl string) (img image.Image) { //FIXME add error return and remove log
	res, err := http.Get(imgUrl)
	if err != nil {
		//log.Printf("http-res %d %s\n", res.StatusCode, res.Status)
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

func GetImageFromPath(imgPath string) (img image.Image, err error) {

	f1, err := os.Open(imgPath)
	if err != nil {
		err = fmt.Errorf("failed to open image due to %s", err)
		return
	}
	defer f1.Close()

	imgType := path.Ext(imgPath)

	imgType = strings.TrimPrefix(imgType, ".")

	switch imgType {

	case "jpeg", "jpg":
		img, err = jpeg.Decode(f1)
	case "png":
		img, err = png.Decode(f1)

	default:
		log.Printf("img type-%s may not be supported", imgType) //TODO bring in the logger
		img, _, err = image.Decode(f1)
	}

	if err != nil {
		err = fmt.Errorf("failed to decode image due to %s", err)
		return
	}

	return
}

func CheckSlice(tiles []image.Image, grid [2]uint) (err error) {
	expectedNoOfTiles := int(grid[0] * grid[1])

	if len(tiles) < expectedNoOfTiles {
		err = fmt.Errorf("expected-%d got-%d", expectedNoOfTiles, len(tiles))
	}
	return
}
