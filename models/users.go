package models

const TableNameUsers = "users" //

const (
	UsersColumn_Id        = "id"
	UsersColumn_CreatedAt = "created_at"
	UsersColumn_UpdatedAt = "updated_at"
	UsersColumn_UserName  = "user_name"
	UsersColumn_Email     = "email"
)

type User struct {
	Id       uint64 `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                          //
	UserName string `json:"user_name" db:"user_name" gorm:"column:user_name;type:varchar(32);uniqueIndex:idx_users_user_name;default:null;" sqlca:"isnull"` //
	Email    string `json:"email" db:"email" gorm:"column:email;type:varchar(64);uniqueIndex:idx_users_email;default:null;" sqlca:"isnull"`                 //
}

func (do User) TableName() string { return "users" }

func (do User) GetId() uint64 { return do.Id }

func (do User) GetUserName() string { return do.UserName }

func (do User) GetEmail() string { return do.Email }

func (do *User) SetId(v uint64) { do.Id = v }

func (do *User) SetUserName(v string) { do.UserName = v }

func (do *User) SetEmail(v string) { do.Email = v }
