package models

import github_com_civet148_db2go_types "github.com/civet148/db2go/types"

const TableNameUserRoleId = "user_role_id" //

const (
	USER_ROLE_ID_COLUMN_ID         = "id"
	USER_ROLE_ID_COLUMN_USER_ID    = "user_id"
	USER_ROLE_ID_COLUMN_ROLE_ID    = "role_id"
	USER_ROLE_ID_COLUMN_CREATED_AT = "created_at"
)

type UserRoleId struct {
	github_com_civet148_db2go_types.BaseModel
	Id        int64  `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                                     //自增ID
	UserId    int64  `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;type:bigint;uniqueIndex:user_role_id_unique;comment:用户ID;"`                                      //用户ID
	RoleId    int64  `json:"role_id,omitempty" db:"role_id" gorm:"column:role_id;type:bigint;uniqueIndex:user_role_id_unique;comment:角色ID;"`                                      //角色ID
	CreatedAt string `json:"created_at,omitempty" db:"created_at" gorm:"column:created_at;type:datetime;autoCreateTime;default:CURRENT_TIMESTAMP;comment:创建时间;" sqlca:"readonly"` //创建时间
}

func (do UserRoleId) TableName() string { return "user_role_id" }

func (do UserRoleId) GetId() int64         { return do.Id }
func (do UserRoleId) GetUserId() int64     { return do.UserId }
func (do UserRoleId) GetRoleId() int64     { return do.RoleId }
func (do UserRoleId) GetCreatedAt() string { return do.CreatedAt }

func (do *UserRoleId) SetId(v int64)         { do.Id = v }
func (do *UserRoleId) SetUserId(v int64)     { do.UserId = v }
func (do *UserRoleId) SetRoleId(v int64)     { do.RoleId = v }
func (do *UserRoleId) SetCreatedAt(v string) { do.CreatedAt = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
