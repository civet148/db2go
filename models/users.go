package models

import (
	"time"
)

const TableNameUsers = "users"

const (
	UsersColumn_Id        = "id"
	UsersColumn_CreatedAt = "created_at"
	UsersColumn_UpdatedAt = "updated_at"
	UsersColumn_UserName  = "user_name"
	UsersColumn_State     = "state"
	UsersColumn_Email     = "email"
	UsersColumn_ExtraData = "extra_data"
)

const (
	UserState_Unknown UserState = iota
	UserState_Active
	UserState_Baned
)

type UserState int

type User struct {
	Id        uint64        `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`
	UserName  string        `json:"user_name" db:"user_name" gorm:"column:user_name;type:varchar(32);uniqueIndex:idx_users_user_name,priority:1;null;" sqlca:"nullable"`
	State     UserState     `json:"state" db:"state" gorm:"column:state;type:tinyint(1);default:0;null;" sqlca:"nullable"`
	Email     string        `json:"email" db:"email" gorm:"column:email;type:varchar(64);uniqueIndex:idx_users_email,priority:1;null;" sqlca:"nullable"`
	ExtraData UserExtraData `json:"extra_data" db:"extra_data" gorm:"column:extra_data;type:json;null;" sqlca:"nullable"`
	Roles     []*Role       `json:"roles,omitempty" db:"-" gorm:"many2many:user_roles;-:migration;"` // 用户角色列表
	Profile   UserProfile   `json:"profile,omitempty" db:"-" gorm:"foreignKey:UserId;-:migration;"`  // 用户资料明细
	BaseModel
}

func (do User) DatabaseName() string { return "test" }

func (do User) TableName() string { return TableNameUsers }

func (do User) GetId() uint64 { return do.Id }

func (do User) GetCreatedAt() time.Time { return do.CreatedAt }

func (do User) GetUpdatedAt() time.Time { return do.UpdatedAt }

func (do User) GetUserName() string { return do.UserName }

func (do User) GetState() UserState { return do.State }

func (do User) GetEmail() string { return do.Email }

func (do User) GetExtraData() UserExtraData { return do.ExtraData }

func (do *User) SetId(v uint64) { do.Id = v }

func (do *User) SetCreatedAt(v time.Time) { do.CreatedAt = v }

func (do *User) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }

func (do *User) SetUserName(v string) { do.UserName = v }

func (do *User) SetState(v UserState) { do.State = v }

func (do *User) SetEmail(v string) { do.Email = v }

func (do *User) SetExtraData(v UserExtraData) { do.ExtraData = v }
