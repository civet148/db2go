package models

import (
	"time"
)

const TableNameRoles = "roles"

const (
	RolesColumn_Id        = "id"
	RolesColumn_CreatedAt = "created_at"
	RolesColumn_UpdatedAt = "updated_at"
	RolesColumn_Name      = "name"
)

type Role struct {
	Id        uint64     `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`
	CreatedAt *time.Time `json:"created_at" db:"created_at" gorm:"column:created_at;type:timestamp;autoCreateTime;index:idx_roles_created_at,priority:1;default:CURRENT_TIMESTAMP;" sqlca:"readonly"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at" gorm:"column:updated_at;type:timestamp;autoUpdateTime;index:idx_roles_updated_at,priority:1;default:CURRENT_TIMESTAMP;" sqlca:"readonly"`
	Name      string     `json:"name" db:"name" gorm:"column:name;type:varchar(64);uniqueIndex:idx_roles_name,priority:1;default:null;" sqlca:"nullable"`
}

func (do Role) DatabaseName() string {
	return "test"
}

func (do Role) TableName() string {
	return TableNameRoles
}

func (do Role) GetId() uint64 { return do.Id }

func (do Role) GetCreatedAt() *time.Time { return do.CreatedAt }

func (do Role) GetUpdatedAt() *time.Time { return do.UpdatedAt }

func (do Role) GetName() string { return do.Name }

func (do *Role) SetId(v uint64) { do.Id = v }

func (do *Role) SetCreatedAt(v *time.Time) { do.CreatedAt = v }

func (do *Role) SetUpdatedAt(v *time.Time) { do.UpdatedAt = v }

func (do *Role) SetName(v string) { do.Name = v }
