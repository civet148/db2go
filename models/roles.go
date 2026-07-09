package models

const TableNameRoles = "roles" //

const (
	ROLES_COLUMN_ID         = "id"
	ROLES_COLUMN_CREATED_AT = "created_at"
	ROLES_COLUMN_UPDATED_AT = "updated_at"
	ROLES_COLUMN_NAME       = "name"
)

type Role struct {
	Id   uint64 `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                         //
	Name string `json:"name,omitempty" db:"name" gorm:"column:name;type:varchar(64);uniqueIndex:idx_roles_name;" sqlca:"isnull"` //
}

func (do Role) TableName() string { return "roles" }

func (do Role) GetId() uint64 { return do.Id }

func (do Role) GetName() string { return do.Name }

func (do *Role) SetId(v uint64) { do.Id = v }

func (do *Role) SetName(v string) { do.Name = v }
