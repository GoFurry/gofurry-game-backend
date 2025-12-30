package dao

import (
	"github.com/GoFurry/gofurry-game-backend/apps/game/models"
	gm "github.com/GoFurry/gofurry-game-backend/apps/recommend/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/abstract"
	"gorm.io/gorm"
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

func (dao gameDao) GetLatestGame(num int) (res []int64, err common.GFError) {
	db := dao.Gm.Table(models.TableNameGfgGame).Select("id").Order("release_date DESC").Limit(num)
	db.Find(&res)
	if dbErr := db.Error; dbErr != nil {
		return res, common.NewDaoError(dbErr.Error())
	}
	return res, nil
}

func (dao gameDao) GetRecentGame(num int) (res []int64, err common.GFError) {
	db := dao.Gm.Table(models.TableNameGfgGame).Select("id").Order("create_time DESC").Limit(num)
	db.Find(&res)
	if dbErr := db.Error; dbErr != nil {
		return res, common.NewDaoError(dbErr.Error())
	}
	return res, nil
}

func (dao gameDao) GetFreeGame(num int) (res []int64, err common.GFError) {
	db := dao.Gm.Table(models.TableNameGfgGameRecord).Select("game_id")
	db.Where("lang=? AND initial=? AND final=?", "en", 0, 0).Limit(num)
	db.Find(&res)
	if dbErr := db.Error; dbErr != nil {
		return res, common.NewDaoError(dbErr.Error())
	}
	return res, nil
}

func (dao gameDao) GetPlayerPeak(num int) (res []models.PlayerTopCountVo, err common.GFError) {
	countTable := models.TableNameGfgGamePlayerCount
	gameTable := models.TableNameGfgGame

	// 获取每个game_id的最新记录时间
	latestTimeSubQuery := dao.Gm.Table(countTable).
		Select("game_id, MAX(create_time) AS latest_create_time").
		Group("game_id")

	// 通过最新时间关联, 获取每个game_id的最新count和collect_time
	latestCountSubQuery := dao.Gm.Table(countTable).
		Joins("JOIN (?) AS latest_time ON "+countTable+".game_id = latest_time.game_id AND "+countTable+".create_time = latest_time.latest_create_time", latestTimeSubQuery).
		Select(countTable + ".game_id, " + countTable + ".count AS count_recent, " + countTable + ".create_time AS collect_time")

	db := dao.Gm.Table(countTable).
		// 关联最新count子查询
		Joins("JOIN (?) AS latest_count ON "+countTable+".game_id = latest_count.game_id", latestCountSubQuery).
		// 关联游戏表获取名称和封面图
		Joins("LEFT JOIN " + gameTable + " ON " + countTable + ".game_id = " + gameTable + ".id").
		Select(`
			CAST(` + countTable + `.game_id AS VARCHAR) AS id,
			COALESCE(` + gameTable + `.name, '未知游戏') AS name,
			MAX(` + countTable + `.count) AS count_peak,
			latest_count.count_recent AS count_recent,
			EXTRACT(EPOCH FROM latest_count.collect_time)::INT8 AS collect_time,
			COALESCE(` + gameTable + `.header, '') AS header
		`).
		// 按game_id分组
		Group(countTable + ".game_id, " + gameTable + ".name, " + gameTable + ".header, latest_count.count_recent, latest_count.collect_time").
		// 按峰值降序排序
		Order("count_peak DESC").
		// 限制返回条数
		Limit(num)

	// 执行查询并处理错误
	if dbErr := db.Find(&res).Error; dbErr != nil {
		return res, common.NewDaoError("获取游戏在线人数峰值失败: " + dbErr.Error())
	}
	if len(res) == 0 {
		return res, common.NewDaoError("暂无游戏在线人数数据")
	}

	return res, nil
}

func (dao gameDao) GetTopPrice(num int) (res []models.TopPriceVo, err common.GFError) {
	recordTable := models.TableNameGfgGameRecord
	gameTable := models.TableNameGfgGame

	// 获取en区记录
	enSubQuery := dao.Gm.Table(recordTable).
		Select("game_id, final AS global_price, discount").
		Where("lang = ?", "en")

	// 获取zh区记录
	zhSubQuery := dao.Gm.Table(recordTable).
		Select("game_id, final AS china_price").
		Where("lang = ?", "zh")

	// 关联en区、zh区、游戏表
	db := dao.Gm.Table("(?) AS en_data", enSubQuery).
		// 左关联zh区数据
		Joins("LEFT JOIN (?) AS zh_data ON en_data.game_id = zh_data.game_id", zhSubQuery).
		// 左关联游戏表获取名称和封面图
		Joins("LEFT JOIN " + gameTable + " ON en_data.game_id = " + gameTable + ".id").
		Select(`
			CAST(en_data.game_id AS VARCHAR) AS id,
			COALESCE(` + gameTable + `.name, '未知游戏') AS name,
			en_data.global_price,
			COALESCE(zh_data.china_price, 0) AS china_price,
			en_data.discount,
			COALESCE(` + gameTable + `.header, '') AS header 
		`).
		// 按全球区价格降序排序
		Order("en_data.global_price DESC").
		// 限制返回条数
		Limit(num)

	// 执行查询
	if dbErr := db.Find(&res).Error; dbErr != nil {
		return res, common.NewDaoError("获取高价游戏失败: " + dbErr.Error())
	}
	if len(res) == 0 {
		return res, common.NewDaoError("暂无游戏价格数据")
	}

	return res, nil
}

func (dao gameDao) GetUpdateNews(num int, lang string) (res []models.UpdateNewsModels, err common.GFError) {
	newsTable := models.TableNameGfgGameNews
	gameTable := models.TableNameGfgGame

	db := dao.Gm.Table(newsTable).
		// 左关联游戏表
		Joins("LEFT JOIN "+gameTable+" ON "+newsTable+".game_id = "+gameTable+".id").
		Select(`
			CAST(`+newsTable+`.game_id AS VARCHAR) AS id,
			COALESCE(`+gameTable+`.name, '未知游戏') AS name,
			`+newsTable+`.post_time,
			`+newsTable+`.headline,
			COALESCE(`+gameTable+`.header, '') AS header,
			`+newsTable+`.author,
			`+newsTable+`.url,
			SUBSTRING(`+newsTable+`.content, 1, 200) AS content
		`).
		Where("lang = ?", lang).
		Order(newsTable + ".post_time DESC").
		Limit(num)

	// 执行查询并处理错误
	if dbErr := db.Find(&res).Error; dbErr != nil {
		return res, common.NewDaoError("获取最新游戏新闻失败: " + dbErr.Error())
	}
	if len(res) == 0 {
		return res, common.NewDaoError("暂无游戏新闻数据")
	}

	return res, nil
}

func (dao gameDao) GetTagList(lang string) (res []models.TagModelVo, err common.GFError) {
	var countSubQuery *gorm.DB
	countSubQuery = dao.Gm.Table(gm.TableNameGfgTagMap).
		Select("tag_id, COUNT(*) as game_count").
		Group("tag_id")

	nameField := "gfg_tag.name AS name"
	if lang == "en" {
		nameField = "gfg_tag.name_en AS name"
	}

	db := dao.Gm.Table(gm.TableNameGfgTag).
		Joins("LEFT JOIN (?) AS tag_count ON gfg_tag.id = tag_count.tag_id", countSubQuery).
		Select(
			"CAST(gfg_tag.id AS VARCHAR)",
			nameField,
			"CAST(gfg_tag.prefix AS VARCHAR) AS prefix",
			"COALESCE(tag_count.game_count, 0) AS game_count",
		).Order("game_count DESC")

	if dbErr := db.Find(&res).Error; dbErr != nil {
		return res, common.NewDaoError("获取标签记录失败: " + dbErr.Error())
	}

	return res, nil
}
