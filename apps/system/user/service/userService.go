package service

import (
	"github.com/GoFurry/gofurry-game-backend/apps/system/user/models"
	"github.com/GoFurry/gofurry-game-backend/common"
)

type userService struct{}

var userSingleton = new(userService)

func GetUserService() *userService { return userSingleton }

// 用户登录
func (svc *userService) Login(req models.UserLoginRequest) (tokenStr string, err common.GFError) {
	return
}
