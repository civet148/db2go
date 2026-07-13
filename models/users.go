package models

const TableNameUsers = "users" //

const (
	USERS_COLUMN_ID         = "id"
	USERS_COLUMN_CREATED_AT = "created_at"
	USERS_COLUMN_UPDATED_AT = "updated_at"
	USERS_COLUMN_USER_NAME  = "user_name"
	USERS_COLUMN_EMAIL      = "email"
)

type User struct {
	Id       uint64 `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                          // 用户ID
	UserName string `json:"user_name" db:"user_name" gorm:"column:user_name;type:varchar(32);uniqueIndex:idx_users_user_name;default:null;" sqlca:"isnull"` // 用户名
	Email    string `json:"email" db:"email" gorm:"column:email;type:varchar(64);uniqueIndex:idx_users_email;default:null;" sqlca:"isnull"`                 // 邮箱地址
	BaseModel
	Roles   []*Role     `json:"roles" db:"-" gorm:"many2many:user_roles"`
	Profile UserProfile `json:"profile" db:"-" gorm:"foreignKey:UserId;"`
}

func (do User) TableName() string { return "users" }

func (do User) GetId() uint64 { return do.Id }

func (do User) GetUserName() string { return do.UserName }

func (do User) GetEmail() string { return do.Email }

func (do *User) SetId(v uint64) { do.Id = v }

func (do *User) SetUserName(v string) { do.UserName = v }

func (do *User) SetEmail(v string) { do.Email = v }
