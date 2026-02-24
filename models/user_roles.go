package models

import github_com_civet148_db2go_types "github.com/civet148/db2go/types"

const TableNameUserRoles = "user_roles" //

const (
	USER_ROLES_COLUMN_ROLE_ID = "role_id"
	USER_ROLES_COLUMN_USER_ID = "user_id"
)

type UserRole struct {
	github_com_civet148_db2go_types.BaseModel
	RoleId uint64 `json:"role_id,omitempty" db:"role_id" gorm:"column:role_id;primaryKey;autoIncrement;"` //
	UserId uint64 `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;primaryKey;autoIncrement;"` //
}

func (do UserRole) TableName() string { return "user_roles" }

func (do UserRole) GetRoleId() uint64 { return do.RoleId }
func (do UserRole) GetUserId() uint64 { return do.UserId }

func (do *UserRole) SetRoleId(v uint64) { do.RoleId = v }
func (do *UserRole) SetUserId(v uint64) { do.UserId = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
