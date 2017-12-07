package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/seemywingz/in3D"
)

var objects []*in3D.DrawnObject
var exploding bool
var texture uint32

func main() {

	in3D.Init(800, 600, "Simple Cube in3D")
	window := in3D.GetWindow()

	light := in3D.NewLight()
	light.Position = in3D.NewPosition(0, 100, 0)
	light.Radius = 10000

	in3D.GetCamera().Position = in3D.NewPosition(0, 2, 10)

	in3D.Enable(in3D.Physics, true)
	in3D.Enable(in3D.PointerLock, true)
	in3D.Enable(in3D.FlyMode, true)

	in3D.SetRelPath("../assets/textures")
	texture = in3D.NoTexture
	// texture = in3D.NewTexture("box.jpg")
	explode()

	for !in3D.ShouldClose() {
		in3D.Update()

		// Press e
		if window.GetKey(glfw.KeyE) == in3D.Press && !exploding {
			explode()
		}

		for _, obj := range objects {
			obj.Draw()
		}

		in3D.SwapBuffers()
	}
}

func explode() {
	for i := 1; i < 10; i++ {
		obj := in3D.NewPointsObject(
			in3D.NewPosition(0, float32(i*10), -20),
			in3D.Cube,
			texture,
			[]float32{1, 1, 1},
			// []float32{in3D.Randomf(), in3D.Randomf(), in3D.Randomf()},
			in3D.Shader["phong"],
		)
		objects = append(objects, obj)
	}
}
