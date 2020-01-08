package resolv

import (
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/ClessLi/2d-game-engin/src/render"
	"github.com/go-gl/mathgl/mgl32"
)

type Shape interface {
	IsColliding(Shape) bool
	WouldBeColliding(Shape, int32, int32) bool
	GetTags() []string
	ClearTags()
	AddTags(...string)
	RemoveTags(...string)
	HasTags(...string) bool
	GetData() interface{}
	SetData(interface{})
	GetXY() (int32, int32)
	GetXY2() (int32, int32)
	SetXY(int32, int32)
	Move(int32, int32)
	Draw(*render.SpriteRenderer)
	GetFriction() float32
	SetFriction(float32)
	GetMaxSpd() float32
	SetMaxSpd(float32)
	GetSpd() (float32, float32)
	SetSpd(float32, float32)
}

// BasicShape isn't to be used directly; it just has some basic functions and data, common to all structs that embed it, like
// position and tags. It is embedded in other Shapes.
type BasicShape struct {
	X, Y       int32
	tags       []string
	Data       interface{}
	texture    *resource.Texture2D
	rotate     float32
	color      *mgl32.Vec3
	isXReverse bool
	friction   float32
	multiple   float32
}

// GetTags returns a reference to the the string array representing the tags on the BasicShape.
func (b *BasicShape) GetTags() []string {
	return b.tags
}

// AddTags adds the specified tags to the BasicShape.
func (b *BasicShape) AddTags(tags ...string) {
	if b.tags == nil {
		b.tags = []string{}
	}
	b.tags = append(b.tags, tags...)
}

// RemoveTags removes the specified tags from the BasicShape.
func (b *BasicShape) RemoveTags(tags ...string) {

	for _, t := range tags {

		for i := len(b.tags) - 1; i >= 0; i-- {

			if t == b.tags[i] {
				b.tags = append(b.tags[:i], b.tags[i+1:]...)
			}

		}

	}

}

// ClearTags clears the tags active on the BasicShape.
func (b *BasicShape) ClearTags() {
	b.tags = []string{}
}

// HasTags returns true if the Shape has all of the tags provided.
func (b *BasicShape) HasTags(tags ...string) bool {

	hasTags := true

	for _, t1 := range tags {
		found := false
		for _, shapeTag := range b.tags {
			if t1 == shapeTag {
				found = true
				continue
			}
		}
		if !found {
			hasTags = false
			break
		}
	}

	return hasTags
}

// GetData returns the data on the Shape.
func (b *BasicShape) GetData() interface{} {
	return b.Data
}

// SetData sets the data on the Shape.
func (b *BasicShape) SetData(data interface{}) {
	b.Data = data
}

// GetXY returns the position of the Shape.
func (b *BasicShape) GetXY() (int32, int32) {
	return b.X, b.Y
}

// SetXY sets the position of the Shape.
func (b *BasicShape) SetXY(x, y int32) {
	b.X = x
	b.Y = y
}

// Move moves the Shape by the delta X and Y values provided.
func (b *BasicShape) Move(x, y int32) {
	b.X += x
	b.Y += y
}

// 使对象转换为水平镜像
func (b *BasicShape) ReverseX() {
	b.isXReverse = true
}

// 使对象转换为非水平镜像
func (b *BasicShape) ForWordX() {
	b.isXReverse = false
}

// 获取 friction
func (b *BasicShape) GetFriction() float32 {
	return b.friction
}

// 设置 friction
func (b *BasicShape) SetFriction(friction float32) {
	b.friction = friction
}

// 实例化 BasicShape
func NewBasicShape(x, y int32, texture *resource.Texture2D, rotate float32, color *mgl32.Vec3, friction, multiple float32) *BasicShape {
	return &BasicShape{
		X:          x,
		Y:          y,
		tags:       nil,
		Data:       nil,
		texture:    texture,
		rotate:     rotate,
		color:      color,
		isXReverse: false,
		friction:   friction,
		multiple:   multiple,
	}
}
