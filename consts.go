package in3D

import "github.com/go-gl/glfw/v3.2/glfw"

const (
	// MaxLights :
	MaxLights = 10
	// Press :
	Press = glfw.Press
	// Releasese :
	Releasese = glfw.Release

	// TimeStep :
	TimeStep = 1.0 / 60.0
	// VelocityIterations :
	VelocityIterations = 8
	// PositionIterations :
	PositionIterations = 3
)

const (
	_ = iota
	// FlyMode : Allow "Flying" Through Scene
	FlyMode
	// PointerLock :
	PointerLock
	// Look :
	Look
	// Move :
	Move
	// NoTexture :
	NoTexture
	// Physics :
	Physics
)
