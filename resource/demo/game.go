package demo

import (
	"github.com/ClessLi/2d-game-engin/core/resolv"
	"github.com/ClessLi/2d-game-engin/core/scene"
	"github.com/ClessLi/2d-game-engin/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// NewDemo, 游戏 demo 版框架初始化函数
// 参数:
//     w, h: 传入显示分辨率
// 返回值:
//     scene.Scene 类指针
func NewDemo(w, h float32) *scene.Scene {
	var (
		sceneW  float32 = 1600
		sceneH  float32 = 800
		screenW float32 = 800
		screenH float32 = 600
		cellW   float32 = 16
		cellH   float32 = 16
	)

	xF := w / screenW
	yF := h / screenH
	//fmt.Printf("1) xF: %f, yF: %f, sW: %f, sH: %f\n", xF, yF, screenW, screenH)

	sceneW = xF * sceneW
	sceneH = yF * sceneH
	screenW = xF * screenW
	screenH = yF * screenH
	cellW = xF * cellW
	cellH = yF * cellH
	//fmt.Printf("2) xF: %f, yF: %f, sW: %f, sH: %f\n", xF, yF, screenW, screenH)

	game := scene.NewScene(
		sceneW,
		sceneH,
		nil,
		nil,
		nil,
		nil)

	// 定义game.Init函数
	game.Init = func() {
		//加载资源
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/platformLine.png", "platformLine")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/firebolt.png", "FireBolt")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/line.png", "line")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/spike.png", "spike")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/wall.png", "wall")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/x.png", "x")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/0.png", "0")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/1.png", "1")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/2.png", "2")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/3.png", "3")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/4.png", "4")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/5.png", "5")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/6.png", "6")
		resource.LoadTexture(gl.TEXTURE0, "./resource/image/bat/7.png", "7")

		game.Map = resolv.NewSpace()
		game.Map.Clear()

		// 创建游戏角色
		player := scene.NewPlayer(
			int32(sceneW)/2, int32(sceneH)/2,
			100, 100,
			0.5,
			2,
			[]string{"0", "1", "2", "3", "4", "5", "6", "7"},
			[]string{"0", "1", "2", "3", "4", "5", "6", "7"})
		player.SetMaxSpd(5)
		player.Weapon = scene.NewFireBolt()
		game.Camera = scene.NewDefaultCamera(float32(player.X), float32(player.Y), screenW, screenH)
		game.Player = player
		game.Map.Add(player)

		// A ramp
		line := resolv.NewLine(
			int32(game.W/4+cellW),
			int32(game.H-cellH*4),
			int32(game.W/4+cellW*11),
			int32(game.H-cellH*10),
			0.5,
			1,
			resource.GetTexturesByName("line"),
			nil)
		line.AddTags("ramp")
		game.Map.Add(line)

		line = resolv.NewLine(
			int32(game.W/4+cellW*11),
			int32(game.H-cellH*10),
			int32(game.W/4+cellW*40),
			int32(game.H-cellH*10),
			0.5,
			1,
			resource.GetTexturesByName("line"),
			nil)
		line.AddTags("ramp")
		game.Map.Add(line)

		line = resolv.NewLine(
			int32(game.W/4+cellW*40),
			int32(game.H-cellH*10),
			int32(game.W/4+cellW*50),
			int32(game.H-cellH*4),
			0.5,
			1,
			resource.GetTexturesByName("line"),
			nil)
		line.AddTags("ramp")
		game.Map.Add(line)

		// 来点阻碍的线段
		line = resolv.NewLine(
			int32(game.W/4+cellW*10),
			int32(game.H-cellH*25),
			int32(game.W/4+cellW*10),
			int32(game.H-cellH*20),
			0.5,
			1,
			resource.GetTexturesByName("line"),
			nil)
		line.AddTags("ramp")
		game.Map.Add(line)

		for y := float32(0); y < game.H; y += cellH {

			for x := float32(0); x < game.W; x += cellW {

				// 构建四周的墙
				if y <= cellH*4 || y >= game.H-cellH*4 || x <= cellW*4 || x >= game.W-cellW*4 {
					wall := resolv.NewRectangle(
						int32(x),
						int32(y),
						int32(cellW),
						int32(cellH),
						0.5,
						1,
						nil,
						resource.GetTexturesByName("wall"))
					wall.AddTags("isWall", "solid", "ramp")
					game.Map.Add(wall)

				}

				// 构建顶部尖刺
				if y == cellH*5 && x > cellW*4 && x < game.W-cellW*4 {
					spike := resolv.NewRectangle(
						int32(x),
						int32(y),
						int32(cellW),
						int32(cellH),
						0.01,
						1,
						nil,
						resource.GetTexturesByName("spike"))
					spike.AddTags("dangerous", "isSpike")
					game.Map.Add(spike)
				}

			}

		}
	}

	return game
}
