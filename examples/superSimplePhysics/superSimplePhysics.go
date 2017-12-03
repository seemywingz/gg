package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/seemywingz/in3D"
)

var objects []*in3D.DrawnObject

func main() {

	in3D.Init(800, 600, "Simple Cube in3D")
	window := in3D.GetWindow()
	light := in3D.NewLight()
	light.Position = in3D.NewPosition(0, 1, 10)
	light.Radius = 1000

	in3D.GetCamera().Position = in3D.NewPosition(0, 80, 400)

	in3D.Enable(in3D.Physics, true)
	in3D.Enable(in3D.PointerLock, true)
	in3D.Enable(in3D.FlyMode, true)

	explode()

	for !in3D.ShouldClose() {
		in3D.Update()

		// Press w
		if window.GetKey(glfw.KeyF) == glfw.Press {
			explode()
		}

		for _, obj := range objects {
			obj.Draw()
		}

		in3D.SwapBuffers()
	}
}

func explode() {
	for i := 1; i < 1000; i++ {
		obj := in3D.NewPointsObject(
			in3D.NewPosition(0, float32(i), -20),
			in3D.Cube,
			in3D.NoTexture,
			[]float32{0, 1, 1},
			in3D.Shader["phong"],
		)
		objects = append(objects, obj)
	}
}
