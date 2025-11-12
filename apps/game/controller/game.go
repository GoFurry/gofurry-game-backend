package controller

import (
	"github.com/GoFurry/gofurry-game-backend/apps/game/service"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/gofiber/fiber/v2"
)

type gameApi struct{}

var GameApi *gameApi

func init() {
	GameApi = &gameApi{}
}

// @Summary 获取所有游戏的记录
// @Schemes
// @Description 获取所有游戏记录
// @Tags Game
// @Accept json
// @Produce json
// @Param num query string true "请求数量"
// @Param lang query string true "语言"
// @Success 200 {object} []models.GameRespVo
// @Router /api/game/info/list [Get]
func (api *gameApi) GetGameList(c *fiber.Ctx) error {
	num := c.Query("num", "100")
	lang := c.Query("lang", "zh")
	data, err := service.GetGameService().GetGameList(num, lang)
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}

	return common.NewResponse(c).SuccessWithData(data)
}
