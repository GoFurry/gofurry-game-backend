package dao

import (
	"github.com/GoFurry/gofurry-game-backend/apps/game/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/abstract"
)

var newGameDao = new(gameDao)

func init() {
	newGameDao.Init()
}

type gameDao struct{ abstract.Dao }

func GetGameDao() *gameDao { return newGameDao }

func (dao gameDao) GetGameList(num int) (res []models.GfgGame, err common.GFError) {
	db := dao.Gm.Table(models.TableNameGfgGame)
	db.Order("weight ASC").Limit(num)
	db.Find(&res)
	if dbErr := db.Error; dbErr != nil {
		return res, common.NewDaoError(dbErr.Error())
	}
	return
}

func (dao gameDao) GetByNum(randomInt int) (res models.GfgGame, err common.GFError) {
	db := dao.Gm.Table(models.TableNameGfgGame).Order("id DESC")
	db.Offset(randomInt).Limit(1)
	db.Take(&res)
	if dbErr := db.Error; dbErr != nil {
		return res, common.NewDaoError(dbErr.Error())
	}
	return res, nil
}
