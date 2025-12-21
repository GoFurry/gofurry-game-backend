package controller

import (
	"github.com/GoFurry/gofurry-game-backend/apps/search/models"
	"github.com/GoFurry/gofurry-game-backend/apps/search/service"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/gofiber/fiber/v2"
)

type searchApi struct{}

var SearchApi *searchApi

func init() {
	SearchApi = &searchApi{}
}

// @Summary 简易搜索
// @Schemes
// @Description 简易搜索
// @Tags Search
// @Accept json
// @Produce json
// @Param body body models.SearchRequest true "请求body"
// @Success 200 {object} models.SearchGameVo
// @Router /api/search/game/simple [POST]
func (api *searchApi) SimpleSearch(c *fiber.Ctx) error {
	req := models.SearchRequest{}
	if err := c.BodyParser(&req); err != nil {
		return common.NewResponse(c).Error("解析请求体失败")
	}
	data, err := service.GetSearchService().SimpleSearchQuery(req)
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}

	return common.NewResponse(c).SuccessWithData(data)
}
