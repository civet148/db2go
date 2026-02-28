package models

const TableNameRoleUsers = "role_users" //

const (
	ROLE_USERS_COLUMN_USER_ID = "user_id"
	ROLE_USERS_COLUMN_ROLE_ID = "role_id"
)

type RoleUser struct {
	BaseModel
	UserId uint64 `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;type:bigint unsigned;uniqueIndex:PRIMARY;"`                          //
	RoleId uint64 `json:"role_id,omitempty" db:"role_id" gorm:"column:role_id;type:bigint unsigned;index:fk_role_users_role;uniqueIndex:PRIMARY;"` //
}

func (do RoleUser) TableName() string { return "role_users" }

func (do RoleUser) GetUserId() uint64 { return do.UserId }
func (do RoleUser) GetRoleId() uint64 { return do.RoleId }

func (do *RoleUser) SetUserId(v uint64) { do.UserId = v }
func (do *RoleUser) SetRoleId(v uint64) { do.RoleId = v }

/*
CREATE TABLE `role_users` (
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `fk_role_users_role` (`role_id`),
  CONSTRAINT `fk_role_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `fk_role_users_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
