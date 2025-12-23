package types

type BaseModel struct {
	CreateTime string `json:"create_time,omitempty" db:"create_time" gorm:"column:create_time;default:CURRENT_TIMESTAMP;autoCreateTime"` //创建时间
	UpdateTime string `json:"update_time,omitempty" db:"update_time" gorm:"column:update_time;default:CURRENT_TIMESTAMP;autoUpdateTime"` //更新时间
	IsDeleted  int8   `json:"is_deleted,omitempty" db:"is_deleted" gorm:"column:is_deleted"`                                             //删除状态(0: 未删除 1: 已删除)
}
