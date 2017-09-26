package gg

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// LightManager :
type LightManager struct {
	Lights  []*Light
	Program uint32
}

// Light : struct to hold light data
type Light struct {
	Radius  float32
	Iamb    []float32
	Idif    []float32
	Ispec   []float32
	LRadID  int32
	LPosID  int32
	IambID  int32
	IdifID  int32
	IspecID int32
	SceneData
	Draw        bool
	DrawnObject *DrawnObject
}

// NewLightManager :
func NewLightManager() *LightManager {
	lightManager = &LightManager{
		[]*Light{},
		Shader["phong"],
	}
	return lightManager
}

// NewLight :
func NewLight(color []float32) *Light {
	return NewCustomLight(
		NewPosition(0, 0, 0), // position
		50,                   // radius
		[]float32{0.2, 0.2, 0.2}, // ambiant intensity
		color, // diffuse intensity
		[]float32{0.8, 0.8, 0.8}, // specular intensity
		false, // draw
	)
}

// NewCustomLight :
func NewCustomLight(position Position, radius float32, iamb, idif, ispec []float32, draw bool) *Light {
	n := len(lightManager.Lights)
	if n > MaxLights {
		EoE("Error adding New Light:", errors.New("Max lights reached "+string(MaxLights)))
	}
	fmt.Println("Adding Light:", n)
	uniform := fmt.Sprintf("Light[%v]", n)
	LRadID := gl.GetUniformLocation(lightManager.Program, gl.Str(uniform+".lightRad\x00"))
	LPosID := gl.GetUniformLocation(lightManager.Program, gl.Str(uniform+".lightPos\x00"))
	IambID := gl.GetUniformLocation(lightManager.Program, gl.Str(uniform+".Iamb\x00"))
	IdifID := gl.GetUniformLocation(lightManager.Program, gl.Str(uniform+".Idif\x00"))
	IspecID := gl.GetUniformLocation(lightManager.Program, gl.Str(uniform+".Ispec\x00"))

	drawnObject := NewPointsObject(position, Cube, NoTexture, idif, Shader["color"])

	light := &Light{
		radius,
		iamb,
		idif,
		ispec,
		LRadID,
		LPosID,
		IambID,
		IdifID,
		IspecID,
		SceneData{},
		draw,
		drawnObject,
	}
	light.Position = position
	lightManager.Lights = append(lightManager.Lights, light)
	return light
}

// Update :
func (l *LightManager) Update() {

	for _, light := range l.Lights {
		if light.SceneLogic != nil {
			light.SceneLogic(&light.SceneData)
		}

		if light.Draw {
			light.DrawnObject.Position = light.Position
			light.DrawnObject.Draw()
		}
		gl.UseProgram(l.Program)
		gl.Uniform1f(light.LRadID, light.Radius)
		gl.Uniform3fv(light.LPosID, 1, &[]float32{light.X, light.Y, light.Z}[0])
		gl.Uniform3fv(light.IambID, 1, &light.Iamb[0])
		gl.Uniform3fv(light.IdifID, 1, &light.Idif[0])
		gl.Uniform3fv(light.IspecID, 1, &light.Ispec[0])

	}
}
