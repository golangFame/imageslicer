module github.com/goferHiro/image-slicer/examples/websocket-server

go 1.19

replace github.com/goferHiro/image-slicer => ../..


require github.com/goferHiro/image-slicer v0.9-alpha.2

retract (
    v0.1-alpha.1
)