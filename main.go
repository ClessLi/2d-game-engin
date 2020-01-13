package main

import (
	"github.com/ClessLi/2d-game-engin/resource/demo"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const (
	Width  = 800
	Height = 600
)

var (
	windowName = "Test Game"
	game       = demo.NewDemo(Width, Height)
	deltaTime  = 0.0
	lastFrame  = 0.0
)

func main() {
	runtime.LockOSThread()
	window := initGlfw()
	defer glfw.Terminate()
	initOpenGL()
	game.Create()

	for !window.ShouldClose() {
		currFrame := glfw.GetTime()
		deltaTime = currFrame - lastFrame
		lastFrame = currFrame

		glfw.PollEvents()
		game.Update(deltaTime)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		game.Draw()
		window.SwapBuffers()
	}
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Viewport(0, 0, Width, Height)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(Width, Height, windowName, nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(KeyCallback)

	window.MakeContextCurrent()
	return window
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		game.SetKeyDown(key)
	case glfw.Release:
		game.ReleaseKey(key)
	}
}
