package service

import (
	"github.com/GoFurry/gofurry-game-backend/apps/game/dao"
	"github.com/GoFurry/gofurry-game-backend/apps/game/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/log"
	"github.com/GoFurry/gofurry-game-backend/common/util"
	"github.com/goccy/go-json"
)

type gameService struct{}

var gameSingleton = new(gameService)

func GetGameService() *gameService { return gameSingleton }

// 查询 weight 前 num 条游戏记录
func (s gameService) GetGameList(num string, lang string) (gameVo []models.GameRespVo, err common.GFError) {
	intNum, parseErr := util.String2Int(num)
	if parseErr != nil {
		return gameVo, common.NewServiceError("入参转换错误")
	}
	gameList, err := dao.GetGameDao().GetGameList(intNum)
	if err != nil {
		return
	}

	for _, v := range gameList {
		newGameVo := models.GameRespVo{
			ID:          util.Int642String(v.ID),
			CreateTime:  v.CreateTime,
			UpdateTime:  v.UpdateTime,
			ReleaseDate: v.ReleaseDate,
			Appid:       util.Int642String(v.Appid),
			Header:      v.Header,
		}
		switch lang {
		case "zh":
			newGameVo.Name = v.Name
			newGameVo.Info = v.Info
		case "en":
			newGameVo.Name = v.NameEn
			newGameVo.Info = v.InfoEn
		default:
			newGameVo.Name = v.Name
			newGameVo.Info = v.Info
		}
		jsonErr := json.Unmarshal([]byte(v.Developers), &newGameVo.Developers)
		if jsonErr != nil {
			log.Error(v.Name, " ([]byte(*v.Developers), &newGameVo.Developers) err: ", jsonErr)
		}
		jsonErr = json.Unmarshal([]byte(v.Publishers), &newGameVo.Publishers)
		if jsonErr != nil {
			log.Error(v.Name, " ([]byte(*v.Publishers), &newGameVo.Publishers) err: ", jsonErr)
		}
		if v.Resources != nil {
			jsonErr = json.Unmarshal([]byte(*v.Resources), &newGameVo.Resources)
			if jsonErr != nil {
				log.Error(v.Name, " ([]byte(*v.Resources), &newGameVo.Resources) err: ", jsonErr)
			}
		}
		if v.Groups != nil {
			jsonErr = json.Unmarshal([]byte(*v.Groups), &newGameVo.Groups)
			if jsonErr != nil {
				log.Error(v.Name, " ([]byte(*v.Groups), &newGameVo.Groups) err: ", jsonErr)
			}
		}
		if v.Links != nil {
			jsonErr = json.Unmarshal([]byte(*v.Links), &newGameVo.Links)
			if jsonErr != nil {
				log.Error(v.Name, " ([]byte(*v.Links), &newGameVo.Links) err: ", jsonErr)
			}
		}
		gameVo = append(gameVo, newGameVo)
	}
	return
}
