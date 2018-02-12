package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const (
	vertexProgramramSource = `
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
		  frag_colour = vec4(1, 1, 1, 1.0);
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

	vertexProgram := compileShaderProgram(vertexProgramramSource, gl.VERTEX_SHADER)
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

func main() {
	runtime.LockOSThread()

	window := initWindow()
	defer glfw.Terminate()

	program := compileGLProgram()

	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		// Thought: what if we had one big 'quad' that covers the screen and set it
		// draw from a texture that we update based on an internal representation
		// that we maintain as a plain old array. Sort of recreating the mem-mapped
		// IO that we used back in the bad old days?
		// it'd be nice to avoid making a couple of vertex lists for coords and
		// colours if we it's possible. I feel like "map this memory to a picture is
		// probably easier to understand
		window.SwapBuffers()
	}
}
