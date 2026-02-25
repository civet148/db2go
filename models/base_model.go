package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" gorm:"column:created_at;type:timestamp;not null;index;default:CURRENT_TIMESTAMP;autoCreatedAt"` //创建时间
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" gorm:"column:updated_at;type:timestamp;not null;index;default:CURRENT_TIMESTAMP;autoUpdatedAt"` //更新时间
	isExist   bool      `gorm:"-" db:"-"`                                                                                                                            //数据是否在数据库中存在
}
