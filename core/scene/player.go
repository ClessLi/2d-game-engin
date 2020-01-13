package scene

import (
	"github.com/ClessLi/2d-game-engin/core/resolv"
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/mathgl/mgl32"
)

// Player, 玩家角色对象，暂以方形作为角色的形状对象，包含了武器类型与攻击矢量
type Player struct {
	resolv.Rectangle
	Weapon Weapon
	AtkVec mgl32.Vec2
}

// NewPlayer, Player 类实例初始化函数
// 参数:
//     x, y: 角色坐标
//     w, h: 角色长宽尺寸
//     friction: 角色阻力值
//     drawMulti: 角色渲染缩放系数
//     moveList: 动态 Texture 对象名列表
//     standList: 静态 Texture 对象名列表
// 返回值:
//     Player 类指针
func NewPlayer(x, y, w, h int32, friction, drawMulti float32, moveList, standList []string) *Player {
	r := resolv.NewRectangle(x, y, w, h,
		friction, drawMulti,
		resource.GetTexturesByName(moveList...),
		resource.GetTexturesByName(standList...))
	p := &Player{
		Rectangle: *r,
		Weapon:    nil,
	}
	return p
}

// Attack, Player 类攻击方法
// 返回值:
//     resolve.Shape 类，武器攻击作用形状对象
func (p *Player) Attack() resolv.Shape {
	x, y := p.GetXY()
	if !p.IsXReverse {
		x += p.W * 2 / 3
	} else {
		x += p.W / 3
	}
	y += p.H / 4

	switch p.Weapon.(type) {
	case *LongRangeWeapon:
		return p.Weapon.Attack(x, y, p.AtkVec, p.IsXReverse)
	}
	return nil
}
