package dao

import (
	"github.com/GoFurry/gofurry-game-backend/apps/recommend/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/abstract"
)

var newRecommendDao = new(recommendDao)

func init() {
	newRecommendDao.Init()
}

type recommendDao struct{ abstract.Dao }

func GetRecommendDao() *recommendDao { return newRecommendDao }

func (dao recommendDao) GetTagMappingList() (res []models.GfgTagMap, gfError common.GFError) {
	db := dao.Gm.Table(models.TableNameGfgTagMap).Find(&res)
	if err := db.Error; err != nil {
		return res, common.NewDaoError(err.Error())
	}
	return res, nil
}

func (dao recommendDao) GetTagList() (res []models.GfgTag, gfError common.GFError) {
	db := dao.Gm.Table(models.TableNameGfgTag).Find(&res)
	if err := db.Error; err != nil {
		return res, common.NewDaoError(err.Error())
	}
	return res, nil
}
