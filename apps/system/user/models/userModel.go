package models

import (
	"github.com/GoFurry/gofurry-game-backend/common/abstract"
	cm "github.com/GoFurry/gofurry-game-backend/common/models"
)

const TableNameGfUser = "gf_user"

// GfUser mapped from table <gf_user>
type GfUser struct {
	abstract.DefaultModel
	Password   string       `gorm:"column:password;type:character varying(255);not null;comment:密码" json:"password"`                  // 密码
	CreateTime cm.LocalTime `gorm:"column:create_time;type:int;type:unsigned;not null;autoCreateTime;comment:创建时间" json:"createTime"` // 创建时间
}

// TableName GfUser's table name
func (*GfUser) TableName() string {
	return TableNameGfUser
}

type CurrentUser struct {
	Name string `json:"name"`
	ID   int64  `json:"id,string"`
}

type UserLoginRequest struct {
	Name     string `form:"name" json:"name" validate:"required,min=3,max=36"`
	Password string `form:"password" json:"password" validate:"required"`
}
