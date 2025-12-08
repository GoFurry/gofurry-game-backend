package dao

import (
	gm "github.com/GoFurry/gofurry-game-backend/apps/game/models"
	"github.com/GoFurry/gofurry-game-backend/apps/review/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/abstract"
)

var newReviewDao = new(reviewDao)

func init() {
	newReviewDao.Init()
}

type reviewDao struct{ abstract.Dao }

func GetReviewDao() *reviewDao { return newReviewDao }

func (dao reviewDao) GetHotGame(num int) (res []models.AvgScoreResult, err common.GFError) {
	var results []models.AvgScoreResult
	commentTable := models.TableNameGfgGameComment
	gameTable := gm.TableNameGfgGame

	db := dao.Gm.Table(commentTable).
		Joins("LEFT JOIN "+gameTable+" ON "+commentTable+".game_id = "+gameTable+".id").
		Select(
			commentTable+".game_id",
			"AVG("+commentTable+".score) AS avg_score",
			"COUNT(*) AS comment_count",
			gameTable+".name",
			gameTable+".name_en",
			gameTable+".info",
			gameTable+".info_en",
			gameTable+".header",
		).
		Group(commentTable + ".game_id, " + gameTable + ".name, " + gameTable + ".name_en, " +
			gameTable + ".info, " + gameTable + ".info_en, " + gameTable + ".header").
		Order("avg_score DESC").
		Limit(num)

	if dbErr := db.Find(&results).Error; dbErr != nil {
		return res, common.NewDaoError(dbErr.Error())
	}

	return results, nil
}

func (dao reviewDao) GetScoreById(id int64) (res models.AvgScoreResult, err common.GFError) {
	commentTable := models.TableNameGfgGameComment
	gameTable := gm.TableNameGfgGame

	db := dao.Gm.Table(commentTable).
		Joins("LEFT JOIN "+gameTable+" ON "+commentTable+".game_id = "+gameTable+".id").
		Select(
			commentTable+".game_id",
			"AVG("+commentTable+".score) AS avg_score",
			"COUNT(*) AS comment_count",
			gameTable+".name",
			gameTable+".name_en",
			gameTable+".info",
			gameTable+".info_en",
			gameTable+".header",
		).
		Where(commentTable+".game_id = ?", id).
		Group(commentTable + ".game_id, " + gameTable + ".name, " + gameTable + ".name_en, " + gameTable +
			".info, " + gameTable + ".info_en, " + gameTable + ".header")

	if dbErr := db.Take(&res).Error; dbErr != nil {
		return res, common.NewDaoError(dbErr.Error())
	}

	return res, nil
}
