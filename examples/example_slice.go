package main

import (
	"bytes"
	"fmt"
	"github.com/goferHiro/imageslicer"
	"github.com/gorilla/websocket"
	"image"
	"image/jpeg"
	"log"
	"net/http"
)

var width, height int
var m image.Image

func getBytes(i image.Image) (b []byte) {
	var outWriter bytes.Buffer

	err := jpeg.Encode(&outWriter, i, nil)
	if err != nil {
		fmt.Println(err)
	}
	b = outWriter.Bytes()

	return
}

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

		if mt != websocket.TextMessage {
			c.WriteMessage(websocket.TextMessage, []byte("invalid msg"))
			continue
		}

		imageUrl := string(msg)

		inputImage := getImageFromUrl(imageUrl)

		c.WriteMessage(websocket.BinaryMessage, getBytes(inputImage))

		if err != nil {
			c.WriteMessage(websocket.CloseMessage, []byte("unable to decode the given image"))
			return
		}
		grid := [2]uint{4, 4}
		tiles := imageslicer.Slice(inputImage, grid)

		for _, tile := range tiles {
			c.WriteMessage(websocket.BinaryMessage, getBytes(tile))
		}
		outImage, err := imageslicer.Join(tiles, grid)

		if err != nil {
			c.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
			return
		}
		c.WriteMessage(websocket.TextMessage, []byte("joined image"))
		c.WriteMessage(websocket.BinaryMessage, getBytes(outImage))
	}
}

func getImageFromUrl(url string) (img image.Image) {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		// handle errors
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

func main() {
	fmt.Println("starting the server")
	http.HandleFunc("/ws/split", splitImage)
	http.ListenAndServe(":128", nil)
}
