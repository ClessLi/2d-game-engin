package scene

import (
	"fmt"
	"github.com/ClessLi/2d-game-engin/core/resolv"
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/mathgl/mgl32"
)

type Weapon interface {
	Attack(X, Y int32, vec2 mgl32.Vec2, isXReverse bool) resolv.Shape
	CoolDown(delta float64)
}

type LongRangeWeapon struct {
	BoltName   string
	CD         float64
	CDDelta    float64
	Speed      float32
	BoltRadius int32
}

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

func (lw *LongRangeWeapon) CoolDown(delta float64) {
	lw.CDDelta -= delta
}

func (lw *LongRangeWeapon) initSpd(vec2 mgl32.Vec2) (float32, float32) {
	vecLen := vec2.Len()
	return vec2[0] * lw.Speed / vecLen, vec2[1] * lw.Speed / vecLen
}

func NewFireBolt() *LongRangeWeapon {
	weapon := &LongRangeWeapon{
		BoltName:   "FireBolt",
		CD:         1.0,
		CDDelta:    0,
		Speed:      20,
		BoltRadius: 10,
	}
	return weapon
}
