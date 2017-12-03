package main

import (
	"github.com/seemywingz/in3D"
)

var objects []*in3D.DrawnObject

func main() {

	in3D.Init(800, 600, "Simple Cube in3D")
	in3D.NewLight().Position = in3D.NewPosition(10, 1, 10)
	in3D.GetCamera().Position = in3D.NewPosition(0, 8, 0)
	in3D.Enable(in3D.Physics, true)

	in3D.SetRelPath("../assets/textures")
	texture := in3D.NewTexture("seemywingz.jpg")

	for i := 0; i < 10; i++ {
		obj := in3D.NewPointsObject(
			in3D.NewPosition(0, 30, -20),
			in3D.Cube,
			texture,
			[]float32{0, 1, 1},
			in3D.Shader["phong"],
		)
		objects = append(objects, obj)
	}

	for !in3D.ShouldClose() {
		in3D.Update()

		for _, obj := range objects {
			obj.Draw()
		}

		in3D.SwapBuffers()
	}
}
