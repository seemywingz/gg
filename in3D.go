package in3D

import (
	"fmt"
	_ "image/jpeg" // include jpeg support
	_ "image/png"  // include png support
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Init : initializes glfw and returns a Window to use, then InitGL
func Init(width, height int, title string) {
	runtime.LockOSThread()

	EoE("Error Initializing GLFW", glfw.Init())
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var err error
	if width == 0 || height == 0 {
		mode := glfw.GetPrimaryMonitor().GetVideoMode()
		width = mode.Width
		height = mode.Height
		window, err = glfw.CreateWindow(width, height, title, glfw.GetPrimaryMonitor(), nil)
	} else {
		window, err = glfw.CreateWindow(width, height, title, nil, nil)
	}
	EoE("Error Creating GLFW Window", err)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetInputMode(glfw.StickyMouseButtonsMode, 1)
	TogglePointerLock()
	window.MakeContextCurrent()
	InitGL()
	InitFeatures()
	NewCamera()
	NewLightManager()
}

// InitGL : initialize GL setting and print version
func InitGL() {
	EoE("Error Initializing OpenGL", gl.Init())

	gl.Enable(gl.DEPTH_TEST)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	gl.DepthFunc(gl.LESS)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	println("OpenGL version", version)
	InitShaders()
}

// InitFeatures :
func InitFeatures() {
	Feature = make(map[int]bool)
	Feature[Look] = false
	Feature[Move] = false
	Feature[FlyMode] = false
	Feature[PointerLock] = false
}

// InitShaders :
func InitShaders() {
	Shader = make(map[string]uint32)
	SetRelPath("shaders")
	Shader["basic"] = NewShader("Vert.glsl", "basicFrag.glsl")
	Shader["color"] = NewShader("Vert.glsl", "colorFrag.glsl")
	Shader["texture"] = NewShader("Vert.glsl", "textureFrag.glsl")
	Shader["fixedLight"] = NewShader("Vert.glsl", "fixedLightFrag.glsl")
	Shader["phong"] = NewShader("Vert.glsl", "blinnPhongFrag.glsl")
	Shader["normalMap"] = NewShader("normalMapVert.glsl", "normalMapFrag.glsl")
	Shader["in3D"] = NewShader("in3dVert.glsl", "in3DFrag.glsl")
}

// MakeVAO : initializes and returns a vertex array from the points provided.
func MakeVAO(points []float32, program uint32) uint32 {

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(0))

	vertTexCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(vertTexCoordAttrib)
	gl.VertexAttribPointer(vertTexCoordAttrib, 2, gl.FLOAT, false, 11*4, gl.PtrOffset(3*4))

	vertNormalAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertNormal\x00")))
	gl.EnableVertexAttribArray(vertNormalAttrib)
	gl.VertexAttribPointer(vertNormalAttrib, 3, gl.FLOAT, true, 11*4, gl.PtrOffset(5*4))

	vertTangentAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTangent\x00")))
	gl.EnableVertexAttribArray(vertTangentAttrib)
	gl.VertexAttribPointer(vertTangentAttrib, 3, gl.FLOAT, true, 11*4, gl.PtrOffset(8*4))

	return vao
}

// CompileShader :
func CompileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		EoE("Failed to Compile Source ", fmt.Errorf("failed to compile %v: %v", source, log))
	}

	return shader
}

// CompileShaderFromFile : create gl shader from source string
func CompileShaderFromFile(sourceFile string, shaderType uint32) uint32 {

	source, err := ioutil.ReadFile(sourceFile)
	EoE("Error Reading Source File", err)

	return CompileShader(string(source)+"\x00", shaderType)
}

// NewShader : create GL shader program from provided GLSL source files
func NewShader(vertexShaderSourceFile, fragmentShaderSourceFile string) uint32 {

	vertexShader := CompileShaderFromFile(vertexShaderSourceFile, gl.VERTEX_SHADER)
	fragmentShader := CompileShaderFromFile(fragmentShaderSourceFile, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)

	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		EoE("Error Linking Shader Program", fmt.Errorf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return program
}

// ShouldClose : wraper for glfw
func ShouldClose() bool {
	return window.ShouldClose()
}

// SwapBuffers : wrapper for glfw
func SwapBuffers() {
	window.SwapBuffers()
}

// Update :
func Update() {
	camera.Update()
	lightManager.Update()
}

// GetCamera : return pointer to gg camera
func GetCamera() *Camera {
	return camera
}

// GetWindow : return pounter to gg window
func GetWindow() *glfw.Window {
	return window
}

// SetCameraPosition :
func SetCameraPosition(position Position) {
	camera.Position = position
}

// SetCameraSpeed :
func SetCameraSpeed(speed float32) {
	camera.Speed = speed
}

// SetClearColor :
func SetClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

// TogglePointerLock :
func TogglePointerLock() {
	fmt.Println("PointerLock Enabled:", Feature[PointerLock])
	if Feature[PointerLock] {
		x, y := window.GetCursorPos()
		camera.LastX = x
		camera.LastY = y
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else {
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}
}

// Enable :
func Enable(feature int, enabled bool) {

	Feature[feature] = enabled
	switch feature {
	case Look:
		x, y := window.GetCursorPos()
		camera.LastX = x
		camera.LastY = y
	case PointerLock:
		TogglePointerLock()
	case FlyMode:
		Feature[Look] = enabled
		Feature[Move] = enabled
	}
}
