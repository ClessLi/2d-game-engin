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
	camera     *Camera
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
//     camera: Camera 镜头类指针，用于 Scene 场景类绑定 camera 子类
//     init: Init 函数，用于场景 Create() 方法初始化时调用
// 返回值:
//     Scene 类指针
func NewScene(w, h float32, p *Player, sp *resolv.Space, camera *Camera, init func()) *Scene {
	sp.Add(p)
	return &Scene{
		Player:     p,
		Map:        sp,
		renderer:   nil,
		camera:     camera,
		Keys:       [1024]bool{},
		LockedKeys: [1024]bool{},
		Init:       init,
		W:          w,
		H:          h,
	}
}

// resetSceneSize, Scene 类重置场景边界的包内方法
// 参数:
//     width, height: 场景尺寸
func (gm *Scene) resetSceneSize(width, height float32) {
	gm.W = width
	gm.H = height
}

// Create, Scene 类场景创建方法
// TODO: 定义游戏场景接口，并将其作为接口方法实现
func (gm *Scene) Create() {
	//初始化着色器
	resource.LoadShader("../../resource/glsl/shader.vs", "../../resource/glsl/shader.fs", "sprite")
	shader := resource.GetShader("sprite")
	shader.Use()
	shader.SetInt("image", 0)
	//初始化精灵渲染器
	gm.renderer = render.NewSpriteRenderer(shader)
	//设置投影
	projection := mgl32.Ortho(0, gm.W, gm.H, 0, -1, 1)
	shader.SetMatrix4fv("projection", &projection[0])

	// 初始化地图
	gm.Init()
}

// Update, Scene 类场景更新方法
// TODO: 定义游戏场景接口，并将其作为接口方法实现
func (gm *Scene) Update(delta float64) {
	gm.Player.SpeedY += 0.5

	// 更新移动物体
	gm.updateMove()

	// Check for a collision downwards by just attempting a resolution downwards and seeing if it collides with something.
	down := gm.Map.Filter(func(shape resolv.Shape) bool {
		if (shape.HasTags("solid") || shape.HasTags("ramp")) && !shape.HasTags("destroyed") {
			return true
		}
		return false
	}).Resolve(gm.Player, 0, 4)
	onGround := down.Colliding()
	gm.Player.IsMove = false

	// 角色左右移动
	gm.playerMove(down)

	// JUMP
	gm.playerJump(onGround)

	// Attack
	gm.playerAttack(delta)

	if !gm.Player.IsMove {
		gm.Player.ToStand(float32(delta))
	} else {
		gm.Player.ToMove(float32(delta))
	}

	x := int32(gm.Player.SpeedX)
	y := int32(gm.Player.SpeedY)

	solids := gm.Map.FilterByTags("solid")
	ramps := gm.Map.FilterByTags("ramp")
	dangers := gm.Map.FilterByTags("dangerous")

	//fmt.Println("check player is dead or not.")
	// 判断用户是否已死亡
	if res := dangers.Resolve(gm.Player, x, y); res.Colliding() {
		//fmt.Println("player is dead.")
		gm.Player.AddTags("isDead")
	}

	// X-movement. We only want to collide with solid objects (not ramps) because we want to be able to move up them
	// and don't need to be inhibited on the x-axis when doing so.

	if res := solids.Resolve(gm.Player, x, 0); res.Colliding() {
		x = res.ResolveX
		gm.Player.SpeedX = 0
	}

	gm.Player.X += x

	// Y movement. We check for ramp collision first; if we find it, then we just automatically will
	// slide up the ramp because the player is moving into it.

	// We look for ramps a little aggressively downwards because when walking down them, we want to stick to them.
	// If we didn't do this, then you would "bob" when walking down the ramp as the Player moves too quickly out into
	// space for gravity to push back down onto the ramp.
	res := ramps.Resolve(gm.Player, 0, y+4)

	if y < 0 || (res.Teleporting && res.ResolveY < -gm.Player.H/2) {
		res = resolv.Collision{}
	}

	if !res.Colliding() {
		res = solids.Resolve(gm.Player, 0, y)
	}

	if res.Colliding() {
		y = res.ResolveY
		gm.Player.SpeedY = 0
	}

	gm.Player.Y += y

	//if gm.Player.HasTags("isDead") {
	//	gm.Player.SpeedX = 0
	//}

}

