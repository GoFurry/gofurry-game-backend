package task

import (
	"time"

	gd "github.com/GoFurry/gofurry-game-backend/apps/game/dao"
	gm "github.com/GoFurry/gofurry-game-backend/apps/game/models"
	rd "github.com/GoFurry/gofurry-game-backend/apps/review/dao"
	rm "github.com/GoFurry/gofurry-game-backend/apps/review/models"
	"github.com/GoFurry/gofurry-game-backend/apps/schedule/models"
	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/common/log"
	cs "github.com/GoFurry/gofurry-game-backend/common/service"
	"github.com/GoFurry/gofurry-game-backend/common/util"
	"github.com/bytedance/sonic"
)

func UpdateMainInfoCache() {
	log.Info("StatTask UpdateMainInfoCache 开始...")

	info := map[string]any{
		"latest": LatestInfo{models.InfoModel{Key: "game-info:latest", Num: 8, Duration: 3 * time.Hour}},
		"recent": RecentInfo{models.InfoModel{Key: "game-info:recent", Num: 8, Duration: 3 * time.Hour}},
		"free":   FreeInfo{models.InfoModel{Key: "game-info:free", Num: 8, Duration: 3 * time.Hour}},
		"hot":    HotInfo{models.InfoModel{Key: "game-info:hot", Num: 8, Duration: 3 * time.Hour}},
	}
	for k, v := range info {
		switch k {
		case "latest":
			latestInfo := v.(LatestInfo)
			latestInfo.cacheGameInfo()
		case "recent":
			recentInfo := v.(RecentInfo)
			recentInfo.cacheGameInfo()
		case "free":
			freeInfo := v.(FreeInfo)
			freeInfo.cacheGameInfo()
		case "hot":
			hotInfo := v.(HotInfo)
			hotInfo.cacheGameInfo()
		}
	}
	log.Info("StatTask UpdateMainInfoCache 结束...")
}

type HotInfo struct {
	models.InfoModel
}

func (r *HotInfo) cacheGameInfo() common.GFError {
	res, err := rd.GetReviewDao().GetHotGame(r.Num)
	if err != nil {
		return err
	}
	if idList, jsonErr := sonic.Marshal(res); jsonErr == nil {
		cs.SetExpire(r.Key, string(idList), r.Duration)
	}
	return nil
}

type FreeInfo struct {
	models.InfoModel
}

func (r *FreeInfo) cacheGameInfo() common.GFError {
	res, err := gd.GetGameDao().GetFreeGame(r.Num)
	if err != nil {
		return err
	}
	infoRecord := []rm.AvgScoreResult{}
	for _, v := range res {
		newRecord, gfError := rd.GetReviewDao().GetScoreById(v)
		if gfError != nil && gfError.GetMsg() == "record not found" {
			newRecord = rm.AvgScoreResult{GameID: util.Int642String(v), AvgScore: 0.0, CommentCount: 0}
			gameRecord := gm.GfgGame{}
			gd.GetGameDao().GetById(v, &gameRecord)
			newRecord.Name = gameRecord.Name
			newRecord.NameEn = gameRecord.NameEn
			newRecord.Info = gameRecord.Info
			newRecord.InfoEn = gameRecord.InfoEn
			newRecord.Header = gameRecord.Header
		} else if gfError != nil {
			return gfError
		}
		infoRecord = append(infoRecord, newRecord)
	}
	if idList, jsonErr := sonic.Marshal(infoRecord); jsonErr == nil {
		cs.SetExpire(r.Key, string(idList), r.Duration)
	}
	return nil
}

type RecentInfo struct {
	models.InfoModel
}

