package core

import (
	"github.com/ClessLi/2d-game-engin/core/resolv"
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	resolv.Rectangle
	weapon Weapon
	AtkVec mgl32.Vec2
}

func NewPlayer(x, y, w, h int32, friction, drawMulti float32, moveList, standList []string) *Player {
	r := resolv.NewRectangle(x, y, w, h,
		friction, drawMulti,
		resource.GetTexturesByName(moveList...),
		resource.GetTexturesByName(standList...))
	p := &Player{
		Rectangle: *r,
		weapon:    nil,
	}
	return p
}

func (p *Player) Attack() resolv.Shape {
	x, y := p.GetXY()
	if !p.IsXReverse {
		x += p.W * 2 / 3
	} else {
		x += p.W / 3
	}
	y += p.H / 4

	switch p.weapon.(type) {
	case *LongRangeWeapon:
		return p.weapon.Attack(x, y, p.AtkVec, p.IsXReverse)
	}
	return nil
}
