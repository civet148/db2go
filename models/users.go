package models

import "time"
import github_com_civet148_db2go_types "github.com/civet148/db2go/types"

const TableNameUsers = "users" //

const (
	USERS_COLUMN_ID          = "id"
	USERS_COLUMN_CREATE_TIME = "create_time"
	USERS_COLUMN_CREATE_ID   = "create_id"
	USERS_COLUMN_CREATE_NAME = "create_name"
	USERS_COLUMN_UPDATE_ID   = "update_id"
	USERS_COLUMN_UPDATE_NAME = "update_name"
	USERS_COLUMN_UPDATE_TIME = "update_time"
	USERS_COLUMN_USER_NAME   = "user_name"
	USERS_COLUMN_EMAIL       = "email"
)

type User struct {
	github_com_civet148_db2go_types.BaseModel
	Id         uint64 `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                             //
	CreateId   uint64 `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id;type:bigint unsigned;" sqlca:"isnull"`                             //
	CreateName string `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name;type:longtext;" sqlca:"isnull"`                              //
	UpdateId   uint64 `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id;type:bigint unsigned;" sqlca:"isnull"`                             //
	UpdateName string `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name;type:longtext;" sqlca:"isnull"`                              //
	UserName   string `json:"user_name,omitempty" db:"user_name" gorm:"column:user_name;type:varchar(32);uniqueIndex:idx_users_user_name;" sqlca:"isnull"` //
	Email      string `json:"email,omitempty" db:"email" gorm:"column:email;type:varchar(64);uniqueIndex:idx_users_email;" sqlca:"isnull"`                 //
}

func (do User) TableName() string { return "users" }

func (do User) GetId() uint64            { return do.Id }
func (do User) GetCreateTime() time.Time { return do.CreateTime }
func (do User) GetCreateId() uint64      { return do.CreateId }
func (do User) GetCreateName() string    { return do.CreateName }
func (do User) GetUpdateId() uint64      { return do.UpdateId }
func (do User) GetUpdateName() string    { return do.UpdateName }
func (do User) GetUpdateTime() time.Time { return do.UpdateTime }
func (do User) GetUserName() string      { return do.UserName }
func (do User) GetEmail() string         { return do.Email }

func (do *User) SetId(v uint64)            { do.Id = v }
func (do *User) SetCreateTime(v time.Time) { do.CreateTime = v }
func (do *User) SetCreateId(v uint64)      { do.CreateId = v }
func (do *User) SetCreateName(v string)    { do.CreateName = v }
func (do *User) SetUpdateId(v uint64)      { do.UpdateId = v }
func (do *User) SetUpdateName(v string)    { do.UpdateName = v }
func (do *User) SetUpdateTime(v time.Time) { do.UpdateTime = v }
func (do *User) SetUserName(v string)      { do.UserName = v }
func (do *User) SetEmail(v string)         { do.Email = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