// Draw, Scene 类场景渲染方法
// TODO: 定义游戏场景接口，并将其作为接口方法实现
func (gm *Scene) Draw() {

	resource.GetShader("sprite").SetMatrix4fv("view", gm.camera.GetViewMatrix())
	// 若角色处于死亡状态，则调整角色 Texture 为死亡态的
	if gm.Player.HasTags("isDead") {
		gm.Player.Texture = resource.GetTexture("x")
	}
	gm.Player.Draw(gm.renderer)

	//摄像头跟随
	playerSize := gm.Player.GetSize()
	screenX := float32(gm.Player.X) - gm.camera.W/2 + playerSize[0]
	screenY := float32(gm.Player.Y) - gm.camera.H/2 + playerSize[1]
	gm.camera.InPosition(screenX, screenY, gm.W, gm.H)

	// TODO: 由于渲染依赖camera，暂时将space内各个对象渲染放在这个位置
	for _, shape := range *gm.Map {
		if shape != gm.Player && gm.isInCamera(shape) && !shape.HasTags("hide") && !shape.HasTags("destroyed") && !shape.HasTags("init") {
			shape.Draw(gm.renderer)
		}

		if shape.HasTags("destroy") {
			shape.RemoveTags("destroy")
			shape.AddTags("destroyed")
		}

	}

	//if gm.DrawHelpText {
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
func (gm *Scene) Destroy() {
	gm.Map.Clear()
}

// isInCamera, Scene 类判断 shape 对象是否在镜头内的包内方法
// 参数:
//     shape: resolv.Shape 接口对象
// 返回值:
//     bool 类型， true 为在镜头内， false 为在镜头外
func (gm *Scene) isInCamera(shape resolv.Shape) bool {
	position := gm.camera.GetPosition()
	x := int32(position.X())
	y := int32(position.Y())
	cameraRec := resolv.NewRectangle(x, y, int32(gm.camera.W), int32(gm.camera.H), 0, 0, nil, nil)
	return cameraRec.IsColliding(shape)
}

// SetKeyDown, Scene 类设置控制器按键按下的方法
// 参数:
//     key: glfw.Key 类，对应控制器按键
func (gm *Scene) SetKeyDown(key glfw.Key) {
	gm.Keys[key] = true
}

// IsPressed, Scene 类判断控制器按键是否处于按下状况的方法
// 参数:
//     keys: glfw.Key 类参数列表，对应查询按键
// 返回值:
//     bool 类型， true 为查询按键列表中存在处于按下状态的按键， false 为查询列表中不含按下状态的按键
func (gm *Scene) IsPressed(keys ...glfw.Key) bool {
	for _, key := range keys {
		if gm.LockedKeys[key] {
			return true
		}
	}
	return false
}

// PressedKey, Scene 类设置控制器按键为按下状态的方法
// 参数:
//     key: glfw.Key 类，对应控制器按键
func (gm *Scene) PressedKey(key glfw.Key) {
	gm.LockedKeys[key] = true
}

// ReleaseKey, Scene 类设置控制器按键释放并解除按下状态的方法
func (gm *Scene) ReleaseKey(key glfw.Key) {
	gm.Keys[key] = false
	gm.LockedKeys[key] = false
}

// HasOneKeyDown, Scene 类判断控制器按键列表中是否存在至少一个按键已按下的方法
// 参数:
//     keys: glfw.Key 类参数列表，对应查询按键
// 返回值:
//     bool 类型， true 为存在， false 为不存在
func (gm *Scene) HasOneKeyDown(keys ...glfw.Key) bool {
	for _, key := range keys {
		if gm.Keys[key] {
			return true
		}
	}
	return false
}

// playerJump, Scene 类玩家角色跳跃的包内方法
// 参数:
//     onGround: bool 类，角色是否着陆
func (gm *Scene) playerJump(onGround bool) {
	if gm.HasOneKeyDown(glfw.KeyUp, glfw.KeyW) && !gm.IsPressed(glfw.KeyUp, glfw.KeyW) && onGround && !gm.Player.HasTags("isDead") {
		gm.Player.IsMove = true
		// 现在跳跃按键按下后重复跳跃
		if gm.HasOneKeyDown(glfw.KeyUp) {
			gm.PressedKey(glfw.KeyUp)
		}
		if gm.HasOneKeyDown(glfw.KeyW) {
			gm.PressedKey(glfw.KeyW)
		}
		gm.Player.SpeedY = -16
	}
}

// playerMove, Scene 类玩家角色移动的包内方法
// 参数:
//     down: resolv.Collision 类，角色着陆点所在对象
func (gm *Scene) playerMove(down resolv.Collision) {
	onGround := down.Colliding()
	friction := float32(0.01)
	if onGround {
		ground := down.ShapeB
		if ground.GetFriction() <= gm.Player.GetFriction() {
			friction = ground.GetFriction()
		} else {
			friction = gm.Player.GetFriction()
		}
	}
	accel := gm.Player.GetFriction() + friction

	if gm.Player.SpeedX > friction {
		gm.Player.SpeedX -= friction
	} else if gm.Player.SpeedX < -friction {
		gm.Player.SpeedX += friction
	} else {
		gm.Player.SpeedX = 0
	}

	if gm.HasOneKeyDown(glfw.KeyLeft, glfw.KeyRight, glfw.KeyA, glfw.KeyD) {
		gm.Player.IsMove = true
	}

	if gm.HasOneKeyDown(glfw.KeyRight, glfw.KeyD) && onGround && !gm.Player.HasTags("isDead") {
		gm.Player.IsXReverse = false
		gm.Player.SpeedX += accel
	}

	if gm.HasOneKeyDown(glfw.KeyLeft, glfw.KeyA) && onGround && !gm.Player.HasTags("isDead") {
		gm.Player.IsXReverse = true
		gm.Player.SpeedX -= accel
	}

	//fmt.Println(gm.Player.SpeedX)
	if gm.Player.SpeedX > gm.Player.GetMaxSpd() {
		gm.Player.SpeedX = gm.Player.GetMaxSpd()
	}

	if gm.Player.SpeedX < -gm.Player.GetMaxSpd() {
		gm.Player.SpeedX = -gm.Player.GetMaxSpd()
	}
	//fmt.Println(gm.Player.SpeedX)
}

// playerAttack, Scene 类玩家角色攻击的包内方法
// 参数:
//      delta: float64 类型，与上次更新的时延度量
func (gm *Scene) playerAttack(delta float64) {
	gm.Player.weapon.CoolDown(delta)

	// 调整角色攻击矢量
	if gm.Player.IsXReverse {
		gm.Player.AtkVec[0] = -1
	} else {
		gm.Player.AtkVec[0] = 1
	}

	if gm.HasOneKeyDown(glfw.KeyUp, glfw.KeyW) {
		gm.Player.AtkVec[1] = -1
	} else if gm.HasOneKeyDown(glfw.KeyDown, glfw.KeyS) {
		gm.Player.AtkVec[1] = 1
	} else {
		gm.Player.AtkVec[1] = 0
	}

	if gm.HasOneKeyDown(glfw.KeyJ, glfw.KeySpace) && !gm.Player.HasTags("isDead") {
		bolt := gm.Player.Attack()
		if bolt != nil {
			gm.Map.Add(bolt)
		}
	}
}

// updateMove, Scene 类 Update() 方法调用，用于更新“移动物体”位置的包内方法
func (gm *Scene) updateMove() {
	move := gm.Map.FilterByTags("isMove")
	for i := 0; i < move.Length(); i++ {
		shape := move.Get(i)
		X, Y := shape.GetXY()
		x, y := shape.GetSpd()
		if res := gm.Map.FilterByTags("solid").Resolve(shape, int32(x), int32(y)); res.Colliding() {
			x = float32(res.ResolveX)
			y = float32(res.ResolveY)
			shape.SetSpd(x, y)
			shape.AddTags("destroy")
		}
		shape.SetXY(X+int32(x), Y+int32(y))
	}
}
