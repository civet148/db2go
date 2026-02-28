package models

import "time"

const TableNameRoles = "roles" //

const (
	ROLES_COLUMN_ID         = "id"
	ROLES_COLUMN_CREATED_AT = "created_at"
	ROLES_COLUMN_UPDATED_AT = "updated_at"
	ROLES_COLUMN_NAME       = "name"
)

type Role struct {
	BaseModel
	Id   uint64 `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                         //
	Name string `json:"name,omitempty" db:"name" gorm:"column:name;type:varchar(64);uniqueIndex:idx_roles_name;" sqlca:"isnull"` //
}

func (do Role) TableName() string { return "roles" }

func (do Role) GetId() uint64           { return do.Id }
func (do Role) GetCreatedAt() time.Time { return do.CreatedAt }
func (do Role) GetUpdatedAt() time.Time { return do.UpdatedAt }
func (do Role) GetName() string         { return do.Name }

func (do *Role) SetId(v uint64)           { do.Id = v }
func (do *Role) SetCreatedAt(v time.Time) { do.CreatedAt = v }
func (do *Role) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }
func (do *Role) SetName(v string)         { do.Name = v }

/*
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_roles_name` (`name`),
  KEY `idx_roles_created_at` (`created_at`),
  KEY `idx_roles_updated_at` (`updated_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
