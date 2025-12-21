package dao

import (
	"strings"

	gm "github.com/GoFurry/gofurry-game-backend/apps/game/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/abstract"
)

var newSearchDao = new(searchDao)

func init() {
	newSearchDao.Init()
}

type searchDao struct{ abstract.Dao }

func GetSearchDao() *searchDao { return newSearchDao }

func (dao searchDao) GetGameListByText(text string, limit int) (res []gm.GfgGame, err common.GFError) {
	db := dao.Gm.Table(gm.TableNameGfgGame)
	if db.Error != nil {
		return res, common.NewDaoError(db.Error.Error())
	}

	// 不区分大小写模糊匹配
	searchText := "%" + strings.TrimSpace(text) + "%"
	db.Where("name ILIKE ? OR name_en ILIKE ? OR info ILIKE ? OR info_en ILIKE ?",
		searchText, searchText, searchText, searchText)
	db.Order("weight ASC, update_time DESC").Limit(limit)

	if errDb := db.Find(&res).Error; errDb != nil {
		return res, common.NewDaoError(errDb.Error())
	}
	return res, nil
}
