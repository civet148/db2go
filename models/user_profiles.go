package models

const TableNameUserProfiles = "user_profiles" //

const (
	UserProfilesColumn_Id        = "id"
	UserProfilesColumn_CreatedAt = "created_at"
	UserProfilesColumn_UpdatedAt = "updated_at"
	UserProfilesColumn_UserId    = "user_id"
	UserProfilesColumn_Avatar    = "avatar"
	UserProfilesColumn_Address   = "address"
)

type UserProfile struct {
	Id      uint64 `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                              //
	UserId  uint64 `json:"user_id" db:"user_id" gorm:"column:user_id;type:bigint unsigned;uniqueIndex:idx_user_profiles_user_id;default:null;" sqlca:"isnull"` //
	Avatar  string `json:"avatar" db:"avatar" gorm:"column:avatar;type:varchar(512);default:null;" sqlca:"isnull"`                                             //
	Address string `json:"address" db:"address" gorm:"column:address;type:varchar(128);default:null;" sqlca:"isnull"`                                          //
}

func (do UserProfile) TableName() string { return "user_profiles" }

func (do UserProfile) GetId() uint64 { return do.Id }

func (do UserProfile) GetUserId() uint64 { return do.UserId }

func (do UserProfile) GetAvatar() string { return do.Avatar }

func (do UserProfile) GetAddress() string { return do.Address }

func (do *UserProfile) SetId(v uint64) { do.Id = v }

func (do *UserProfile) SetUserId(v uint64) { do.UserId = v }

func (do *UserProfile) SetAvatar(v string) { do.Avatar = v }

func (do *UserProfile) SetAddress(v string) { do.Address = v }
