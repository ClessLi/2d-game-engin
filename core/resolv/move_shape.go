package resolv

import (
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/mathgl/mgl32"
)

// 可移动的游戏对象
type MoveShape struct {
	BasicShape
	SpeedX      float32
	SpeedY      float32
	maxSpd      float32
	BounceFrame float32
	// 对象是否为移动状态
	isMove bool
	// 移动时的动画纹理
	moveTextures []*resource.Texture2D
	// 静止时的动画纹理
	standTextures []*resource.Texture2D
	//当前静止帧
	standIndex int
	//静止帧之间的切换阈值
	standDelta float32
	//当前运动帧
	moveIndex int
	//运动帧之间的切换阈值
	moveDelta float32
}

func NewMoveShape(x, y int32, rotate, friction, multiple float32, moveTextures []*resource.Texture2D, standTextures []*resource.Texture2D) *MoveShape {
	var texture *resource.Texture2D
	if len(standTextures) > 0 {
		texture = standTextures[0]
	} else if len(moveTextures) > 0 {
		texture = moveTextures[0]
	} else {
	}

	bs := NewBasicShape(x, y, texture, rotate, &mgl32.Vec3{1, 1, 1}, friction, multiple)

	return &MoveShape{
		BasicShape:    *bs,
		isMove:        false,
		moveTextures:  moveTextures,
		standTextures: standTextures,
		standIndex:    0,
		standDelta:    0,
		moveIndex:     0,
		moveDelta:     0,
	}
}

//恢复静止
func (m *MoveShape) ToStand(delta float32) {
	if m.standIndex >= len(m.standTextures) {
		m.standIndex = 0
	}
	m.standDelta += delta
	if m.standDelta > 0.1 {
		m.standDelta = 0
		m.texture = m.standTextures[m.standIndex]
		m.standIndex += 1
	}
}

//由用户主动发起的运动
func (m *MoveShape) ToMove(delta float32) {
	if m.moveIndex >= len(m.moveTextures) {
		m.moveIndex = 0
	}
	m.moveDelta += delta
	if m.moveDelta > 0.05 {
		m.moveDelta = 0
		m.texture = m.moveTextures[m.moveIndex]
		m.moveIndex += 1
	}
}

// 获取最大速度
func (m *MoveShape) GetMaxSpd() float32 {
	return m.maxSpd
}

// 设置最大速度
func (m *MoveShape) SetMaxSpd(spd float32) {
	m.maxSpd = spd
}

// 获取移动速度
func (m *MoveShape) GetSpd() (float32, float32) {
	return m.SpeedX, m.SpeedY
}

// 设置移动速度
func (m *MoveShape) SetSpd(spdX, spdY float32) {
	m.SpeedX = spdX
	m.SpeedY = spdY
}
