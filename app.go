package main

import (
	"./packages/env"
	"./packages/shader"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
	"syscall"
)

const defaultWindowWidth = 640
const defaultWindowHeight = 480

func init() {
	runtime.LockOSThread()
}

func dropCallback(w *glfw.Window, files []string) {
	for index, file := range files {
		fmt.Println(index, file)
	}
}

type Shader struct {
	vertex      string
	fragment    string
	program     uint32
	last_update int64
}

func newShader(vertex, fragment string) *Shader {
	sh := &Shader{vertex: vertex, fragment: fragment}
	sh.setProgram()
	sh.setLastUpdate()
	return sh
}

func (self *Shader) setProgram() {
	program, err := shader.LoadShader(self.vertex, self.fragment)
	self.program = program
	if err != nil {
		panic(err)
	}
}

func (self *Shader) setLastUpdate() {
	var s syscall.Stat_t
	syscall.Stat(self.fragment, &s)
	sec, _ := s.Mtimespec.Unix()
	self.last_update = sec
}

func (self *Shader) getLastUpdate() int64 {
	var s syscall.Stat_t
	syscall.Stat(self.fragment, &s)
	sec, _ := s.Mtimespec.Unix()
	return sec
}

func (self *Shader) changeShader(fragment string) {
	self.fragment = fragment
	self.setProgram()
	self.setLastUpdate()
}

func (self *Shader) useProgram() {
	gl.UseProgram(self.program)
}

func (self *Shader) checkUpdate() {
	if self.last_update != self.getLastUpdate() {
		self.setProgram()
		self.setLastUpdate()
	}
}

func (self *Shader) dropCallback(w *glfw.Window, files []string) {
	for _, file := range files {
		if program, err := shader.LoadShader(self.vertex, file); err == nil {
			self.fragment = file
			self.program = program
			self.setLastUpdate()
		}
	}
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

	sh := newShader(env.ResourcePath+"glsl/default.vs", env.ResourcePath+"glsl/default.fs")
	window.SetDropCallback(sh.dropCallback)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		sh.checkUpdate()
		sh.useProgram()

		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

		gl.DisableVertexAttribArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
