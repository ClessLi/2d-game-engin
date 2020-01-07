package render

import (
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type SpriteRenderer struct {
	shader *resource.Shader
	vao    uint32
}

func NewSpriteRenderer(shader *resource.Shader) *SpriteRenderer {
	spriteRenderer := SpriteRenderer{
		shader: shader,
	}
	spriteRenderer.initRenderData()
	return &spriteRenderer
}

func (sr *SpriteRenderer) DrawSprite(texture *resource.Texture2D, position *mgl32.Vec2, size *mgl32.Vec2, rotate float32, color *mgl32.Vec3, isReverseX bool) {
	model := mgl32.Translate3D(position[0], position[1], 0).Mul4(mgl32.Translate3D(0.5*size[0], 0.5*size[1], 0))
	model = model.Mul4(mgl32.HomogRotate3D(rotate, mgl32.Vec3{0, 0, 1}))
	model = model.Mul4(mgl32.Translate3D(-0.5*size[0], -0.5*size[1], 0))
	model = model.Mul4(mgl32.Scale3D(size[0], size[1], 1))

	sr.shader.SetMatrix4fv("model", &model[0])
	if isReverseX {
		sr.shader.SetInt("reverseX", -1)
	} else {
		sr.shader.SetInt("reverseX", 1)
	}
	sr.shader.SetVector3f("spriteColor", *color)
	texture.Use()

	gl.BindVertexArray(sr.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}

func (sr *SpriteRenderer) initRenderData() {
	var vbo uint32
	vertices := []float32{
		0.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0,

		0.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
	}
	gl.GenVertexArrays(1, &sr.vao)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(sr.vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}
