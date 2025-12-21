package models

import (
	rm "github.com/GoFurry/gofurry-game-backend/apps/review/models"
	cm "github.com/GoFurry/gofurry-game-backend/common/models"
)

const TableNameGfgGame = "gfg_game"

// GfgGame mapped from table <gfg_game>
type GfgGame struct {
	ID          int64        `gorm:"column:id;type:bigint;primaryKey;comment:游戏表ID" json:"id"`                                         // 游戏表ID
	Name        string       `gorm:"column:name;type:character varying(255);not null;comment:游戏名称" json:"name"`                        // 游戏名称
	NameEn      string       `gorm:"column:name_en;type:character varying(255);not null;comment:游戏英文名称" json:"nameEn"`                 // 游戏英文名称
	Info        string       `gorm:"column:info;type:character varying(300);not null;comment:游戏简介" json:"info"`                        // 游戏简介
	InfoEn      string       `gorm:"column:info_en;type:character varying(300);not null;comment:游戏英文简介" json:"infoEn"`                 // 游戏英文简介
	CreateTime  cm.LocalTime `gorm:"column:create_time;type:int;type:unsigned;not null;autoCreateTime;comment:创建时间" json:"createTime"` // 创建时间
	UpdateTime  cm.LocalTime `gorm:"column:update_time;type:int;type:unsigned;not null;autoUpdateTime;comment:更新时间" json:"updateTime"` // 更新时间
	Resources   *string      `gorm:"column:resources;type:json;comment:游戏相关资源" json:"resources"`                                       // 游戏相关资源
	Groups      *string      `gorm:"column:groups;type:json;comment:游戏相关社群" json:"groups"`                                             // 游戏相关社群
	ReleaseDate string       `gorm:"column:release_date;type:character varying(255);not null;comment:发行日期" json:"releaseDate"`         // 发行日期
	Developers  string       `gorm:"column:developers;type:json;not null;comment:开发商" json:"developers"`                               // 开发商
	Publishers  string       `gorm:"column:publishers;type:json;not null;comment:发行商" json:"publishers"`                               // 发行商
	Appid       int64        `gorm:"column:appid;type:bigint;not null;comment:SteamAPI appid" json:"appid"`                            // SteamAPI appid
	Header      string       `gorm:"column:header;type:character varying(255);not null;comment:游戏封面图" json:"header"`                   // 游戏封面图
	Links       *string      `gorm:"column:links;type:json;comment:三方网站链接" json:"links"`                                               // 三方网站链接
	Weight      int64        `gorm:"column:weight;type:bigint;not null;comment:权重" json:"weight"`                                      // 权重
}

// TableName GfgGame's table name
func (*GfgGame) TableName() string {
	return TableNameGfgGame
}

type GameRespVo struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Info        string        `json:"info"`
	CreateTime  cm.LocalTime  `json:"create_time"`
	UpdateTime  cm.LocalTime  `json:"update_time"`
	Resources   *[]cm.KvModel `json:"resources"`
	Groups      *[]cm.KvModel `json:"groups"`
	ReleaseDate string        `json:"release_date"`
	Developers  []string      `json:"developers"`
	Publishers  []string      `json:"publishers"`
	Appid       string        `json:"appid"`
	Header      string        `json:"header"`
	Links       *[]cm.KvModel `json:"links"`
}

const TableNameGfgGameRecord = "gfg_game_record"

// GfgGameRecord mapped from table <gfg_game_record>
type GfgGameRecord struct {
	ID          int64  `gorm:"column:id;type:bigint;primaryKey;comment:游戏记录表id" json:"id"`                              // 游戏记录表id
	GameID      int64  `gorm:"column:game_id;type:bigint;not null;comment:游戏表id" json:"gameId,string"`                  // 游戏表id
	Language    string `gorm:"column:language;type:text;not null;comment:支持语言" json:"language"`                         // 支持语言
	ReleaseDate string `gorm:"column:release_date;type:character varying(30);not null;comment:发行时间" json:"releaseDate"` // 发行时间
	Platform    string `gorm:"column:platform;type:character varying(50);not null;comment:支持平台" json:"platform"`        // 支持平台
	Developer   string `gorm:"column:developer;type:character varying(100);not null;comment:开发商" json:"developer"`      // 开发商
	Publisher   string `gorm:"column:publisher;type:character varying(100);not null;comment:发行商" json:"publisher"`      // 发行商
	Info        string `gorm:"column:info;type:text;not null;comment:游戏概述" json:"info"`                                 // 游戏概述
	Cover       string `gorm:"column:cover;type:character varying(255);comment:封面图" json:"cover"`                       // 封面图
	HotIndex    int64  `gorm:"column:hot_index;type:bigint;not null;comment:热度指数" json:"hotIndex"`                      // 热度指数
	Lang        string `gorm:"column:lang;type:character varying(20);not null;comment:记录的语言" json:"lang"`               // 记录的语言
	PriceList   string `gorm:"column:price_list;type:json;not null;comment:游戏价格列表" json:"priceList"`                    // 游戏价格列表
	Initial     int64  `gorm:"column:initial;type:bigint;not null;comment:游戏价格" json:"initial"`                         // 游戏价格
	Final       int64  `gorm:"column:final;type:bigint;not null;comment:当前价格" json:"final"`                             // 当前价格
	Discount    int64  `gorm:"column:discount;type:bigint;not null;comment:折扣百分比" json:"discount"`                      // 折扣百分比
}

