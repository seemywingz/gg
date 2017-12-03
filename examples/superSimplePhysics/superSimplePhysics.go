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
	color := []float32{1, 1, 1}

	obj := in3D.NewPointsObject(
		in3D.NewPosition(0, 10, -20),
		in3D.Cube,
		texture,
		color,
		in3D.Shader["phong"],
	)

	for !in3D.ShouldClose() {
		in3D.Update()
		obj.Draw()
		in3D.SwapBuffers()
	}
}
