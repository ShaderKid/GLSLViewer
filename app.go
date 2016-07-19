package main

import (
	//"fmt"
	"./packages/shader"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const defaultWindowWidth = 640
const defaultWindowHeight = 480

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(defaultWindowWidth, defaultWindowHeight, "GLSLViewer", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexBufferData := []float32{
		-1.0, -1.0, 0.0,
		-1.0, 1.0, 0.0,
		1.0, -1.0, 0.0,
		1.0, 1.0, 0.0,
	}

	var vertexArrayID uint32
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBufferData)*4, gl.Ptr(vertexBufferData), gl.STATIC_DRAW)

	program, err := shader.LoadShader("./glsl/default.vs", "./glsl/default.fs")
	if err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

		gl.DisableVertexAttribArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
