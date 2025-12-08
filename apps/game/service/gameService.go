package service

import (
	"github.com/GoFurry/gofurry-game-backend/apps/game/dao"
	"github.com/GoFurry/gofurry-game-backend/apps/game/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/log"
	cs "github.com/GoFurry/gofurry-game-backend/common/service"
	"github.com/GoFurry/gofurry-game-backend/common/util"
	"github.com/bytedance/sonic"
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
		jsonErr := sonic.Unmarshal([]byte(v.Developers), &newGameVo.Developers)
		if jsonErr != nil {
			log.Error(v.Name, " ([]byte(*v.Developers), &newGameVo.Developers) err: ", jsonErr)
		}
		jsonErr = sonic.Unmarshal([]byte(v.Publishers), &newGameVo.Publishers)
		if jsonErr != nil {
			log.Error(v.Name, " ([]byte(*v.Publishers), &newGameVo.Publishers) err: ", jsonErr)
		}
		if v.Resources != nil {
			jsonErr = sonic.Unmarshal([]byte(*v.Resources), &newGameVo.Resources)
			if jsonErr != nil {
				log.Error(v.Name, " ([]byte(*v.Resources), &newGameVo.Resources) err: ", jsonErr)
			}
		}
		if v.Groups != nil {
			jsonErr = sonic.Unmarshal([]byte(*v.Groups), &newGameVo.Groups)
			if jsonErr != nil {
				log.Error(v.Name, " ([]byte(*v.Groups), &newGameVo.Groups) err: ", jsonErr)
			}
		}
		if v.Links != nil {
			jsonErr = sonic.Unmarshal([]byte(*v.Links), &newGameVo.Links)
			if jsonErr != nil {
				log.Error(v.Name, " ([]byte(*v.Links), &newGameVo.Links) err: ", jsonErr)
			}
		}
		gameVo = append(gameVo, newGameVo)
	}
	return
}

func (s gameService) GetGameMainList() (res models.GameMainInfoVo, err common.GFError) {
	jsonStr, err := cs.GetString("game-info:latest")
	if err != nil {
		return res, err
	}
	jsonErr := sonic.Unmarshal([]byte(jsonStr), &res.Latest)
	if jsonErr != nil {
		return res, common.NewServiceError(err.GetMsg())
	}

	jsonStr, err = cs.GetString("game-info:recent")
	if err != nil {
		return res, err
	}
	jsonErr = sonic.Unmarshal([]byte(jsonStr), &res.Recent)
	if jsonErr != nil {
		return res, common.NewServiceError(err.GetMsg())
	}

	jsonStr, err = cs.GetString("game-info:hot")
	if err != nil {
		return res, err
	}
	jsonErr = sonic.Unmarshal([]byte(jsonStr), &res.Hot)
	if jsonErr != nil {
		return res, common.NewServiceError(err.GetMsg())
	}

	jsonStr, err = cs.GetString("game-info:free")
	if err != nil {
		return res, err
	}
	jsonErr = sonic.Unmarshal([]byte(jsonStr), &res.Free)
	if jsonErr != nil {
		return res, common.NewServiceError(err.GetMsg())
	}
	return
}

func (s gameService) GetPanelMainList() (res models.GameMainPanelVo, err common.GFError) {
	jsonStr, err := cs.GetString("game-panel:top-player-count")
	if err != nil {
		return res, err
	}
	jsonErr := sonic.Unmarshal([]byte(jsonStr), &res.CountVo)
	if jsonErr != nil {
		return res, common.NewServiceError(err.GetMsg())
	}

	jsonStr, err = cs.GetString("game-panel:top-price")
	if err != nil {
		return res, err
	}
	jsonErr = sonic.Unmarshal([]byte(jsonStr), &res.PriceVo)
	if jsonErr != nil {
		return res, common.NewServiceError(err.GetMsg())
	}

	return
}

func (s gameService) GetUpdateNews() (res models.UpdateNewsVo, err common.GFError) {
	jsonStr, err := cs.GetString("game-news:latest")
	if err != nil {
		return res, err
	}
	jsonErr := sonic.Unmarshal([]byte(jsonStr), &res)
	if jsonErr != nil {
		return res, common.NewServiceError(err.GetMsg())
	}

	return
}
