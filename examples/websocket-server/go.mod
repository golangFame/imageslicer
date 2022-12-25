module github.com/goferHiro/image-slicer/examples/websocket-server

go 1.19

replace github.com/goferHiro/image-slicer => ../..

require (
	github.com/goferHiro/image-slicer v1.0.1-alpha
	github.com/gorilla/websocket v1.5.0
)

retract v0.0.0-20221225000820-aabbc67a701f
