package models

const TableNameUserRoles = "user_roles" //

const (
	UserRolesColumn_UserId    = "user_id"
	UserRolesColumn_RoleId    = "role_id"
	UserRolesColumn_CreatedAt = "created_at"
	UserRolesColumn_UpdatedAt = "updated_at"
)

type UserRole struct {
	UserId uint64 `json:"user_id" db:"user_id" gorm:"column:user_id;type:bigint unsigned;;default:null;"`                          //
	RoleId uint64 `json:"role_id" db:"role_id" gorm:"column:role_id;type:bigint unsigned;index:fk_user_roles_role;;default:null;"` //
}

func (do UserRole) TableName() string { return "user_roles" }

func (do UserRole) GetUserId() uint64 { return do.UserId }

func (do UserRole) GetRoleId() uint64 { return do.RoleId }

func (do *UserRole) SetUserId(v uint64) { do.UserId = v }

func (do *UserRole) SetRoleId(v uint64) { do.RoleId = v }
