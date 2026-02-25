package models

const TableNameUserRoles = "user_roles" //

const (
	USER_ROLES_COLUMN_USER_ID = "user_id"
	USER_ROLES_COLUMN_ROLE_ID = "role_id"
)

type UserRole struct {
	BaseModel
	UserId uint64 `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;type:bigint unsigned;uniqueIndex:PRIMARY;"`                          //
	RoleId uint64 `json:"role_id,omitempty" db:"role_id" gorm:"column:role_id;type:bigint unsigned;index:fk_user_roles_role;uniqueIndex:PRIMARY;"` //
}

func (do UserRole) TableName() string { return "user_roles" }

func (do UserRole) GetUserId() uint64 { return do.UserId }
func (do UserRole) GetRoleId() uint64 { return do.RoleId }

func (do *UserRole) SetUserId(v uint64) { do.UserId = v }
func (do *UserRole) SetRoleId(v uint64) { do.RoleId = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
