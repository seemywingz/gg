package main

import (
	"github.com/seemywingz/in3D"
)

func main() {

	var objects []*in3D.DrawnObject

	in3D.Init(800, 600, "Wavefront Loader")
	in3D.SetClearColor(0.1, 0.1, 0.1, 1)
	in3D.SetCameraPosition(in3D.NewPosition(0, 0.55, 2))
	in3D.SetCameraSpeed(0.1)
	in3D.Enable(in3D.PointerLock, true)
	in3D.Enable(in3D.FlyMode, true)

	light := in3D.NewLight()
	light.Ambient = []float32{0.5, 0.5, 0.5}
	light.Specular = []float32{10, 10, 10}
	light.Position = in3D.NewPosition(1, 1, 0)
	light.Draw = true
	light.DrawnObject.Scale = 0.05
	light.Radius = 10

	model := "sky"
	skyShader := in3D.Shader["texture"]
	in3D.SetRelPath("../assets/models/" + model)
	skymesh := in3D.LoadObject(model+".obj", skyShader)
	sky := in3D.NewMeshObject(in3D.Position{}, skymesh, skyShader)
	sky.Scale = 10000
	objects = append(objects, sky)

	// all models are from: https://www.blendswap.com/
	model = "buddha"
	meshShader := in3D.Shader["normalMap"]
	in3D.SetRelPath("../assets/models/" + model)
	mesh := in3D.LoadObject(model+".obj", meshShader)
	buddha := in3D.NewMeshObject(in3D.Position{}, mesh, meshShader)
	buddha.SceneLogic = func(s *in3D.SceneData) {
		s.YRotation += 0.1
	}
	objects = append(objects, buddha)
	// objects = append(objects, in3D.NewPointsObject(in3D.NewPosition(1, 1, -10), in3D.Cube, in3D.NoTexture, []float32{1, 1, 1}, in3D.Shader["normalMap"]))

	for !in3D.ShouldClose() {
		in3D.Update()
		for _, o := range objects {
			o.Draw()
		}
		in3D.SwapBuffers()
	}
}
