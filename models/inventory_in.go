package models

import github_com_civet148_db2go_types "github.com/civet148/db2go/types"
import "github.com/civet148/sqlca/v3"

const TableNameInventoryIn = "inventory_in" //入库主表

const (
	INVENTORY_IN_COLUMN_ID          = "id"
	INVENTORY_IN_COLUMN_CREATE_ID   = "create_id"
	INVENTORY_IN_COLUMN_CREATE_NAME = "create_name"
	INVENTORY_IN_COLUMN_CREATE_TIME = "create_time"
	INVENTORY_IN_COLUMN_UPDATE_ID   = "update_id"
	INVENTORY_IN_COLUMN_UPDATE_NAME = "update_name"
	INVENTORY_IN_COLUMN_UPDATE_TIME = "update_time"
	INVENTORY_IN_COLUMN_IS_DELETED  = "is_deleted"
	INVENTORY_IN_COLUMN_DELETE_TIME = "delete_time"
	INVENTORY_IN_COLUMN_PRODUCT_ID  = "product_id"
	INVENTORY_IN_COLUMN_ORDER_NO    = "order_no"
	INVENTORY_IN_COLUMN_USER_ID     = "user_id"
	INVENTORY_IN_COLUMN_USER_NAME   = "user_name"
	INVENTORY_IN_COLUMN_QUANTITY    = "quantity"
	INVENTORY_IN_COLUMN_WEIGHT      = "weight"
	INVENTORY_IN_COLUMN_REMARK      = "remark"
)

type InventoryIn struct {
	github_com_civet148_db2go_types.BaseModel
	Id         uint64        `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                //主键ID
	CreateId   uint64        `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id;type:bigint unsigned;default:0;"`     //创建人ID
	CreateName string        `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name;type:varchar(64);"`             //创建人姓名
	UpdateId   uint64        `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id;type:bigint unsigned;default:0;"`     //更新人ID
	UpdateName string        `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name;type:varchar(64);"`             //更新人姓名
	IsDeleted  int8          `json:"is_deleted,omitempty" db:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;"`       //删除状态(0: 未删除 1: 已删除)
	DeleteTime string        `json:"delete_time,omitempty" db:"delete_time" gorm:"column:delete_time;type:datetime;" sqlca:"isnull"` //删除时间
	ProductId  uint64        `json:"product_id,omitempty" db:"product_id" gorm:"column:product_id;type:bigint unsigned;"`            //产品ID
	OrderNo    string        `json:"order_no,omitempty" db:"order_no" gorm:"column:order_no;type:varchar(64);"`                      //入库单号
	UserId     uint64        `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;type:bigint unsigned;default:0;"`           //交货人ID
	UserName   string        `json:"user_name,omitempty" db:"user_name" gorm:"column:user_name;type:varchar(64);"`                   //交货人姓名
	Quantity   sqlca.Decimal `json:"quantity,omitempty" db:"quantity" gorm:"column:quantity;type:decimal(16,6);default:0.000000;"`   //数量
	Weight     sqlca.Decimal `json:"weight,omitempty" db:"weight" gorm:"column:weight;type:decimal(16,6);default:0.000000;"`         //净重
	Remark     string        `json:"remark,omitempty" db:"remark" gorm:"column:remark;type:varchar(512);"`                           //备注
}

func (do *InventoryIn) GetId() uint64               { return do.Id }
func (do *InventoryIn) SetId(v uint64)              { do.Id = v }
func (do *InventoryIn) GetCreateId() uint64         { return do.CreateId }
func (do *InventoryIn) SetCreateId(v uint64)        { do.CreateId = v }
func (do *InventoryIn) GetCreateName() string       { return do.CreateName }
func (do *InventoryIn) SetCreateName(v string)      { do.CreateName = v }
func (do *InventoryIn) GetCreateTime() string       { return do.CreateTime }
func (do *InventoryIn) SetCreateTime(v string)      { do.CreateTime = v }
func (do *InventoryIn) GetUpdateId() uint64         { return do.UpdateId }
func (do *InventoryIn) SetUpdateId(v uint64)        { do.UpdateId = v }
func (do *InventoryIn) GetUpdateName() string       { return do.UpdateName }
func (do *InventoryIn) SetUpdateName(v string)      { do.UpdateName = v }
func (do *InventoryIn) GetUpdateTime() string       { return do.UpdateTime }
func (do *InventoryIn) SetUpdateTime(v string)      { do.UpdateTime = v }
func (do *InventoryIn) GetIsDeleted() int8          { return do.IsDeleted }
func (do *InventoryIn) SetIsDeleted(v int8)         { do.IsDeleted = v }
func (do *InventoryIn) GetDeleteTime() string       { return do.DeleteTime }
func (do *InventoryIn) SetDeleteTime(v string)      { do.DeleteTime = v }
func (do *InventoryIn) GetProductId() uint64        { return do.ProductId }
func (do *InventoryIn) SetProductId(v uint64)       { do.ProductId = v }
func (do *InventoryIn) GetOrderNo() string          { return do.OrderNo }
func (do *InventoryIn) SetOrderNo(v string)         { do.OrderNo = v }
func (do *InventoryIn) GetUserId() uint64           { return do.UserId }
func (do *InventoryIn) SetUserId(v uint64)          { do.UserId = v }
func (do *InventoryIn) GetUserName() string         { return do.UserName }
func (do *InventoryIn) SetUserName(v string)        { do.UserName = v }
func (do *InventoryIn) GetQuantity() sqlca.Decimal  { return do.Quantity }
func (do *InventoryIn) SetQuantity(v sqlca.Decimal) { do.Quantity = v }
func (do *InventoryIn) GetWeight() sqlca.Decimal    { return do.Weight }
func (do *InventoryIn) SetWeight(v sqlca.Decimal)   { do.Weight = v }
func (do *InventoryIn) GetRemark() string           { return do.Remark }
func (do *InventoryIn) SetRemark(v string)          { do.Remark = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
