package models

import "time"

const TableNameUsers = "users" //

const (
	USERS_COLUMN_ID         = "id"
	USERS_COLUMN_CREATED_AT = "created_at"
	USERS_COLUMN_UPDATED_AT = "updated_at"
	USERS_COLUMN_USER_NAME  = "user_name"
	USERS_COLUMN_EMAIL      = "email"
)

type User struct {
	BaseModel
	Id       uint64 `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                             //
	UserName string `json:"user_name,omitempty" db:"user_name" gorm:"column:user_name;type:varchar(32);uniqueIndex:idx_users_user_name;" sqlca:"isnull"` //
	Email    string `json:"email,omitempty" db:"email" gorm:"column:email;type:varchar(64);uniqueIndex:idx_users_email;" sqlca:"isnull"`                 //
}

func (do User) TableName() string { return "users" }

func (do User) GetId() uint64           { return do.Id }
func (do User) GetCreatedAt() time.Time { return do.CreatedAt }
func (do User) GetUpdatedAt() time.Time { return do.UpdatedAt }
func (do User) GetUserName() string     { return do.UserName }
func (do User) GetEmail() string        { return do.Email }

func (do *User) SetId(v uint64)           { do.Id = v }
func (do *User) SetCreatedAt(v time.Time) { do.CreatedAt = v }
func (do *User) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }
func (do *User) SetUserName(v string)     { do.UserName = v }
func (do *User) SetEmail(v string)        { do.Email = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
