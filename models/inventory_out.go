package models

import "time"
import "github.com/civet148/sqlca/v3"

const TableNameInventoryOut = "inventory_out" //

const (
	INVENTORY_OUT_COLUMN_ID          = "id"
	INVENTORY_OUT_COLUMN_CREATED_AT  = "created_at"
	INVENTORY_OUT_COLUMN_UPDATED_AT  = "updated_at"
	INVENTORY_OUT_COLUMN_IS_DELETED  = "is_deleted"
	INVENTORY_OUT_COLUMN_DELETE_TIME = "delete_time"
	INVENTORY_OUT_COLUMN_PRODUCT_ID  = "product_id"
	INVENTORY_OUT_COLUMN_ORDER_NO    = "order_no"
	INVENTORY_OUT_COLUMN_USER_ID     = "user_id"
	INVENTORY_OUT_COLUMN_USER_NAME   = "user_name"
	INVENTORY_OUT_COLUMN_QUANTITY    = "quantity"
	INVENTORY_OUT_COLUMN_WEIGHT      = "weight"
	INVENTORY_OUT_COLUMN_REMARK      = "remark"
	INVENTORY_OUT_COLUMN_CREATE_ID   = "create_id"
	INVENTORY_OUT_COLUMN_CREATE_NAME = "create_name"
	INVENTORY_OUT_COLUMN_UPDATE_ID   = "update_id"
	INVENTORY_OUT_COLUMN_UPDATE_NAME = "update_name"
)

type InventoryOut struct {
	BaseModel
	Id         uint64        `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                                            //
	IsDeleted  int8          `json:"is_deleted,omitempty" db:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;" sqlca:"isnull"`                                                    //
	DeleteTime *time.Time    `json:"delete_time,omitempty" db:"delete_time" gorm:"column:delete_time;type:datetime;" sqlca:"isnull"`                                                             //
	ProductId  uint64        `json:"product_id,omitempty" db:"product_id" gorm:"column:product_id;type:bigint unsigned;index:i_product_id;uniqueIndex:UNIQ_PROD_USER;default:0;" sqlca:"isnull"` //
	OrderNo    string        `json:"order_no,omitempty" db:"order_no" gorm:"column:order_no;type:varchar(64);uniqueIndex:UNIQ_ORDER_NO;" sqlca:"isnull"`                                         //
	UserId     uint64        `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;type:bigint unsigned;index:i_user_id;uniqueIndex:UNIQ_PROD_USER;default:0;" sqlca:"isnull"`             //
	UserName   string        `json:"user_name,omitempty" db:"user_name" gorm:"column:user_name;type:varchar(64);index:FULTXT_user_name;" sqlca:"isnull"`                                         //
	Quantity   sqlca.Decimal `json:"quantity,omitempty" db:"quantity" gorm:"column:quantity;type:decimal(16,6);default:0.000000;" sqlca:"isnull"`                                                //
	Weight     sqlca.Decimal `json:"weight,omitempty" db:"weight" gorm:"column:weight;type:decimal(16,6);default:0.000000;" sqlca:"isnull"`                                                      //
	Remark     string        `json:"remark,omitempty" db:"remark" gorm:"column:remark;type:varchar(512);" sqlca:"isnull"`                                                                        //
	CreateId   uint64        `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id;type:bigint unsigned;default:0;" sqlca:"isnull"`                                                  //
	CreateName string        `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name;type:varchar(64);" sqlca:"isnull"`                                                          //
	UpdateId   uint64        `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id;type:bigint unsigned;default:0;" sqlca:"isnull"`                                                  //
	UpdateName string        `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name;type:varchar(64);" sqlca:"isnull"`                                                          //
}

func (do InventoryOut) TableName() string { return "inventory_out" }

func (do InventoryOut) GetId() uint64              { return do.Id }
func (do InventoryOut) GetCreatedAt() time.Time    { return do.CreatedAt }
func (do InventoryOut) GetUpdatedAt() time.Time    { return do.UpdatedAt }
func (do InventoryOut) GetIsDeleted() int8         { return do.IsDeleted }
func (do InventoryOut) GetDeleteTime() *time.Time  { return do.DeleteTime }
func (do InventoryOut) GetProductId() uint64       { return do.ProductId }
func (do InventoryOut) GetOrderNo() string         { return do.OrderNo }
func (do InventoryOut) GetUserId() uint64          { return do.UserId }
func (do InventoryOut) GetUserName() string        { return do.UserName }
func (do InventoryOut) GetQuantity() sqlca.Decimal { return do.Quantity }
func (do InventoryOut) GetWeight() sqlca.Decimal   { return do.Weight }
func (do InventoryOut) GetRemark() string          { return do.Remark }
func (do InventoryOut) GetCreateId() uint64        { return do.CreateId }
func (do InventoryOut) GetCreateName() string      { return do.CreateName }
func (do InventoryOut) GetUpdateId() uint64        { return do.UpdateId }
func (do InventoryOut) GetUpdateName() string      { return do.UpdateName }

func (do *InventoryOut) SetId(v uint64)              { do.Id = v }
func (do *InventoryOut) SetCreatedAt(v time.Time)    { do.CreatedAt = v }
func (do *InventoryOut) SetUpdatedAt(v time.Time)    { do.UpdatedAt = v }
func (do *InventoryOut) SetIsDeleted(v int8)         { do.IsDeleted = v }
func (do *InventoryOut) SetDeleteTime(v *time.Time)  { do.DeleteTime = v }
func (do *InventoryOut) SetProductId(v uint64)       { do.ProductId = v }
func (do *InventoryOut) SetOrderNo(v string)         { do.OrderNo = v }
func (do *InventoryOut) SetUserId(v uint64)          { do.UserId = v }
func (do *InventoryOut) SetUserName(v string)        { do.UserName = v }
func (do *InventoryOut) SetQuantity(v sqlca.Decimal) { do.Quantity = v }
func (do *InventoryOut) SetWeight(v sqlca.Decimal)   { do.Weight = v }
func (do *InventoryOut) SetRemark(v string)          { do.Remark = v }
func (do *InventoryOut) SetCreateId(v uint64)        { do.CreateId = v }
func (do *InventoryOut) SetCreateName(v string)      { do.CreateName = v }
func (do *InventoryOut) SetUpdateId(v uint64)        { do.UpdateId = v }
func (do *InventoryOut) SetUpdateName(v string)      { do.UpdateName = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