// TableName GfgGameRecord's table name
func (*GfgGameRecord) TableName() string {
	return TableNameGfgGameRecord
}

type GameMainInfoVo struct {
	Latest []rm.AvgScoreResult `json:"latest"`
	Recent []rm.AvgScoreResult `json:"recent"`
	Hot    []rm.AvgScoreResult `json:"hot"`
	Free   []rm.AvgScoreResult `json:"free"`
}

const TableNameGfgGamePlayerCount = "gfg_game_player_count"

// GfgGamePlayerCount mapped from table <gfg_game_player_count>
type GfgGamePlayerCount struct {
	ID         int64        `gorm:"column:id;type:bigint;primaryKey;comment:在线人数表ID" json:"id"`                                       // 在线人数表ID
	GameID     int64        `gorm:"column:game_id;type:bigint;not null;comment:游戏表ID" json:"gameId,string"`                           // 游戏表ID
	Count_     int64        `gorm:"column:count;type:bigint;not null;comment:在线人数" json:"count"`                                      // 在线人数
	CreateTime cm.LocalTime `gorm:"column:create_time;type:int;type:unsigned;not null;autoCreateTime;comment:创建时间" json:"createTime"` // 创建时间
}

// TableName GfgGamePlayerCount's table name
func (*GfgGamePlayerCount) TableName() string {
	return TableNameGfgGamePlayerCount
}

type PlayerTopCountVo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CountPeak   int64  `json:"count_peak"`
	CountRecent int64  `json:"count_recent"`
	CollectTime int64  `json:"collect_time"`
	Header      string `json:"header"`
}

type TopPriceVo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	GlobalPrice int64  `json:"global_price"`
	ChinaPrice  int64  `json:"china_price"`
	Discount    int64  `json:"discount"`
	Header      string `json:"header"`
}

type GameMainPanelVo struct {
	CountVo []PlayerTopCountVo `json:"count_vo"`
	PriceVo []TopPriceVo       `json:"price_vo"`
}

type UpdateNewsModels struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	PostTime cm.LocalTime `json:"post_time"`
	Headline string       `json:"headline"`
	Header   string       `json:"header"`
	Author   string       `json:"author"`
	Content  string       `json:"content"`
	Url      string       `json:"url"`
}

type UpdateNewsVo struct {
	NewsZh []UpdateNewsModels `json:"news_zh"`
	NewsEn []UpdateNewsModels `json:"news_en"`
}

const TableNameGfgGameNews = "gfg_game_news"

// GfgGameNews mapped from table <gfg_game_news>
type GfgGameNews struct {
	ID         int64        `gorm:"column:id;type:bigint;primaryKey;comment:游戏更新公告记录表id" json:"id"`                                   // 游戏更新公告记录表id
	GameID     int64        `gorm:"column:game_id;type:bigint;not null;comment:游戏表id" json:"gameId,string"`                           // 游戏表id
	Headline   string       `gorm:"column:headline;type:character varying(255);not null;comment:更新公告标题" json:"headline"`              // 更新公告标题
	Content    string       `gorm:"column:content;type:text;not null;comment:更新公告内容" json:"content"`                                  // 更新公告内容
	Index      int64        `gorm:"column:index;type:bigint;not null;comment:更新公告编号" json:"index"`                                    // 更新公告编号
	PostTime   cm.LocalTime `gorm:"column:post_time;type:timestamp(0) without time zone;not null;comment:更新公告上传日期" json:"postTime"`   // 更新公告上传日期
	CreateTime cm.LocalTime `gorm:"column:create_time;type:int;type:unsigned;not null;autoCreateTime;comment:采集时间" json:"createTime"` // 采集时间
	Author     string       `gorm:"column:author;type:character varying(50);not null;comment:公告作者" json:"author"`                     // 公告作者
	URL        string       `gorm:"column:url;type:character varying(255);not null;comment:更新公告原始地址" json:"url"`                      // 更新公告原始地址
	Total      int64        `gorm:"column:total;type:bigint;not null;comment:公告总数" json:"total"`                                      // 公告总数
	Lang       string       `gorm:"column:lang;type:character varying(30);not null;comment:记录的语言" json:"lang"`                        // 记录的语言
}

// TableName GfgGameNews's table name
func (*GfgGameNews) TableName() string {
	return TableNameGfgGameNews
}
