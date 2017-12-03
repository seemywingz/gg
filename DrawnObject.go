package in3D

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ianremmler/ode"
)

// DrawnObject : a struct to hold openGL object data
type DrawnObject struct {
	Mesh           *Mesh
	Body           *ode.Body
	MVPID          int32
	ModelMatrixID  int32
	NormalMatrixID int32
	IambID         int32
	IdifID         int32
	IspecID        int32
	ShininessID    int32
	TextureID      int32
	NormalMapID    int32
	Scale          float32
	SceneData
}

// NewPointsObject :
func NewPointsObject(position Position, points []float32, texture uint32, color []float32, program uint32) *DrawnObject {
	vao := MakeVAO(points, program)
	mat := &Material{
		"default",
		color,
		color,
		color,
		1,
		texture,
		NoTexture,
	}
	mg := make(map[string]*MaterialGroup)
	mg["dfault"] = &MaterialGroup{
		mat,
		[]*Face{},
		vao,
		int32(len(points)),
	}
	mesh := &Mesh{mg}
	return NewMeshObject(position, mesh, program)
}

// NewMeshObject : Create new DrawnObject
func NewMeshObject(position Position, mesh *Mesh, program uint32) *DrawnObject {

	ModelMatrixID := gl.GetUniformLocation(program, gl.Str("MODEL\x00"))
	NormalMatrixID := gl.GetUniformLocation(program, gl.Str("NormalMatrix\x00"))
	MVPID := gl.GetUniformLocation(program, gl.Str("MVP\x00"))

	uniform := "Material"
	IambID := gl.GetUniformLocation(program, gl.Str(uniform+".Iamb\x00"))
	IdifID := gl.GetUniformLocation(program, gl.Str(uniform+".Idif\x00"))
	IspecID := gl.GetUniformLocation(program, gl.Str(uniform+".Ispec\x00"))
	ShininessID := gl.GetUniformLocation(program, gl.Str(uniform+".shininess\x00"))
	TextureID := gl.GetUniformLocation(program, gl.Str("TEXTURE\x00"))
	NoramalMapID := gl.GetUniformLocation(program, gl.Str("NORMAL_MAP\x00"))
	// println(TextureID, NoramalMapID)

	d := &DrawnObject{
		mesh,
		initBodyPhisics(position),
		MVPID,
		ModelMatrixID,
		NormalMatrixID,
		IambID,
		IdifID,
		IspecID,
		ShininessID,
		TextureID,
		NoramalMapID,
		1,
		SceneData{},
	}
	d.Position = position
	d.Program = program
	return d
}

func initBodyPhisics(pos Position) *ode.Body {
	if !Feature[Physics] {
		return nil
	}

	body := world.NewBody()
	body.SetPosition(ode.V3(float64(pos.X), float64(pos.Y), float64(pos.Z)))
	body.SetAngularVelocity(ode.V3(0, 0, 0))
	body.SetAngularDamping(0.02)
	body.SetLinearDamping(0.02)

	mass := ode.NewMass()
	mass.SetBox(1, ode.V3(1, 1, 1))
	// mass.Adjust(sphereMass)
	body.SetMass(mass)

	box := space.NewBox(ode.V3(1, 1, 1))
	box.SetBody(body)

	return &body
}

func (d *DrawnObject) translateRotate() *mgl32.Mat4 {
	model := mgl32.Translate3D(d.X, d.Y, d.Z).
		Mul4(mgl32.Scale3D(d.Scale, d.Scale, d.Scale))
	xrotMatrix := mgl32.HomogRotate3DX(mgl32.DegToRad(d.XRotation))
	yrotMatrix := mgl32.HomogRotate3DY(mgl32.DegToRad(d.YRotation))
	zrotMatrix := mgl32.HomogRotate3DZ(mgl32.DegToRad(d.ZRotation))
	if Feature[Physics] {
		p := d.Body.Position()
		// fmt.Println(p)
		d.Position.X = float32(p[0])
		d.Position.Y = float32(p[1])
		d.Position.Z = float32(p[2])
		// zrotMatrix = mgl32.HomogRotate3DZ(float32(d.Body.GetAngle()))
	}
	final := model.Mul4(xrotMatrix.Mul4(yrotMatrix.Mul4(zrotMatrix)))
	return &final
}

// Draw : draw the object
func (d *DrawnObject) Draw() {

	if d.SceneLogic != nil {
		d.SceneLogic(&d.SceneData)
	}

	modelMatrix := d.translateRotate()
	normalMatrix := modelMatrix.Inv().Transpose()

	gl.UseProgram(d.Program)
	gl.UniformMatrix4fv(d.MVPID, 1, false, &camera.MVP[0])
	gl.UniformMatrix4fv(d.ModelMatrixID, 1, false, &modelMatrix[0])
	gl.UniformMatrix4fv(d.NormalMatrixID, 1, false, &normalMatrix[0])

	for _, m := range d.Mesh.MaterialGroups {
		gl.UseProgram(d.Program)
		gl.BindVertexArray(m.VAO)

		// Material
		gl.Uniform3fv(d.IambID, 1, &m.Material.Ambient[0])
		gl.Uniform3fv(d.IspecID, 1, &m.Material.Specular[0])
		gl.Uniform3fv(d.IdifID, 1, &m.Material.Diffuse[0])
		gl.Uniform1f(d.ShininessID, m.Material.Shininess)

		gl.Uniform1i(d.TextureID, 0)
		gl.Uniform1i(d.NormalMapID, 1)

		// Bind our diffuse texture in Texture Unit 0
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, m.Material.DiffuseTex)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, m.Material.NormalTex)

		gl.DrawArrays(gl.TRIANGLES, 0, m.VertCount)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}

}
