package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
)

const (
	// These should come into distinct files, and then use a loader to pull
	// them in to be used
	vertexProgramSource = `
		#version 410
		in vec3 vp;
		void main() {
		  gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentProgramSource = `
		#version 410
		out vec4 frag_colour;
		void main() {
		  frag_colour = vec4(1.0, 0.0, 0.0, 1.0);
		}
	` + "\x00"
)

func initWindow() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(640, 320, "Chip-8", nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(windowKeyCallback)
	window.MakeContextCurrent()

	return window
}

func windowKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, modifierKeys glfw.ModifierKey) {
	// TODO: handle keyboard stuff 0-9 a-f and maybe esc for quitting
}

func compileGLProgram() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexProgram := compileShaderProgram(vertexProgramSource, gl.VERTEX_SHADER)
	fragmentProgram := compileShaderProgram(fragmentProgramSource, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexProgram)
	gl.AttachShader(program, fragmentProgram)

	gl.LinkProgram(program)
	return program
}

func compileShaderProgram(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)

	sources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, sources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		panic("failed to compile shader")
	}

	return shader
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

var (
	// pulling this out into a file might be worth while too
	quad = []float32{
		-1.0, 1.0, 0.0,
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		-1.0, 1.0, 0.0,
		1.0, -1.0, 0.0,
		1.0, 1.0, 0.0,
	}
)

func main() {
	log.Println("Starting up…")

	runtime.LockOSThread()

	window := initWindow()
	defer glfw.Terminate()

	program := compileGLProgram()
	vao := makeVao(quad)

	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))
		// render 'video memory' to a PBO, and use that to texture the quad
		// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl
		window.SwapBuffers()
	}

	log.Println("Shutting down…")
}
