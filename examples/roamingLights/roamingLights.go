package main

import (
	"math/rand"
	"time"

	"github.com/seemywingz/in3D"
)

var (
	texture      map[string]uint32
	lights       []in3D.Light
	drawnObjects []*in3D.DrawnObject
)

func randObjects(numberOfObjects, min, max int, points []float32, textr, shadr uint32) {

	for i := 0; i < numberOfObjects; i++ {

		rand.Seed(time.Now().UnixNano())
		x, y, z := in3D.Random(min, max), in3D.Random(min, max), in3D.Random(min, max)
		rx, ry, rz := in3D.Randomf(), in3D.Randomf(), in3D.Randomf()
		if i%101 == 0 {
			color := []float32{rx, ry, rz}
			roamingLight := in3D.NewColorLight([]float32{0.1, 0.1, 0.1}, color, color)
			roamingLight.Position = in3D.NewPosition(float32(x), float32(y), float32(z))
			roamingLight.Draw = true
			roamingLight.SceneLogic = func(s *in3D.SceneData) {
				s.X += rx * 3
				s.Y += ry * 3
				s.Z += rz * 3
				if s.X > float32(max) || s.X < float32(min) {
					rx = -rx
				}
				if s.Y > float32(max) || s.Y < float32(min) {
					ry = -ry
				}
				if s.Z > float32(max) || s.Z < float32(min) {
					rz = -rz
				}
			}
		} else {
			var color []float32
			// if textr == in3D.NoTexture {
			// 	color = []float32{rx, ry, rz}
			// } else {
			color = []float32{1, 1, 1}
			// }
			d := in3D.NewPointsObject(
				in3D.NewPosition(float32(x), float32(y), float32(z)),
				points,
				textr,
				color,
				shadr)
			// d.SceneLogic = func(s *in3D.SceneData) {
			// s.XRotation += rx
			// s.YRotation += ry
			// s.ZRotation += rz
			// }
			drawnObjects = append(drawnObjects, d)
		}
	}
}

func loadTextures() {
	texture = make(map[string]uint32)
	texture["none"] = in3D.NoTexture
	in3D.SetRelPath("../assets/textures")
	texture["box"] = in3D.NewTexture("box.jpg")
	texture["box1"] = in3D.NewTexture("box1.jpg")
}

func main() {

	in3D.Init(000, 000, "Roaming Light")
	in3D.SetCameraPosition(in3D.NewPosition(0, 15, 130))
	in3D.Enable(in3D.PointerLock, true)
	in3D.Enable(in3D.FlyMode, true)

	loadTextures()

	min, max := -100, 100
	randObjects(800, min, max, in3D.Cube, texture["none"], in3D.Shader["phong"])
	randObjects(100, min, max, in3D.Cube, texture["box"], in3D.Shader["phong"])
	randObjects(100, min, max, in3D.Cube, texture["box1"], in3D.Shader["phong"])

	centerLight := in3D.NewLight()
	centerLight.Draw = true
	println(centerLight.DrawnObject.IdifID)

	for !in3D.ShouldClose() {
		in3D.Update()

		for _, obj := range drawnObjects {
			obj.Draw()
		}

		in3D.SwapBuffers()
	}
}
