package main

import (
	"encoding/json"
	"fmt"
	"github.com/goferHiro/image-slicer"
	"github.com/gorilla/websocket"
	"image"
	"image/color"
	"image/draw"
	"log"
	"net/http"
)

var width, height int
var m image.Image

func splitImage(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	c.WriteMessage(websocket.TextMessage, []byte("lets start"))

	for {
		c.WriteMessage(websocket.TextMessage, []byte("send the image url"))

		mt, msg, err := c.ReadMessage()
		grid := imageslicer.Grid{4, 4}

		var imageUrl string

		var req struct {
			Url    string           `json:"url"`
			Grid   imageslicer.Grid `json:"grid"`
			Border int              `json:"border"`
		}

		switch mt {
		case websocket.TextMessage:
			//imageUrl = string(msg)
			err := json.Unmarshal(msg, &req)
			if err != nil {
				c.WriteMessage(websocket.TextMessage, []byte("prefer json msgs"))
				imageUrl = string(msg)
			} else {
				imageUrl = req.Url
				grid = req.Grid
			}

		default:
			c.WriteMessage(websocket.TextMessage, []byte("invalid msg"))
			continue

		}

		inputImage := imageslicer.GetImageFromUrl(imageUrl)

		if inputImage == nil {
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("get image from url failed for %s", imageUrl)))
			continue
		}

		log.Println("recvd image")

		c.WriteMessage(websocket.BinaryMessage, imageslicer.GetBytes(inputImage))

		tiles := imageslicer.Slice(inputImage, grid)

		for i, tile := range tiles {
			border := req.Border
			if border != 0 {
				tiles[i] = drawBorder(tile, border)
				tile = tiles[i]
			}
			//c.WriteMessage(websocket.BinaryMessage, imageslicer.GetBytes(tile))
		}
		outImage, err := imageslicer.Join(tiles, grid)

		if err != nil {
			log.Fatalf("join failed -%v", err)
			c.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
			return
		}
		c.WriteMessage(websocket.TextMessage, []byte("joined image"))
		c.WriteMessage(websocket.BinaryMessage, imageslicer.GetBytes(outImage))
	}
}

func main() {
	fmt.Println("starting the server")
	http.HandleFunc("/ws/split", splitImage)
	http.ListenAndServe(":128", nil)
}

func drawBorder(img image.Image, borderSize int) (borderedImg image.Image) {

	bounds := img.Bounds()

	newBounds := image.Rect(bounds.Min.X-borderSize, bounds.Min.Y-borderSize, bounds.Max.X+borderSize, bounds.Max.Y+borderSize)
	newImg := image.NewRGBA(newBounds)

	// Draw the original image onto the new image
	draw.Draw(newImg, bounds, img, bounds.Min, draw.Src)

	// Draw the border onto the new image
	borderColor := color.RGBA{0, 0, 0, 0}
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{borderColor}, image.ZP, draw.Over)

	borderedImg = newImg

	return
	/*	// Draw the border onto the new image
		borderColor := color.RGBA{0, 0, 255, 1}
		draw.Draw(img.(interface {
			//RGBA64At(x int, y int) color.RGBA64
			//PixOffset(x int, y int) int
			//RGBAAt(x int, y int) color.RGBA
			//SetRGBA64(x int, y int, c color.RGBA64)
			//SetRGBA(x int, y int, c color.RGBA)

			Bounds() image.Rectangle
			At(x int, y int) color.Color

			ColorModel() color.Model

			Set(x int, y int, c color.Color)

			SubImage(r image.Rectangle) image.Image
			Opaque() bool
		}), img.Bounds(), &image.Uniform{borderColor}, image.ZP, draw.Src)*/

}
