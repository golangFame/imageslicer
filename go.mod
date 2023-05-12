module github.com/golangFame/imageslicer

go 1.19

retract v0.0.0-20221225000820-aabbc67a701f

retract (
	v1.4.0-beta-1
	v1.4.0-alpha-2
	v1.4.0-alpha
	v1.3.1-beta

	v1.3.0-alpha
	v1.2.2 //dep issues for eg
	v1.2.1 //dep issues for examples

	v1.2.0 //dep issues for examples
	v1.0.1-beta //cannot find package issue
	v1.0.1-alpha

	v1.0.0 //doesn't support examples due to replace
	//unstable versions
	v0.8.1-alpha.1
	v0.8.1-alpha.1
)

require github.com/google/gofuzz v1.2.0
