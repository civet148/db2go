package types

import "time"

type BaseModel struct {
	CreateTime time.Time `json:"create_time,omitempty" db:"create_time" gorm:"column:create_time;default:CURRENT_TIMESTAMP;autoCreateTime"` //创建时间
	UpdateTime time.Time `json:"update_time,omitempty" db:"update_time" gorm:"column:update_time;default:CURRENT_TIMESTAMP;autoUpdateTime"` //更新时间
}
