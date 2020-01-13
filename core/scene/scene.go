// scene 包，该包包含了多个关于场景的结构体，定义了相关结构体的方法及相关函数
// 创建人： ClessLi
// 创建时间： 2020-1-8
package scene

import (
	"github.com/ClessLi/2d-game-engin/core/render"
	"github.com/ClessLi/2d-game-engin/core/resolv"
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Scene struct {
	Player *Player
	Map    *resolv.Space
	//精灵渲染器
	renderer *render.SpriteRenderer
	//摄像头
	Camera     *Camera
	Keys       [1024]bool
	LockedKeys [1024]bool
	Init       func()
	W, H       float32
}

// NewScene, 初始化 Scene 类实例函数
// 参数:
//     w, h: 场景尺寸
//     p: Player 玩家角色类指针
//     sp: resolv.Space 空间集合类指针，用于定义场景地图
//     Camera: Camera 镜头类指针，用于 Scene 场景类绑定 Camera 子类
//     init: Init 函数，用于场景 Create() 方法初始化时调用
// 返回值:
//     Scene 类指针
func NewScene(sceneW, sceneH float32, p *Player, sp *resolv.Space, camera *Camera, init func()) *Scene {

	if p != nil && sp != nil {
		sp.Add(p)
	}

	return &Scene{
		Player:     p,
		Map:        sp,
		renderer:   nil,
		Camera:     camera,
		Keys:       [1024]bool{},
		LockedKeys: [1024]bool{},
		Init:       init,
		W:          sceneW,
		H:          sceneH,
	}
}

// resetSceneSize, Scene 类重置场景边界的包内方法
// 参数:
//     width, height: 场景尺寸
func (s *Scene) resetSceneSize(width, height float32) {
	s.W = width
	s.H = height
}

// Create, Scene 类场景创建方法
// TODO: 定义游戏场景接口，并将其作为接口方法实现
func (s *Scene) Create() {
	//str, _ := os.Getwd()
	//fmt.Println(str)
	//初始化着色器
	resource.LoadShader("resource/glsl/shader.vs", "resource/glsl/shader.fs", "sprite")
	shader := resource.GetShader("sprite")
	shader.Use()
	shader.SetInt("image", 0)
	//初始化精灵渲染器
	s.renderer = render.NewSpriteRenderer(shader)

	// 初始化地图
	s.Init()

	//设置投影
	// mgl32.Ortho(0, --投影宽度, --投影高度, 0, -1, 1)
	projection := mgl32.Ortho(0, s.Camera.W, s.Camera.H, 0, -1, 1)
	shader.SetMatrix4fv("projection", &projection[0])
}

// Update, Scene 类场景更新方法
// TODO: 定义游戏场景接口，并将其作为接口方法实现
func (s *Scene) Update(delta float64) {
	s.Player.SpeedY += 0.5

	// 更新移动物体
	s.updateMove()

	// Check for a collision downwards by just attempting a resolution downwards and seeing if it collides with something.
	down := s.Map.Filter(func(shape resolv.Shape) bool {
		if (shape.HasTags("solid") || shape.HasTags("ramp")) && !shape.HasTags("destroyed") {
			return true
		}
		return false
	}).Resolve(s.Player, 0, 4)
	onGround := down.Colliding()
	s.Player.IsMove = false

	// 角色左右移动
	s.playerMove(down)

	// JUMP
	s.playerJump(onGround)

	// Attack
	s.playerAttack(delta)

	if !s.Player.IsMove {
		s.Player.ToStand(float32(delta))
	} else {
		s.Player.ToMove(float32(delta))
	}

	x := int32(s.Player.SpeedX)
	y := int32(s.Player.SpeedY)

	solids := s.Map.FilterByTags("solid")
	ramps := s.Map.FilterByTags("ramp")
	dangers := s.Map.FilterByTags("dangerous")

	//fmt.Println("check player is dead or not.")
	// 判断用户是否已死亡
	if res := dangers.Resolve(s.Player, x, y); res.Colliding() {
		//fmt.Println("player is dead.")
		s.Player.AddTags("isDead")
	}

	// X-movement. We only want to collide with solid objects (not ramps) because we want to be able to move up them
	// and don't need to be inhibited on the x-axis when doing so.

	if res := solids.Resolve(s.Player, x, 0); res.Colliding() {
		x = res.ResolveX
		s.Player.SpeedX = 0
	}

	s.Player.X += x

	// Y movement. We check for ramp collision first; if we find it, then we just automatically will
	// slide up the ramp because the player is moving into it.

	// We look for ramps a little aggressively downwards because when walking down them, we want to stick to them.
	// If we didn't do this, then you would "bob" when walking down the ramp as the Player moves too quickly out into
	// space for gravity to push back down onto the ramp.
	res := ramps.Resolve(s.Player, 0, y+4)

	if y < 0 || (res.Teleporting && res.ResolveY < -s.Player.H/2) {
		res = resolv.Collision{}
	}

	if !res.Colliding() {
		res = solids.Resolve(s.Player, 0, y)
	}

	if res.Colliding() {
		y = res.ResolveY
		s.Player.SpeedY = 0
	}

	s.Player.Y += y

	//if s.Player.HasTags("isDead") {
	//	s.Player.SpeedX = 0
	//}

}

// Draw, Scene 类场景渲染方法
// TODO: 定义游戏场景接口，并将其作为接口方法实现
func (s *Scene) Draw() {

	resource.GetShader("sprite").SetMatrix4fv("view", s.Camera.GetViewMatrix())
	// 若角色处于死亡状态，则调整角色 Texture 为死亡态的
	if s.Player.HasTags("isDead") {
		s.Player.Texture = resource.GetTexture("x")
	}
	s.Player.Draw(s.renderer)

	//摄像头跟随
	//playerSize := s.Player.GetSize()
	Px, Py := s.Player.Center()
	//screenX := float32(s.Player.X) - s.Camera.W/2 + playerSize[0]
	//screenY := float32(s.Player.Y) - s.Camera.H/2 + playerSize[1]
	screenX := float32(Px) - s.Camera.W/2
	screenY := float32(Py) - s.Camera.H/2
	//fmt.Printf("3) Px: %d, Py: %d, sx: %f, sy: %f\n", Px, Py, s.Camera.X, s.Camera.Y)
	s.Camera.InPosition(screenX, screenY, s.W, s.H)
	//fmt.Printf("4) Px: %d, Py: %d, sx: %f, sy: %f\n", Px, Py, s.Camera.X, s.Camera.Y)

	// TODO: 由于渲染依赖camera，暂时将space内各个对象渲染放在这个位置
	for _, shape := range *s.Map {
		if shape != s.Player && s.isInCamera(shape) && !shape.HasTags("hide") && !shape.HasTags("destroyed") && !shape.HasTags("init") {
			shape.Draw(s.renderer)
		}

		if shape.HasTags("destroy") {
			shape.RemoveTags("destroy")
			shape.AddTags("destroyed")
		}

	}
	//fmt.Println(s.Player.X, s.Player.Y, s.Camera.X, s.Camera.Y, s.Camera.W, s.Camera.H)

	//if s.DrawHelpText {
	//    DrawText(32, 16,
	//        "-Platformer test-",
	//        "You are the green square.",
	//        "Use the arrow keys to move.",
	//        "Press X to playerJump.",
	//        "You can playerJump through blue ramps / platforms.")
	//}

}

// Destroy, Scene 类场景销毁方法
// TODO: 定义游戏场景接口，并将其作为接口方法实现
func (s *Scene) Destroy() {
	s.Map.Clear()
}

// isInCamera, Scene 类判断 shape 对象是否在镜头内的包内方法
// 参数:
//     shape: resolv.Shape 接口对象
// 返回值:
//     bool 类型， true 为在镜头内， false 为在镜头外
func (s *Scene) isInCamera(shape resolv.Shape) bool {
	//position := s.Camera.GetPosition()
	//x := int32(position.X())
	//y := int32(position.Y())
	//cameraRec := resolv.NewRectangle(x, y, int32(s.Camera.W), int32(s.Camera.H), 0, 0, nil, nil)
	cameraRec := resolv.NewRectangle(int32(s.Camera.X), int32(s.Camera.Y), int32(s.Camera.W), int32(s.Camera.H), 0, 0, nil, nil)
	return cameraRec.IsColliding(shape)
}

// SetKeyDown, Scene 类设置控制器按键按下的方法
// 参数:
//     key: glfw.Key 类，对应控制器按键
func (s *Scene) SetKeyDown(key glfw.Key) {
	s.Keys[key] = true
}

// IsPressed, Scene 类判断控制器按键是否处于按下状况的方法
// 参数:
//     keys: glfw.Key 类参数列表，对应查询按键
// 返回值:
//     bool 类型， true 为查询按键列表中存在处于按下状态的按键， false 为查询列表中不含按下状态的按键
func (s *Scene) IsPressed(keys ...glfw.Key) bool {
	for _, key := range keys {
		if s.LockedKeys[key] {
			return true
		}
	}
	return false
}

// PressedKey, Scene 类设置控制器按键为按下状态的方法
// 参数:
//     key: glfw.Key 类，对应控制器按键
func (s *Scene) PressedKey(key glfw.Key) {
	s.LockedKeys[key] = true
}

// ReleaseKey, Scene 类设置控制器按键释放并解除按下状态的方法
func (s *Scene) ReleaseKey(key glfw.Key) {
	s.Keys[key] = false
	s.LockedKeys[key] = false
}

// HasOneKeyDown, Scene 类判断控制器按键列表中是否存在至少一个按键已按下的方法
// 参数:
//     keys: glfw.Key 类参数列表，对应查询按键
// 返回值:
//     bool 类型， true 为存在， false 为不存在
func (s *Scene) HasOneKeyDown(keys ...glfw.Key) bool {
	for _, key := range keys {
		if s.Keys[key] {
			return true
		}
	}
	return false
}

// playerJump, Scene 类玩家角色跳跃的包内方法
// 参数:
//     onGround: bool 类，角色是否着陆
func (s *Scene) playerJump(onGround bool) {
	if s.HasOneKeyDown(glfw.KeyUp, glfw.KeyW) && !s.IsPressed(glfw.KeyUp, glfw.KeyW) && onGround && !s.Player.HasTags("isDead") {
		s.Player.IsMove = true
		// 现在跳跃按键按下后重复跳跃
		if s.HasOneKeyDown(glfw.KeyUp) {
			s.PressedKey(glfw.KeyUp)
		}
		if s.HasOneKeyDown(glfw.KeyW) {
			s.PressedKey(glfw.KeyW)
		}
		s.Player.SpeedY = -16
	}
}

// playerMove, Scene 类玩家角色移动的包内方法
// 参数:
//     down: resolv.Collision 类，角色着陆点所在对象
func (s *Scene) playerMove(down resolv.Collision) {
	onGround := down.Colliding()
	friction := float32(0.01)
	if onGround {
		ground := down.ShapeB
		if ground.GetFriction() <= s.Player.GetFriction() {
			friction = ground.GetFriction()
		} else {
			friction = s.Player.GetFriction()
		}
	}
	accel := s.Player.GetFriction() + friction

	if s.Player.SpeedX > friction {
		s.Player.SpeedX -= friction
	} else if s.Player.SpeedX < -friction {
		s.Player.SpeedX += friction
	} else {
		s.Player.SpeedX = 0
	}

	if s.HasOneKeyDown(glfw.KeyLeft, glfw.KeyRight, glfw.KeyA, glfw.KeyD) {
		s.Player.IsMove = true
	}

	if s.HasOneKeyDown(glfw.KeyRight, glfw.KeyD) && onGround && !s.Player.HasTags("isDead") {
		s.Player.IsXReverse = false
		s.Player.SpeedX += accel
	}

	if s.HasOneKeyDown(glfw.KeyLeft, glfw.KeyA) && onGround && !s.Player.HasTags("isDead") {
		s.Player.IsXReverse = true
		s.Player.SpeedX -= accel
	}

	//fmt.Println(s.Player.SpeedX)
	if s.Player.SpeedX > s.Player.GetMaxSpd() {
		s.Player.SpeedX = s.Player.GetMaxSpd()
	}

	if s.Player.SpeedX < -s.Player.GetMaxSpd() {
		s.Player.SpeedX = -s.Player.GetMaxSpd()
	}
	//fmt.Println(s.Player.SpeedX)
}

// playerAttack, Scene 类玩家角色攻击的包内方法
// 参数:
//      delta: float64 类型，与上次更新的时延度量
func (s *Scene) playerAttack(delta float64) {
	s.Player.Weapon.CoolDown(delta)

	// 调整角色攻击矢量
	if s.Player.IsXReverse {
		s.Player.AtkVec[0] = -1
	} else {
		s.Player.AtkVec[0] = 1
	}

	if s.HasOneKeyDown(glfw.KeyUp, glfw.KeyW) {
		s.Player.AtkVec[1] = -1
	} else if s.HasOneKeyDown(glfw.KeyDown, glfw.KeyS) {
		s.Player.AtkVec[1] = 1
	} else {
		s.Player.AtkVec[1] = 0
	}

	if s.HasOneKeyDown(glfw.KeyJ, glfw.KeySpace) && !s.Player.HasTags("isDead") {
		bolt := s.Player.Attack()
		if bolt != nil {
			s.Map.Add(bolt)
		}
	}
}

// updateMove, Scene 类 Update() 方法调用，用于更新“移动物体”位置的包内方法
func (s *Scene) updateMove() {
	move := s.Map.FilterByTags("isMove")
	for i := 0; i < move.Length(); i++ {
		shape := move.Get(i)
		X, Y := shape.GetXY()
		x, y := shape.GetSpd()
		if res := s.Map.FilterByTags("solid").Resolve(shape, int32(x), int32(y)); res.Colliding() {
			x = float32(res.ResolveX)
			y = float32(res.ResolveY)
			shape.SetSpd(x, y)
			shape.AddTags("destroy")
		}
		shape.SetXY(X+int32(x), Y+int32(y))
	}
}
