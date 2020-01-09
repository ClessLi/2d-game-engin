package scene

import (
	"fmt"
	"github.com/ClessLi/2d-game-engin/core/resolv"
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/mathgl/mgl32"
)

// Weapon, 武器接口对象，定义武器攻击与冷却方法
type Weapon interface {
	Attack(X, Y int32, vec2 mgl32.Vec2, isXReverse bool) resolv.Shape
	CoolDown(delta float64)
}

// LongRangeWeapon, 远程武器类，定义远程武器及其相关信息
type LongRangeWeapon struct {
	BoltName   string
	CD         float64
	CDDelta    float64
	Speed      float32
	BoltRadius int32
}

// Attack, LongRangeWeapon 类攻击方法， Weapon.Attack(X, Y int32, vec2 mgl32.Vec2, isXReverse bool) resolv.Shape 的实现
// 参数:
//     X, Y: 攻击初始坐标
//     vec2: 攻击矢量
//     isXReverse: 图像是水平镜像向后的
// 返回值:
//     resolv.Shape 类，攻击作用形状对象
func (lw *LongRangeWeapon) Attack(X, Y int32, vec2 mgl32.Vec2, isXReverse bool) resolv.Shape {
	if lw.CDDelta > 0 {
		return nil
	}
	lw.CDDelta = lw.CD

	SpdX, SpdY := lw.initSpd(vec2)

	bolt := resolv.NewCircle(X, Y,
		lw.BoltRadius,
		0, 1.2,
		nil,
		resource.GetTexturesByName(lw.BoltName))

	bolt.IsXReverse = isXReverse
	bolt.SetSpd(SpdX, SpdY)
	bolt.AddTags("isMove")
	fmt.Println("shooting, x:", bolt.X, "y:", bolt.Y, "spdX:", bolt.SpeedX, "spdY:", bolt.SpeedY)
	return bolt
}

// CoolDown, LongRangeWeapon 类攻击冷却方法， Weapon.CoolDown(delta float64) 的实现
// 参数:
//     delta: 上次更新后时延
func (lw *LongRangeWeapon) CoolDown(delta float64) {
	lw.CDDelta -= delta
}

// initSpd, LongRangeWeapon 类定义子弹初速度的包内方法
// 参数:
//     vec2: 攻击方向矢量
// 返回值:
//     float32, float32 类型，子弹水平与垂直方向速度值
func (lw *LongRangeWeapon) initSpd(vec2 mgl32.Vec2) (float32, float32) {
	vecLen := vec2.Len()
	return vec2[0] * lw.Speed / vecLen, vec2[1] * lw.Speed / vecLen
}

// NewFireBolt, 火球武器实例化函数
// 返回值:
//     LongRangeWeapon 类指针
func NewFireBolt() *LongRangeWeapon {
	return &LongRangeWeapon{
		BoltName:   "FireBolt",
		CD:         1.0,
		CDDelta:    0,
		Speed:      20,
		BoltRadius: 10,
	}
}
