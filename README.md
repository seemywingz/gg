# in3D
A relatively simple Go powered OpenGL Graphics Engine

Create a new Window, Get OpenGL Context, Setup Camera Projection, create 3D object, Draw!  
Go Ahead, you can do it yourself... `go get github.com/seemywingz/in3D`
```go
package main

import (
	"github.com/seemywingz/in3D"
)

func main() {

	in3D.Init(800, 600, "Simple Cube in3D")
	in3D.NewLight().Position = in3D.NewPosition(10, 1, 10)

	in3D.SetRelPath("../assets/textures")
	texture := in3D.NewTexture("seemywingz.jpg")
	color := []float32{1, 1, 1}

	obj := in3D.NewPointsObject(
		in3D.NewPosition(0, 0, -7),
		in3D.Cube,
		texture,
		color,
		in3D.Shader["phong"],
	)
	obj.SceneLogic = func(s *in3D.SceneData) {
		s.XRotation++
		s.YRotation++
	}

	for !in3D.ShouldClose() {
		in3D.Update()
		obj.Draw()
		in3D.SwapBuffers()
	}
}
```
![Simple Rotating Cude in3D](./examples/assets/textures/readme.png)
### ME-TODO:
  *  Optimize all the things!  
  *  Add Shadows, Ambient Occulsion and other light related things  
  * Have more fun making weird examples!  

### YOU-TODO:
  * Checkout the other examples to see some more basic functionality

### Installation on Ubuntu
```
apt-get install libgl1-mesa-dev

apt-get install libxrandr-dev

apt-get install libxcursor-dev

apt-get install libxi-dev

apt-get install libxinerama-dev

go get github.com/seemywingz/in3D

```  
  
##### Note:
###### Some Names and method may change until version 1.0 is tagged
###### Texture UVs are imported upside down due to the openGL/texture coordinate systems ( flip your texture vertiacally to render correctly  )

