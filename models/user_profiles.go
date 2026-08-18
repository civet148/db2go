package models

import (
	"time"
)

const TableNameUserProfiles = "user_profiles"

const (
	UserProfilesColumn_Id        = "id"
	UserProfilesColumn_CreatedAt = "created_at"
	UserProfilesColumn_UpdatedAt = "updated_at"
	UserProfilesColumn_UserId    = "user_id"
	UserProfilesColumn_Avatar    = "avatar"
	UserProfilesColumn_Address   = "address"
)

type UserProfile struct {
	Id        uint64    `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"column:created_at;type:timestamp;autoCreateTime;index:idx_user_profiles_created_at,priority:1;default:CURRENT_TIMESTAMP;" sqlca:"readonly"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"column:updated_at;type:timestamp;autoUpdateTime;index:idx_user_profiles_updated_at,priority:1;default:CURRENT_TIMESTAMP;" sqlca:"readonly"`
	UserId    uint64    `json:"user_id" db:"user_id" gorm:"column:user_id;type:bigint unsigned;uniqueIndex:idx_user_profiles_user_id,priority:1;default:null;" sqlca:"isnull"`
	Avatar    string    `json:"avatar" db:"avatar" gorm:"column:avatar;type:varchar(512);default:null;" sqlca:"isnull"`
	Address   string    `json:"address" db:"address" gorm:"column:address;type:varchar(128);default:null;" sqlca:"isnull"`
}

func (do UserProfile) TableName() string {
	return TableNameUserProfiles
}

func (do UserProfile) GetId() uint64 { return do.Id }

func (do UserProfile) GetCreatedAt() time.Time { return do.CreatedAt }

func (do UserProfile) GetUpdatedAt() time.Time { return do.UpdatedAt }

func (do UserProfile) GetUserId() uint64 { return do.UserId }

func (do UserProfile) GetAvatar() string { return do.Avatar }

func (do UserProfile) GetAddress() string { return do.Address }

func (do *UserProfile) SetId(v uint64) { do.Id = v }

func (do *UserProfile) SetCreatedAt(v time.Time) { do.CreatedAt = v }

func (do *UserProfile) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }

func (do *UserProfile) SetUserId(v uint64) { do.UserId = v }

func (do *UserProfile) SetAvatar(v string) { do.Avatar = v }

func (do *UserProfile) SetAddress(v string) { do.Address = v }
