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