func (r *RecentInfo) cacheGameInfo() common.GFError {
	res, err := gd.GetGameDao().GetRecentGame(r.Num)
	if err != nil {
		return err
	}
	infoRecord := []rm.AvgScoreResult{}
	for _, v := range res {
		newRecord, gfError := rd.GetReviewDao().GetScoreById(v)
		if gfError != nil && gfError.GetMsg() == "record not found" {
			newRecord = rm.AvgScoreResult{GameID: util.Int642String(v), AvgScore: 0.0, CommentCount: 0}
			gameRecord := gm.GfgGame{}
			gd.GetGameDao().GetById(v, &gameRecord)
			newRecord.Name = gameRecord.Name
			newRecord.NameEn = gameRecord.NameEn
			newRecord.Info = gameRecord.Info
			newRecord.InfoEn = gameRecord.InfoEn
			newRecord.Header = gameRecord.Header
		} else if gfError != nil {
			return gfError
		}
		infoRecord = append(infoRecord, newRecord)
	}
	if idList, jsonErr := sonic.Marshal(infoRecord); jsonErr == nil {
		cs.SetExpire(r.Key, string(idList), r.Duration)
	}
	return nil
}

type LatestInfo struct {
	models.InfoModel
}

func (l *LatestInfo) cacheGameInfo() common.GFError {
	res, err := gd.GetGameDao().GetLatestGame(l.Num)
	if err != nil {
		return err
	}
	infoRecord := []rm.AvgScoreResult{}
	for _, v := range res {
		newRecord, gfError := rd.GetReviewDao().GetScoreById(v)
		if gfError != nil && gfError.GetMsg() == "record not found" {
			newRecord = rm.AvgScoreResult{GameID: util.Int642String(v), AvgScore: 0.0, CommentCount: 0}
			gameRecord := gm.GfgGame{}
			gd.GetGameDao().GetById(v, &gameRecord)
			newRecord.Name = gameRecord.Name
			newRecord.NameEn = gameRecord.NameEn
			newRecord.Info = gameRecord.Info
			newRecord.InfoEn = gameRecord.InfoEn
			newRecord.Header = gameRecord.Header
		} else if gfError != nil {
			return gfError
		}
		infoRecord = append(infoRecord, newRecord)
	}
	if idList, jsonErr := sonic.Marshal(infoRecord); jsonErr == nil {
		cs.SetExpire(l.Key, string(idList), l.Duration)
	}
	return nil
}

func UpdateGamePanelCache() {
	log.Info("StatTask UpdateGamePanelCache 开始...")
	// 在线人数
	record, err := gd.GetGameDao().GetPlayerPeak(15)
	if err != nil {
		log.Error("GetPlayerPeak err:", err)
		return
	}
	if jsonRecord, jsonErr := sonic.Marshal(record); jsonErr == nil {
		cs.SetExpire("game-panel:top-player-count", string(jsonRecord), 3*time.Hour)
	}
	// 售价
	priceRecord, err := gd.GetGameDao().GetTopPrice(15)
	if err != nil {
		log.Error("GetTopPrice err:", err)
		return
	}
	if jsonRecord, jsonErr := sonic.Marshal(priceRecord); jsonErr == nil {
		cs.SetExpire("game-panel:top-price", string(jsonRecord), 3*time.Hour)
	}

	log.Info("StatTask UpdateGamePanelCache 结束...")
}

func UpdateGameNewsCache() {
	log.Info("StatTask UpdateGameNewsCache 开始...")

	newRecord := gm.UpdateNewsVo{}

	// 最新新闻
	record, err := gd.GetGameDao().GetUpdateNews(15, "zh")
	if err != nil {
		log.Error("GetUpdateNews err:", err)
		return
	}
	newRecord.NewsZh = record

	record, err = gd.GetGameDao().GetUpdateNews(15, "en")
	if err != nil {
		log.Error("GetUpdateNews err:", err)
		return
	}
	newRecord.NewsEn = record

	if jsonRecord, jsonErr := sonic.Marshal(newRecord); jsonErr == nil {
		cs.SetExpire("game-news:latest", string(jsonRecord), 3*time.Hour)
	}

	log.Info("StatTask UpdateGameNewsCache 结束...")
}
