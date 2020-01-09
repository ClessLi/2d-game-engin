package resolv

import (
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/mathgl/mgl32"
)

// MoveShape, 可移动的游戏对象，扩展基础形状对象
type MoveShape struct {
	BasicShape
	SpeedX      float32
	SpeedY      float32
	maxSpd      float32
	BounceFrame float32
	// 对象是否为移动状态
	IsMove bool
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

// NewMoveShape, MoveShape 类实例初始化函数
// 参数:
//     x, y: 形状对象坐标
//     rotate: 形状旋转角度
//     friction: 阻力值
//     multiple: 形状渲染缩放系数
//     moveTextures: 动态 Texture 对象分片
//     standTextures: 静态 Texture 对象分片
// 返回值:
//     MoveShape 类指针
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
		IsMove:        false,
		moveTextures:  moveTextures,
		standTextures: standTextures,
		standIndex:    0,
		standDelta:    0,
		moveIndex:     0,
		moveDelta:     0,
	}
}

// ToStand, MoveShape 类使恢复静止的方法
// 参数:
//     delta: 上次更新后时延
func (m *MoveShape) ToStand(delta float32) {
	if m.standIndex >= len(m.standTextures) {
		m.standIndex = 0
	}
	m.standDelta += delta
	if m.standDelta > 0.1 {
		m.standDelta = 0
		m.Texture = m.standTextures[m.standIndex]
		m.standIndex += 1
	}
}

// ToMove, MoveShape 类使运动的方法
// 参数:
//     delta: 上次更新后时延
func (m *MoveShape) ToMove(delta float32) {
	if m.moveIndex >= len(m.moveTextures) {
		m.moveIndex = 0
	}
	m.moveDelta += delta
	if m.moveDelta > 0.05 {
		m.moveDelta = 0
		m.Texture = m.moveTextures[m.moveIndex]
		m.moveIndex += 1
	}
}

// GetMaxSpd, MoveShape 类获取最大速度的方法， Shape.GetMaxSpd() float32 的实现
// 返回值:
//     float32 类型
func (m *MoveShape) GetMaxSpd() float32 {
	return m.maxSpd
}

// SetMaxSpd, MoveShape 类设置最大速度的方法， Shape.SetMaxSpd(float32) 的实现
// 参数:
//     spd: 速度值
func (m *MoveShape) SetMaxSpd(spd float32) {
	m.maxSpd = spd
}

// GetSpd, MoveShape 类获取移动速度的方法， Shape.GetSpd() (float32, float32) 的实现
// 返回值:
//     float32, float32 类型，水平与垂直方向的移动速度值
func (m *MoveShape) GetSpd() (float32, float32) {
	return m.SpeedX, m.SpeedY
}

// SetSpd, MoveShape 类设置移动速度的方法， Shape.SetSpd(float32, float32) 的实现
// 参数:
//     spdX, spdY: 水平与垂直方向的移动速度值
func (m *MoveShape) SetSpd(spdX, spdY float32) {
	m.SpeedX = spdX
	m.SpeedY = spdY
}
