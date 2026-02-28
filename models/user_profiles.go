package models

import "time"

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
	BaseModel
	Id      uint64 `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                 //
	UserId  uint64 `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;type:bigint unsigned;uniqueIndex:idx_user_profiles_user_id;" sqlca:"isnull"` //
	Avatar  string `json:"avatar,omitempty" db:"avatar" gorm:"column:avatar;type:varchar(512);" sqlca:"isnull"`                                             //
	Address string `json:"address,omitempty" db:"address" gorm:"column:address;type:varchar(128);" sqlca:"isnull"`                                          //
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

/*
CREATE TABLE `user_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` bigint unsigned DEFAULT NULL,
  `avatar` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `address` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_profiles_user_id` (`user_id`),
  KEY `idx_user_profiles_created_at` (`created_at`),
  KEY `idx_user_profiles_updated_at` (`updated_at`),
  CONSTRAINT `fk_users_profile` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
