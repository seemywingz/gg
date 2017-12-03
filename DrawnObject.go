package in3D

import (
	"fmt"

	"github.com/bytearena/box2d"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// DrawnObject : a struct to hold openGL object data
type DrawnObject struct {
	Mesh           *Mesh
	Body           *box2d.B2Body
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

func initBodyPhisics(pos Position) *box2d.B2Body {
	if !Feature[Physics] {
		return nil
	}

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(float64(pos.X), float64(pos.Y))
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	// bd.FixedRotation = true
	bd.AllowSleep = false

	body := world.CreateBody(&bd)

	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(1, 1)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 20.0
	fd.Restitution = 0.2
	body.CreateFixtureFromDef(&fd)
	return body
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

func (d *DrawnObject) translateRotate() *mgl32.Mat4 {
	model := mgl32.Translate3D(d.X, d.Y, d.Z).
		Mul4(mgl32.Scale3D(d.Scale, d.Scale, d.Scale))
	xrotMatrix := mgl32.HomogRotate3DX(mgl32.DegToRad(d.XRotation))
	yrotMatrix := mgl32.HomogRotate3DY(mgl32.DegToRad(d.YRotation))
	zrotMatrix := mgl32.HomogRotate3DZ(mgl32.DegToRad(d.ZRotation))
	final := model.Mul4(xrotMatrix.Mul4(yrotMatrix.Mul4(zrotMatrix)))
	return &final
}

// Draw : draw the object
func (d *DrawnObject) Draw() {

	if Feature[Physics] {
		p := d.Body.GetPosition()
		d.Position.X = float32(p.X)
		d.Position.Y = float32(p.Y)

		a := d.Body.GetAngle()
		d.ZRotation = float32(a)

		fmt.Println(a)
	}

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
