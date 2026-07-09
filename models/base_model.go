package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" gorm:"column:created_at;type:timestamp;not null;index;default:CURRENT_TIMESTAMP;autoCreatedAt"` //创建时间
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" gorm:"column:updated_at;type:timestamp;not null;index;default:CURRENT_TIMESTAMP;autoUpdatedAt"` //更新时间
	isExist   bool      `gorm:"-" db:"-"`                                                                                                                            //数据是否在数据库中存在
}

func (do BaseModel) GetCreatedAt() time.Time { return do.CreatedAt }

func (do BaseModel) GetUpdatedAt() time.Time { return do.UpdatedAt }

func (do *BaseModel) SetCreatedAt(v time.Time) { do.CreatedAt = v }

func (do *BaseModel) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }
