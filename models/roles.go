package models

const TableNameRoles = "roles" //

const (
	RolesColumn_Id        = "id"
	RolesColumn_CreatedAt = "created_at"
	RolesColumn_UpdatedAt = "updated_at"
	RolesColumn_Name      = "name"
)

type Role struct {
	Id   uint64 `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                      //
	Name string `json:"name" db:"name" gorm:"column:name;type:varchar(64);uniqueIndex:idx_roles_name;default:null;" sqlca:"isnull"` //
}

func (do Role) TableName() string { return "roles" }

func (do Role) GetId() uint64 { return do.Id }

func (do Role) GetName() string { return do.Name }

func (do *Role) SetId(v uint64) { do.Id = v }

func (do *Role) SetName(v string) { do.Name = v }
