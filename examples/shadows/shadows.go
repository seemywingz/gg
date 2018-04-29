package main

import (
	"github.com/seemywingz/in3D"
)

func main() {

	var objects []*in3D.DrawnObject

	in3D.Init(000, 600, "Example")
	in3D.SetClearColor(0.1, 0.1, 0.1, 1)
	in3D.SetCameraPosition(in3D.NewPosition(0, 0.55, 2))
	in3D.SetCameraSpeed(0.1)
	in3D.Enable(in3D.PointerLock, true)
	in3D.Enable(in3D.FlyMode, true)

	light := in3D.NewLight()
	// light.Specular = []float32{50, 50, 50}
	light.Position = in3D.NewPosition(0, 1, 1)
	light.Draw = true
	light.DrawnObject.Scale = 0.05
	light.Radius = 10

	dx := float32(0.01)
	n := float32(0)
	light.SceneLogic = func(s *in3D.SceneData) {
		n += dx
		max := float32(5)
		if n > max || n < -max {
			dx = -dx
		}
		s.Position.Z = n
	}

	model := "sky"
	skyShader := in3D.Shader["texture"]
	in3D.SetRelPath("../assets/models/" + model)
	skymesh := in3D.LoadObject(model+".obj", skyShader)
	sky := in3D.NewMeshObject(in3D.Position{}, skymesh, skyShader)
	sky.Scale = 10000
	objects = append(objects, sky)

	sky.SceneLogic = func(s *in3D.SceneData) {
		s.YRotation += 0.01
	}

	// all models are from: https://www.blendswap.com/
	model = "buddha"
	meshShader := in3D.Shader["in3D"]
	in3D.SetRelPath("../assets/models/" + model)
	mesh := in3D.LoadObject("buddha.obj", meshShader)
	meshObject := in3D.NewMeshObject(in3D.NewPosition(-0.5, 0, 0), mesh, meshShader)
	meshObject.YRotation = 90
	objects = append(objects, meshObject)

	meshShader = in3D.Shader["phong"]
	mesh = in3D.LoadObject("buddha.obj", meshShader)
	buddha := in3D.NewMeshObject(in3D.NewPosition(0.5, 0, 0), mesh, meshShader)
	buddha.YRotation = -90
	objects = append(objects, buddha)

	plane := in3D.NewPointsObject(
		in3D.NewPosition(0, 0, 0),
		in3D.Plane, in3D.NoTexture,
		[]float32{1, 1, 1},
		in3D.Shader["phong"])
	plane.XRotation = -90
	plane.Scale = 500
	objects = append(objects, plane)

	for !in3D.ShouldClose() {
		in3D.Update()
		for _, object := range objects {
			object.Draw()
		}
		in3D.SwapBuffers()
	}
}
