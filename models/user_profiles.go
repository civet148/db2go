package models

import "time"
import github_com_civet148_db2go_types "github.com/civet148/db2go/types"

const TableNameUserProfiles = "user_profiles" //

const (
	USER_PROFILES_COLUMN_ID         = "id"
	USER_PROFILES_COLUMN_CREATED_AT = "created_at"
	USER_PROFILES_COLUMN_UPDATED_AT = "updated_at"
	USER_PROFILES_COLUMN_USER_ID    = "user_id"
	USER_PROFILES_COLUMN_AVATAR     = "avatar"
	USER_PROFILES_COLUMN_ADDRESS    = "address"
)

type UserProfile struct {
	github_com_civet148_db2go_types.BaseModel
	Id        uint64    `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                                                            //
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" gorm:"column:created_at;type:timestamp;autoCreateTime;index:idx_user_profiles_created_at;default:CURRENT_TIMESTAMP;" sqlca:"readonly"` //
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" gorm:"column:updated_at;type:timestamp;autoUpdateTime;index:idx_user_profiles_updated_at;default:CURRENT_TIMESTAMP;" sqlca:"readonly"` //
	UserId    uint64    `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;type:bigint unsigned;uniqueIndex:idx_user_profiles_user_id;" sqlca:"isnull"`                                            //
	Avatar    string    `json:"avatar,omitempty" db:"avatar" gorm:"column:avatar;type:varchar(512);" sqlca:"isnull"`                                                                                        //
	Address   string    `json:"address,omitempty" db:"address" gorm:"column:address;type:varchar(128);" sqlca:"isnull"`                                                                                     //
}

func (do UserProfile) TableName() string { return "user_profiles" }

func (do UserProfile) GetId() uint64           { return do.Id }
func (do UserProfile) GetCreatedAt() time.Time { return do.CreatedAt }
func (do UserProfile) GetUpdatedAt() time.Time { return do.UpdatedAt }
func (do UserProfile) GetUserId() uint64       { return do.UserId }
func (do UserProfile) GetAvatar() string       { return do.Avatar }
func (do UserProfile) GetAddress() string      { return do.Address }

func (do *UserProfile) SetId(v uint64)           { do.Id = v }
func (do *UserProfile) SetCreatedAt(v time.Time) { do.CreatedAt = v }
func (do *UserProfile) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }
func (do *UserProfile) SetUserId(v uint64)       { do.UserId = v }
func (do *UserProfile) SetAvatar(v string)       { do.Avatar = v }
func (do *UserProfile) SetAddress(v string)      { do.Address = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
