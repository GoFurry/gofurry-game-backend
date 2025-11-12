package dao

import (
	"errors"
	"github.com/GoFurry/gofurry-game-backend/apps/system/user/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/abstract"
	"gorm.io/gorm"
)

var newUserDao = new(userDao)

func init() {
	newUserDao.Init()
	newUserDao.Mode = models.GfUser{}
}

type userDao struct{ abstract.Dao }

func GetUserDao() *userDao { return newUserDao }

func (dao *userDao) FindOneByName(name string) (record models.GfUser, err common.GFError) {
	db := dao.Gm.Table("gf_user").Where("name = ?", name).Take(&record)
	if err := db.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, common.NewDaoError(common.RETURN_RECORD_NOT_FOUND)
		} else {
			return record, common.NewDaoError(err.Error())
		}
	}
	return
}
