package scene

import "github.com/go-gl/mathgl/mgl32"

// Camera, 镜头对象，定义了镜头（屏幕）画面显示对象的基础信息
type Camera struct {
	X, Y        float32 // 坐标
	W, H        float32 // 场景尺寸、屏幕尺寸
	front, up   mgl32.Vec3
	movementSpd float32 // 镜头转移速度
}

// NewDefaultCamera, Camera 类默认实例初始化的函数
// 参数:
//     X, Y: 镜头坐标
//     W, H: 屏幕显示尺寸
// 返回值:
//     Camera 类指针
func NewDefaultCamera(X, Y, W, H float32) *Camera {
	return &Camera{
		X:           X,
		Y:           Y,
		W:           W,
		H:           H,
		front:       mgl32.Vec3{0, 0, -1},
		up:          mgl32.Vec3{0, 1, 0},
		movementSpd: float32(100),
	}
}

// GetPosition, Camera 类的获取坐标点方法
// 返回值:
//     mgl32.Vec3 类
func (c *Camera) GetPosition() mgl32.Vec3 {
	return mgl32.Vec3{c.X, c.Y, 0}
}

// GetViewMatrix, Camera 类获取view的方法
// 返回值:
//     float32 类型指针
func (c *Camera) GetViewMatrix() *float32 {
	target := c.GetPosition().Add(c.front)
	view := mgl32.LookAtV(c.GetPosition(), target, c.up)
	return &view[0]
}

// resetScreenSize, Camera 类重置屏幕边界的包内方法
// 参数:
//     width, height: 镜头尺寸
func (c *Camera) resetScreenSize(width, height float32) {
	c.W = width
	c.H = height
}

// InPosition, Camera 类根据坐标转换视野的方法
// 参数:
//     x, y: 需转换的坐标
//     W, H: 场景尺寸
func (c *Camera) InPosition(x, y, sceneW, sceneH float32) {
	if x <= 0 {
		c.X = 0
	} else if x+c.W > sceneW {
		c.X = sceneW - c.W
	} else {
		c.X = x
	}

	if y <= 0 {
		c.Y = 0
	} else if y+c.H > sceneH {
		c.Y = sceneH - c.H
	} else {
		c.Y = y
	}
}
