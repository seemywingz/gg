package main

import (
	"github.com/bytearena/box2d"
	"github.com/seemywingz/in3D"
)

var (
	gravity box2d.B2Vec2
	world   box2d.B2World
	body    *box2d.B2Body
)

func main() {

	in3D.Init(800, 600, "Simple Cube in3D")
	in3D.NewLight().Position = in3D.NewPosition(10, 1, 10)
	in3D.Enable(in3D.Physics, true)

	in3D.SetRelPath("../assets/textures")
	texture := in3D.NewTexture("seemywingz.jpg")

	obj := in3D.NewPointsObject(
		in3D.NewPosition(0, 30, -20),
		in3D.Cube,
		texture,
		[]float32{0, 1, 1},
		in3D.Shader["phong"],
	)

	obj2 := in3D.NewPointsObject(
		in3D.NewPosition(1.5, 2, -20),
		in3D.Cube,
		texture,
		[]float32{1, 0, 1},
		in3D.Shader["phong"],
	)

	for !in3D.ShouldClose() {
		in3D.Update()
		obj.Draw()
		obj2.Draw()
		in3D.SwapBuffers()
	}
}
