// resolv 包，该包包含了项目最基础的形状相关对象，及相关方法与函数
// 创建者: SolarLune
// 重构者: ClessLi
// 创建时间: Sep 15, 2018
// 重构时间: 2020-1-8
package resolv

import (
	"github.com/ClessLi/2d-game-engin/core/render"
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/mathgl/mgl32"
)

// Shape is a basic interface that describes a Shape that can be passed to collision testing and resolution functions and
// exist in the same Space.
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
	Texture    *resource.Texture2D
	rotate     float32
	color      *mgl32.Vec3
	IsXReverse bool
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

// ReverseX, BasicShape 类方向转换为水平向后的方法
func (b *BasicShape) ReverseX() {
	b.IsXReverse = true
}

// ForWordX, BasicShape 类方向转换为水平向前的方法
func (b *BasicShape) ForWordX() {
	b.IsXReverse = false
}

// GetFriction, BasicShape 类获取 friction 的方法， Shape.GetFriction() float32 的实现
// 返回值:
//     float32 类型
func (b *BasicShape) GetFriction() float32 {
	return b.friction
}

// SetFriction, BasicShape 类设置 friction 的方法， Shape.SetFriction(float32) 的实现
// 参数:
//     friction: 阻力值
func (b *BasicShape) SetFriction(friction float32) {
	b.friction = friction
}

// NewBasicShape, BasicShape 类的实例初始化函数
func NewBasicShape(x, y int32, texture *resource.Texture2D, rotate float32, color *mgl32.Vec3, friction, multiple float32) *BasicShape {
	return &BasicShape{
		X:          x,
		Y:          y,
		tags:       nil,
		Data:       nil,
		Texture:    texture,
		rotate:     rotate,
		color:      color,
		IsXReverse: false,
		friction:   friction,
		multiple:   multiple,
	}
}
