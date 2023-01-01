module github.com/goferHiro/image-slicer

go 1.19

retract (
	v1.0.1-beta //cannot find package issue
	v1.0.1-alpha

	v1.0.0 //doesn't support examples due to replace
	//unstable versions
	v0.8.1-alpha.1
	v0.8.1-alpha.1

)

require (
	github.com/dvyukov/go-fuzz v0.0.0-20220726122315-1d375ef9f9f6 // indirect
	github.com/dvyukov/go-fuzz-corpus v0.0.0-20190920191254-c42c1b2914c7 // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.1 // indirect
	github.com/stephens2424/writerset v1.0.2 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/tools v0.4.0 // indirect
)
