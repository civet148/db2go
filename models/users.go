package models

const TableNameUsers = "users"

const (
	UsersColumn_Id        = "id"
	UsersColumn_UserName  = "user_name"
	UsersColumn_State     = "state"
	UsersColumn_Email     = "email"
	UsersColumn_ExtraData = "extra_data"
	UsersColumn_LeftVal   = "left_val"
	UsersColumn_RightVal  = "right_val"
	UsersColumn_Depth     = "depth"
	UsersColumn_ParentId  = "parent_id"
	UsersColumn_RootPath  = "root_path"
)

const (
	UserState_Unknown UserState = iota
	UserState_Active
	UserState_Baned
)

type UserState int

type User struct {
	Id        uint64        `json:"id" gorm:"column:id;primaryKey;autoIncrement;"`
	UserName  string        `json:"user_name" gorm:"column:user_name;type:varchar(32);uniqueIndex:idx_users_user_name,priority:1;null;" sqlca:"nullable"`
	State     UserState     `json:"state" gorm:"column:state;type:tinyint(1);default:0;null;" sqlca:"nullable"`
	Email     string        `json:"email" gorm:"column:email;type:varchar(64);uniqueIndex:idx_users_email,priority:1;null;" sqlca:"nullable"`
	ExtraData UserExtraData `json:"extra_data" gorm:"column:extra_data;type:json;null;" sqlca:"nullable"`
	LeftVal   int64         `json:"left_val" gorm:"column:left_val;type:bigint;index:idx_users_left_val,priority:1;null;" sqlca:"nullable"`
	RightVal  int64         `json:"right_val" gorm:"column:right_val;type:bigint;index:idx_users_right_val,priority:1;null;" sqlca:"nullable"`
	Depth     int64         `json:"depth" gorm:"column:depth;type:bigint;null;" sqlca:"nullable"`
	ParentId  int64         `json:"parent_id" gorm:"column:parent_id;type:bigint;default:0;null;" sqlca:"nullable"`
	RootPath  []int64       `json:"root_path" gorm:"column:root_path;type:json;null;" sqlca:"nullable"`
	Roles     []*Role       `json:"roles,omitempty" db:"-" gorm:"many2many:user_roles;-:migration;"` // 用户角色列表(手工添加，自动合并)
	Profile   UserProfile   `json:"profile,omitempty" db:"-" gorm:"foreignKey:UserId;-:migration;"`  // 用户资料明细(手工添加，自动合并)
	BaseModel `json:"base_model" db:"-" gorm:"embedded"`
}

func (do User) DatabaseName() string { return "test" }

func (do User) TableName() string { return TableNameUsers }

func (do User) GetId() uint64 { return do.Id }

func (do User) GetUserName() string { return do.UserName }

func (do User) GetState() UserState { return do.State }

func (do User) GetEmail() string { return do.Email }

func (do User) GetExtraData() UserExtraData { return do.ExtraData }

func (do User) GetLeftVal() int64 { return do.LeftVal }

func (do User) GetRightVal() int64 { return do.RightVal }

func (do User) GetDepth() int64 { return do.Depth }

func (do User) GetParentId() int64 { return do.ParentId }

func (do User) GetRootPath() []int64 { return do.RootPath }

func (do *User) SetId(v uint64) { do.Id = v }

func (do *User) SetUserName(v string) { do.UserName = v }

func (do *User) SetState(v UserState) { do.State = v }

func (do *User) SetEmail(v string) { do.Email = v }

func (do *User) SetExtraData(v UserExtraData) { do.ExtraData = v }

func (do *User) SetLeftVal(v int64) { do.LeftVal = v }

func (do *User) SetRightVal(v int64) { do.RightVal = v }

func (do *User) SetDepth(v int64) { do.Depth = v }

func (do *User) SetParentId(v int64) { do.ParentId = v }

func (do *User) SetRootPath(v []int64) { do.RootPath = v }
