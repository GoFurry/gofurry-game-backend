package controller

import (
	"github.com/GoFurry/gofurry-game-backend/apps/system/user/models"
	"github.com/GoFurry/gofurry-game-backend/apps/system/user/service"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/gofiber/fiber/v2"
)

type userApi struct{}

var UserApi *userApi

func init() {
	UserApi = &userApi{}
}

// @Summary 登录请求
// @Schemes
// @Description 登录
// @Tags System-user
// @Accept json
// @Produce json
// @Param body body models.UserLoginRequest true "请求body"
// @Success 200 {object} common.ResultData
// @Router /api/system/user/login [Post]
func (api *userApi) Login(c *fiber.Ctx) error {
	var req models.UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return common.NewResponse(c).Error("参数错误: " + err.Error())
	}

	token, err := service.GetUserService().Login(req)
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}
	var data = map[string]interface{}{
		"token": token,
	}
	return common.NewResponse(c).SuccessWithData(data)
}
