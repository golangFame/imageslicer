package main

import (
	"encoding/json"
	"fmt"
	imageslicer "github.com/goferHiro/image-slicer"
	"github.com/gorilla/websocket"
	"image"
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
		grid := [2]uint{4, 4}

		var imageUrl string

		var jsonMsg struct {
			Url  string  `json:"url"`
			Grid [2]uint `json:"grid"`
		}

		switch mt {
		case websocket.TextMessage:
			//imageUrl = string(msg)
			err := json.Unmarshal(msg, &jsonMsg)
			if err != nil {
				c.WriteMessage(websocket.TextMessage, []byte("prefer json msgs"))
				imageUrl = string(msg)
			} else {
				imageUrl = jsonMsg.Url
				grid = jsonMsg.Grid
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

		c.WriteMessage(websocket.BinaryMessage, imageslicer.GetBytes(inputImage))

		if err != nil {
			c.WriteMessage(websocket.CloseMessage, []byte("unable to decode the given image"))
			return
		}
		tiles := imageslicer.Slice(inputImage, grid)

		for _, tile := range tiles {
			c.WriteMessage(websocket.BinaryMessage, imageslicer.GetBytes(tile))
		}
		outImage, err := imageslicer.Join(tiles, grid)

		if err != nil {
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
