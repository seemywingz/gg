package main

import (
	"github.com/seemywingz/in3D"
)

func main() {

	in3D.Init(800, 600, "Simple Cube in3D")
	// in3D.NewLight().Position =
	// 	in3D.Position{X: 10, Y: 1, Z: 10}

	in3D.SetRelPath("../assets/textures")
	// texture := in3D.NewTexture("seemywingz.jpg")
	color := []float32{0.1, 0.1, 0.1}

	obj := in3D.NewPointsObject(
		in3D.NewPosition(0, 0, -7),
		in3D.Cube,
		in3D.NoTexture,
		color,
		in3D.Shader["color"],
	)
	obj.SceneLogic = func(s *in3D.SceneData) {
		s.XRotation += 0.1
		s.YRotation += 0.1
	}

	for !in3D.ShouldClose() {
		in3D.Update()
		obj.Draw()
		in3D.SwapBuffers()
	}
}
