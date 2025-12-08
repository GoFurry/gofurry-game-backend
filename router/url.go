package routers

import (
	game "github.com/GoFurry/gofurry-game-backend/apps/game/controller"
	recommend "github.com/GoFurry/gofurry-game-backend/apps/recommend/controller"
	"github.com/gofiber/fiber/v2"
)

/*
 * @Desc: 接口层
 * @author: 福狼
 * @version: v1.0.0
 */

func gameApi(g fiber.Router) {
	g.Get("/info/list", game.GameApi.GetGameList)       // 获取前 num 条游戏记录
	g.Get("/info/main", game.GameApi.GetGameMainList)   // 获取首页展示数据
	g.Get("/panel/main", game.GameApi.GetPanelMainList) // 获取首页面板数据
	g.Get("/update/latest", game.GameApi.GetUpdateNews) // 获取首页更新公告
}

func recommendApi(g fiber.Router) {
	// 基于内容的推荐（Content-based Filtering）
	// 优点: 存储小 速度快 无冷启动 无需用户行为数据
	// 缺点: 需要传入初始物品, 特征值永远为静态, 每次推荐相同
	// 实现重点: 余弦相似度 特征提取-独热编码
	g.Get("/game/CBF", recommend.RecommendApi.RecommendByCBF)     // 用 CBF 返回游戏记录
	g.Get("/game/random", recommend.RecommendApi.GetRandomGameID) // 返回一个随机的游戏记录ID

}
