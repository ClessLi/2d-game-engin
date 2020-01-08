package resolv

import (
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/ClessLi/2d-game-engin/src/render"
	"github.com/go-gl/mathgl/mgl32"
)

// Rectangle represents a rectangle
type Rectangle struct {
	MoveShape
	W, H int32
}

// NewRectangle creates a new Rectangle and returns a pointer to it.
func NewRectangle(x, y, w, h int32, friction, drawMulti float32, moveTextures, standTextures []*resource.Texture2D) *Rectangle {
	r := &Rectangle{
		MoveShape: *NewMoveShape(
			x, y,
			0,
			friction,
			drawMulti,
			moveTextures,
			standTextures),
		W: w, H: h,
	}
	return r
}

// IsColliding returns whether the Rectangle is colliding with the specified other Shape or not, including the other Shape
// being wholly contained within the Rectangle.
func (r *Rectangle) IsColliding(other Shape) bool {

	switch b := other.(type) {
	case *Rectangle:
		return r.X > b.X-r.W && r.Y > b.Y-r.H && r.X < b.X+b.W && r.Y < b.Y+b.H
	default:
		return b.IsColliding(r)
	}

}

// WouldBeColliding returns whether the Rectangle would be colliding with the other Shape if it were to move in the
// specified direction.
func (r *Rectangle) WouldBeColliding(other Shape, dx, dy int32) bool {
	r.X += dx
	r.Y += dy
	isColliding := r.IsColliding(other)
	r.X -= dx
	r.Y -= dy
	return isColliding
}

// Center returns the center point of the Rectangle.
func (r *Rectangle) Center() (int32, int32) {

	x := r.X + r.W/2
	y := r.Y + r.H/2

	return x, y

}

// GetBoundingCircle returns a circle that wholly contains the Rectangle.
func (r *Rectangle) GetBoundingCircle() *Circle {

	x, y := r.Center()
	c := NewCircle(x, y, Distance(x, y, r.X+r.W, r.Y), r.friction, r.multiple, r.moveTextures, r.standTextures)
	return c

}

// 获取第二点横纵坐标
func (r *Rectangle) GetXY2() (x, y int32) {
	return r.X + r.W, r.Y + r.H
}

// 渲染
func (r *Rectangle) Draw(renderer *render.SpriteRenderer) {
	size := &mgl32.Vec2{
		r.multiple * float32(r.W),
		r.multiple * float32(r.H),
	}
	position := &mgl32.Vec2{
		float32(r.X) - size[0]/2,
		float32(r.Y) - size[1]/2,
	}

	renderer.DrawSprite(r.texture, position, size, r.rotate, r.color, r.isXReverse)
}
