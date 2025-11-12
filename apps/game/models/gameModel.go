package models

import (
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
